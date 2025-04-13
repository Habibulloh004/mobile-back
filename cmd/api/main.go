package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"mobilka/config"
	"mobilka/internal/api/routes"
	"mobilka/internal/repository"
	"mobilka/internal/service"
	"mobilka/internal/tasks"
	"mobilka/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize JWT secret
	if err := utils.InitJWTSecret(); err != nil {
		log.Fatalf("Failed to initialize JWT secret: %v", err)
	}

	// Connect to database
	db, err := connectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run database migrations
	migrationPath := "./migrations"
	if err := utils.RunMigrations(context.Background(), db, migrationPath); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Setup super admin account
	if err := setupSuperAdmin(db); err != nil {
		log.Fatalf("Failed to setup super admin: %v", err)
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Fiber App",
		ErrorHandler: errorHandler,
		BodyLimit:    utils.MaxImageSize + 1024*1024, // Max image size + 1MB for other data
	})

	// Setup routes
	routes.SetupRoutes(app, db, cfg)

	// Ensure upload directories exist
	uploadDirs := []string{cfg.ImageUploadPath}
	for _, dir := range uploadDirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Failed to create upload directory %s: %v", dir, err)
		}
	}

	// Start subscription checker task
	subscriptionChecker := setupSubscriptionChecker(db)
	subscriptionChecker.Start()

	// Print startup information
	log.Printf("Server starting on port %d", cfg.ServerPort)
	log.Printf("Environment: %s", cfg.Environment)

	// Start server
	go func() {
		if err := app.Listen(":" + strconv.Itoa(cfg.ServerPort)); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Stop subscription checker
	subscriptionChecker.Stop()

	// Shutdown server with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	log.Println("Server gracefully stopped")
}

// Connect to the database
func connectDB(cfg *config.Config) (*pgxpool.Pool, error) {
	connStr := cfg.GetDBConnString()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to database with retry logic
	var db *pgxpool.Pool
	var err error

	maxRetries := 5
	retryDelay := 3 * time.Second

	for i := 0; i < maxRetries; i++ {
		log.Printf("Attempting to connect to database (attempt %d/%d)...", i+1, maxRetries)

		db, err = pgxpool.New(ctx, connStr)
		if err == nil {
			// Test connection
			if err = db.Ping(ctx); err == nil {
				log.Println("Connected to database successfully")
				return db, nil
			}
		}

		log.Printf("Failed to connect to database: %v. Retrying in %v...", err, retryDelay)
		time.Sleep(retryDelay)
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}

// Setup super admin account
func setupSuperAdmin(db *pgxpool.Pool) error {
	// Create repositories
	superAdminRepo := repository.NewSuperAdminRepository(db)
	adminRepo := repository.NewAdminRepository(db)

	// Create auth service
	authService := service.NewAuthService(superAdminRepo, adminRepo)

	// Setup super admin account
	password, err := authService.SetupDefaultSuperAdmin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to setup super admin: %w", err)
	}

	// Log super admin credentials (only during setup)
	log.Println("====== SUPER ADMIN CREDENTIALS ======")
	log.Println("Login: superadmin")
	log.Println("Password: " + password)
	log.Println("PLEASE SAVE THESE CREDENTIALS SECURELY")
	log.Println("====================================")

	return nil
}

// Setup subscription checker task
func setupSubscriptionChecker(db *pgxpool.Pool) *tasks.SubscriptionChecker {
	// Create repositories needed for the subscription checker
	adminRepo := repository.NewAdminRepository(db)
	paymentRepo := repository.NewPaymentHistoryRepository(db)
	subscriptionTierRepo := repository.NewSubscriptionTierRepository(db)

	// Create payment service
	paymentService := service.NewPaymentService(paymentRepo, adminRepo, subscriptionTierRepo)

	// Create subscription checker with 12-hour interval
	return tasks.NewSubscriptionChecker(paymentService, 12*time.Hour)
}

// Custom error handler
func errorHandler(c *fiber.Ctx, err error) error {
	// Default 500 status code
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		// Override status code if it's a Fiber error
		code = e.Code
	}

	// Return JSON response
	return c.Status(code).JSON(fiber.Map{
		"status":  utils.StatusError,
		"message": err.Error(),
	})
}

package routes

import (
	"mobilka/config"
	"mobilka/internal/api/handlers"
	"mobilka/internal/repository"
	"mobilka/internal/service"
	"mobilka/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SetupRoutes sets up all the routes for the application
func SetupRoutes(app *fiber.App, db *pgxpool.Pool, cfg *config.Config) {
	// Apply global middlewares
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// Create repositories
	superAdminRepo := repository.NewSuperAdminRepository(db)
	adminRepo := repository.NewAdminRepository(db)
	bannerRepo := repository.NewBannerRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	fcmTokenRepo := repository.NewFCMTokenRepository(db)

	// Create services
	authService := service.NewAuthService(superAdminRepo, adminRepo)
	superAdminService := service.NewSuperAdminService(superAdminRepo)
	adminService := service.NewAdminService(adminRepo)
	bannerService := service.NewBannerService(bannerRepo)
	notificationService := service.NewNotificationService(notificationRepo, fcmTokenRepo)
	fcmTokenService := service.NewFCMTokenService(fcmTokenRepo)
	imageService := service.NewImageService(cfg.ImageUploadPath)

	// Create handlers
	authHandler := handlers.NewAuthHandler(authService)
	superAdminHandler := handlers.NewSuperAdminHandler(superAdminService)
	adminHandler := handlers.NewAdminHandler(adminService)
	bannerHandler := handlers.NewBannerHandler(bannerService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	fcmTokenHandler := handlers.NewFCMTokenHandler(fcmTokenService)
	imageHandler := handlers.NewImageHandler(imageService)

	// Setup API routes
	api := app.Group("/api")

	// Setup modular routes
	SetupAuthRoutes(api, authHandler)
	SetupSuperAdminRoutes(api, superAdminHandler, authHandler)
	SetupAdminRoutes(api, adminHandler)
	SetupBannerRoutes(api, bannerHandler)
	SetupNotificationRoutes(api, notificationHandler)
	SetupFCMTokenRoutes(api, fcmTokenHandler)
	SetupImageRoutes(app, api, imageHandler)

	// Setup 404 handler
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Resource not found",
		})
	})
}

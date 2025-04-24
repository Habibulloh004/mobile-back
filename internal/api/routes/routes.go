// Here's how to update the SetupRoutes function in internal/api/routes/routes.go

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
	restaurantRepo := repository.NewRestaurantRepository(db) // Add new repository

	subscriptionTierRepo := repository.NewSubscriptionTierRepository(db)
	paymentRepo := repository.NewPaymentHistoryRepository(db)

	// Create services
	authService := service.NewAuthService(superAdminRepo, adminRepo)
	superAdminService := service.NewSuperAdminService(superAdminRepo)
	adminService := service.NewAdminService(adminRepo)
	bannerService := service.NewBannerService(bannerRepo)
	notificationService := service.NewNotificationService(notificationRepo, fcmTokenRepo)
	fcmTokenService := service.NewFCMTokenService(fcmTokenRepo)
	imageService := service.NewImageService(cfg.ImageUploadPath)
	restaurantService := service.NewRestaurantService(restaurantRepo) // Add new service

	subscriptionTierService := service.NewSubscriptionTierService(subscriptionTierRepo)
	paymentService := service.NewPaymentService(paymentRepo, adminRepo, subscriptionTierRepo)

	// Create handlers
	authHandler := handlers.NewAuthHandler(authService)
	superAdminHandler := handlers.NewSuperAdminHandler(superAdminService)
	adminHandler := handlers.NewAdminHandler(adminService)
	bannerHandler := handlers.NewBannerHandler(bannerService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	fcmTokenHandler := handlers.NewFCMTokenHandler(fcmTokenService)
	imageHandler := handlers.NewImageHandler(imageService)
	restaurantHandler := handlers.NewRestaurantHandler(restaurantService) // Add new handler

	subscriptionTierHandler := handlers.NewSubscriptionTierHandler(subscriptionTierService)
	paymentHandler := handlers.NewPaymentHandler(paymentService)

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
	SetupRestaurantRoutes(api, restaurantHandler) // Add new routes

	// Setup public routes
	publicRoutes := api.Group("/public")
	SetupPublicRoutes(publicRoutes, bannerHandler, notificationHandler, restaurantHandler) // Update public routes

	SetupSubscriptionTierRoutes(api, subscriptionTierHandler)
	SetupPaymentRoutes(api, paymentHandler, subscriptionTierHandler)

	// Setup 404 handler
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Resource not found",
		})
	})
}

// SetupPublicRoutes sets up all the public routes
func SetupPublicRoutes(publicRoutes fiber.Router, bannerHandler *handlers.BannerHandler,
	notificationHandler *handlers.NotificationHandler,
	restaurantHandler *handlers.RestaurantHandler) {

	// Banner routes
	publicRoutes.Get("/banners/admin/:adminID", bannerHandler.GetPublicByAdminID)

	// Notification routes
	publicRoutes.Get("/notifications/admin/:adminID", notificationHandler.GetPublicByAdminID)

	// Restaurant routes
	publicRoutes.Get("/restaurants/admin/:adminID", restaurantHandler.GetPublicByAdminID)
}

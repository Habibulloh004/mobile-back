package routes

import (
	"mobilka/internal/api/handlers"
	"mobilka/internal/api/middlewares"

	"github.com/gofiber/fiber/v2"
)

// SetupSuperAdminRoutes sets up all routes related to super admin operations
func SetupSuperAdminRoutes(api fiber.Router, superAdminHandler *handlers.SuperAdminHandler, authHandler *handlers.AuthHandler) {
	// SuperAdmin routes
	superAdminRoutes := api.Group("/superadmin")
	superAdminRoutes.Use(middlewares.Protected(), middlewares.SuperAdminOnly())
	superAdminRoutes.Get("/profile", superAdminHandler.GetProfile)
	superAdminRoutes.Post("/change-password", authHandler.SuperAdminChangePassword)
}

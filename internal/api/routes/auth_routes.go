package routes

import (
	"mobilka/internal/api/handlers"

	"github.com/gofiber/fiber/v2"
)

// SetupAuthRoutes sets up all routes related to authentication
func SetupAuthRoutes(api fiber.Router, authHandler *handlers.AuthHandler) {
	// Auth routes (no auth required)
	auth := api.Group("/auth")
	auth.Post("/superadmin/login", authHandler.SuperAdminLogin)
	auth.Post("/admin/login", authHandler.AdminLogin)
}

package routes

import (
	"mobilka/internal/api/handlers"
	"mobilka/internal/api/middlewares"

	"github.com/gofiber/fiber/v2"
)

// SetupAdminRoutes sets up all routes related to admin operations
func SetupAdminRoutes(api fiber.Router, adminHandler *handlers.AdminHandler) {
	// Admin routes for super admin
	adminRoutes := api.Group("/admins")
	adminRoutes.Use(middlewares.Protected(), middlewares.SuperAdminOnly())
	adminRoutes.Post("/", adminHandler.Create)
	adminRoutes.Get("/", adminHandler.GetAll)
	adminRoutes.Get("/:id", adminHandler.GetByID)
	adminRoutes.Put("/:id", adminHandler.Update)
	adminRoutes.Delete("/:id", adminHandler.Delete)
	adminRoutes.Post("/:id/system-token", adminHandler.RegenerateSystemToken)
	adminRoutes.Post("/:id/sms-token", adminHandler.RegenerateSmsToken)

	// Admin profile route for regular admins
	adminProfileRoutes := api.Group("/admin")
	adminProfileRoutes.Use(middlewares.Protected(), middlewares.AdminOnly())
	adminProfileRoutes.Get("/profile", adminHandler.GetProfile)
}

package routes

import (
	"mobilka/internal/api/handlers"
	"mobilka/internal/api/middlewares"

	"github.com/gofiber/fiber/v2"
)

// SetupAdminRoutes sets up all routes related to admin operations
func SetupAdminRoutes(api fiber.Router, adminHandler *handlers.AdminHandler) {
	api.Get("/public/admins/:id", adminHandler.GetByIDPublic)
	api.Get("/public/mobileadmin/:id", adminHandler.GetByIDPublicMobile)
	// Admin routes for super admin
	adminRoutes := api.Group("/admins")
	adminRoutes.Use(middlewares.Protected(), middlewares.SuperAdminOnly())
	adminRoutes.Post("/", adminHandler.Create)
	adminRoutes.Get("/", adminHandler.GetAll)
	adminRoutes.Get("/:id", adminHandler.GetByID)
	adminRoutes.Put("/:id", adminHandler.Update)
	adminRoutes.Delete("/:id", adminHandler.Delete)

	// Admin profile route for regular admins
	adminProfileRoutes := api.Group("/admin")
	adminProfileRoutes.Use(middlewares.Protected(), middlewares.AdminOnly())
	adminProfileRoutes.Get("/profile", adminHandler.GetProfile)
}

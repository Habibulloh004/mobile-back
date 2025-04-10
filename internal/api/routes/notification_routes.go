package routes

import (
	"mobilka/internal/api/handlers"
	"mobilka/internal/api/middlewares"

	"github.com/gofiber/fiber/v2"
)

// SetupNotificationRoutes sets up all routes related to notification operations
func SetupNotificationRoutes(api fiber.Router, notificationHandler *handlers.NotificationHandler) {
	// Notification routes
	notificationRoutes := api.Group("/notifications")
	notificationRoutes.Use(middlewares.Protected())
	notificationRoutes.Post("/", notificationHandler.Create)
	notificationRoutes.Get("/", notificationHandler.GetAll)
	notificationRoutes.Get("/:id", notificationHandler.GetByID)
	notificationRoutes.Put("/:id", notificationHandler.Update)
	notificationRoutes.Delete("/:id", notificationHandler.Delete)
}

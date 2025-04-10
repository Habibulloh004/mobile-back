package routes

import (
	"mobilka/internal/api/handlers"
	"mobilka/internal/api/middlewares"

	"github.com/gofiber/fiber/v2"
)

// SetupFCMTokenRoutes sets up all routes related to FCM token operations
func SetupFCMTokenRoutes(api fiber.Router, fcmTokenHandler *handlers.FCMTokenHandler) {
	// FCM token routes
	fcmTokenRoutes := api.Group("/fcm-tokens")
	fcmTokenRoutes.Use(middlewares.Protected())
	fcmTokenRoutes.Post("/", fcmTokenHandler.Create)
	fcmTokenRoutes.Get("/", fcmTokenHandler.GetAll)
	fcmTokenRoutes.Delete("/:id", fcmTokenHandler.Delete)
	fcmTokenRoutes.Post("/delete-by-token", fcmTokenHandler.DeleteByToken)
}

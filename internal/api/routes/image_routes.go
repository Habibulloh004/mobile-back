package routes

import (
	"mobilka/internal/api/handlers"
	"mobilka/internal/api/middlewares"

	"github.com/gofiber/fiber/v2"
)

// SetupImageRoutes sets up all routes related to image operations
func SetupImageRoutes(app *fiber.App, api fiber.Router, imageHandler *handlers.ImageHandler) {
	// Protected image routes
	imageRoutes := api.Group("/images")
	imageRoutes.Use(middlewares.Protected())
	imageRoutes.Post("/", imageHandler.Upload)

	// Public image route (no auth required)
	app.Get("/uploads/images/:filename", imageHandler.Get)
}

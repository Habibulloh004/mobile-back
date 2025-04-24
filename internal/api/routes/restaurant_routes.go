package routes

import (
	"mobilka/internal/api/handlers"
	"mobilka/internal/api/middlewares"

	"github.com/gofiber/fiber/v2"
)

// SetupRestaurantRoutes sets up all routes related to restaurant operations
func SetupRestaurantRoutes(api fiber.Router, restaurantHandler *handlers.RestaurantHandler) {
	// Restaurant routes
	restaurantRoutes := api.Group("/restaurants")
	restaurantRoutes.Use(middlewares.Protected())
	restaurantRoutes.Post("/", restaurantHandler.Create)
	restaurantRoutes.Get("/", restaurantHandler.GetAll)
	restaurantRoutes.Get("/:id", restaurantHandler.GetByID)
	restaurantRoutes.Put("/:id", restaurantHandler.Update)
	restaurantRoutes.Delete("/:id", restaurantHandler.Delete)
}

// SetupPublicRestaurantRoutes sets up public routes for restaurants
func SetupPublicRestaurantRoutes(publicRoutes fiber.Router, restaurantHandler *handlers.RestaurantHandler) {
	// Public restaurant routes (no auth required)
	publicRoutes.Get("/restaurants/admin/:adminID", restaurantHandler.GetPublicByAdminID)
}

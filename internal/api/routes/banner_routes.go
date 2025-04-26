package routes

import (
	"mobilka/internal/api/handlers"
	"mobilka/internal/api/middlewares"

	"github.com/gofiber/fiber/v2"
)

// SetupBannerRoutes sets up all routes related to banner operations
func SetupBannerRoutes(api fiber.Router, bannerHandler *handlers.BannerHandler) {
	api.Get("/public/mobilebanner/:id", bannerHandler.GetByIDPublicMobile)
	// Banner routes
	bannerRoutes := api.Group("/banners")
	bannerRoutes.Use(middlewares.Protected())
	bannerRoutes.Post("/", bannerHandler.Create)
	bannerRoutes.Get("/", bannerHandler.GetAll)
	bannerRoutes.Get("/:id", bannerHandler.GetByID)
	bannerRoutes.Put("/:id", bannerHandler.Update)
	bannerRoutes.Delete("/:id", bannerHandler.Delete)
}

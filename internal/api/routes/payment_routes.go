package routes

import (
	"mobilka/internal/api/handlers"
	"mobilka/internal/api/middlewares"

	"github.com/gofiber/fiber/v2"
)

// SetupSubscriptionTierRoutes sets up all routes related to subscription tier operations
func SetupSubscriptionTierRoutes(api fiber.Router, subscriptionTierHandler *handlers.SubscriptionTierHandler) {
	// Subscription tier routes - super admin only
	subscriptionTierRoutes := api.Group("/subscription-tiers")
	subscriptionTierRoutes.Use(middlewares.Protected(), middlewares.SuperAdminOnly())
	subscriptionTierRoutes.Post("/", subscriptionTierHandler.Create)
	subscriptionTierRoutes.Get("/", subscriptionTierHandler.GetAll)
	subscriptionTierRoutes.Get("/:id", subscriptionTierHandler.GetByID)
	subscriptionTierRoutes.Put("/:id", subscriptionTierHandler.Update)
	subscriptionTierRoutes.Delete("/:id", subscriptionTierHandler.Delete)
}

// SetupPaymentRoutes sets up all routes related to payment operations
func SetupPaymentRoutes(api fiber.Router, paymentHandler *handlers.PaymentHandler, subscriptionTierHandler *handlers.SubscriptionTierHandler) {
	// Public subscription tier routes (for admins to see available tiers)
	api.Get("/public/subscription-tiers", subscriptionTierHandler.GetAll)

	// Admin payment routes
	adminPaymentRoutes := api.Group("/payments")
	adminPaymentRoutes.Use(middlewares.Protected(), middlewares.AdminOnly())
	adminPaymentRoutes.Post("/", paymentHandler.RecordPayment)
	adminPaymentRoutes.Get("/", paymentHandler.GetAdminPayments)
	adminPaymentRoutes.Get("/subscription", paymentHandler.GetSubscriptionInfo)

	// Super admin payment routes
	superadminPaymentRoutes := api.Group("/superadmin/payments")
	superadminPaymentRoutes.Use(middlewares.Protected(), middlewares.SuperAdminOnly())
	superadminPaymentRoutes.Get("/", paymentHandler.GetAllPayments)
	superadminPaymentRoutes.Get("/pending", paymentHandler.GetPendingPayments)
	superadminPaymentRoutes.Get("/:id", paymentHandler.GetPaymentByID)
	superadminPaymentRoutes.Post("/:id/verify", paymentHandler.VerifyPayment)
	superadminPaymentRoutes.Get("/admin/:id/subscription", paymentHandler.GetSubscriptionInfo)
}

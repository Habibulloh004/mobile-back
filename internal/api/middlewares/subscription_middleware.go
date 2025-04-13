package middlewares

import (
	"mobilka/internal/service"
	"mobilka/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// SubscriptionChecker creates middleware to check subscription status
func SubscriptionChecker(paymentService *service.PaymentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip check for super admin
		role, roleOk := c.Locals(utils.ContextUserRole).(string)
		if roleOk && role == utils.RoleSuperAdmin {
			return c.Next()
		}

		// Skip check for public endpoints and auth endpoints
		path := c.Path()
		if isPublicPath(path) || isAuthPath(path) {
			return c.Next()
		}

		// Check if user is authenticated
		adminID, ok := c.Locals(utils.ContextUserID).(int)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Unauthorized",
			})
		}

		// Check subscription access
		hasAccess, err := paymentService.CheckAdminAccess(c.Context(), adminID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Failed to check subscription status",
			})
		}

		// If access is restricted, return subscription needed error
		if !hasAccess {
			return c.Status(fiber.StatusPaymentRequired).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Subscription payment required. Please make a payment to continue using the service.",
				"code":    "SUBSCRIPTION_REQUIRED",
			})
		}

		// Continue to the next middleware or handler
		return c.Next()
	}
}

// isPublicPath checks if a path is a public endpoint
func isPublicPath(path string) bool {
	// Add all public paths here
	publicPaths := []string{
		"/api/public/",
		"/api/uploads/images/",
	}

	for _, p := range publicPaths {
		if len(path) >= len(p) && path[:len(p)] == p {
			return true
		}
	}

	return false
}

// isAuthPath checks if a path is an auth endpoint
func isAuthPath(path string) bool {
	// Add all auth paths here
	authPaths := []string{
		"/api/auth/",
		"/api/payments/subscription", // Allow checking subscription status
	}

	for _, p := range authPaths {
		if len(path) >= len(p) && path[:len(p)] == p {
			return true
		}
	}

	return false
}

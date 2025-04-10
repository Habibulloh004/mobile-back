package middlewares

import (
	"strings"

	"mobilka/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// Protected middleware ensures that the request is authenticated with a valid JWT
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the Authorization header
		authHeader := c.Get("Authorization")

		// Check if the Authorization header is empty
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Unauthorized: Missing authorization header",
			})
		}

		// Extract the token from the Authorization header
		// Format: Bearer <token>
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Unauthorized: Invalid authorization format",
			})
		}

		// Parse and validate the token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Unauthorized: " + err.Error(),
			})
		}

		// Store user ID and role in context for use in handlers
		c.Locals(utils.ContextUserID, claims.ID)
		c.Locals(utils.ContextUserRole, claims.Role)

		// Continue to the next middleware or handler
		return c.Next()
	}
}

// AdminOnly middleware ensures that the request is from an admin
func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the user role from context (set by Protected middleware)
		role, ok := c.Locals(utils.ContextUserRole).(string)
		if !ok || role != utils.RoleAdmin {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Forbidden: Admin access required",
			})
		}

		// Continue to the next middleware or handler
		return c.Next()
	}
}

// SuperAdminOnly middleware ensures that the request is from a super admin
func SuperAdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the user role from context (set by Protected middleware)
		role, ok := c.Locals(utils.ContextUserRole).(string)
		if !ok || role != utils.RoleSuperAdmin {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Forbidden: Super admin access required",
			})
		}

		// Continue to the next middleware or handler
		return c.Next()
	}
}

// GetUserID gets the user ID from the context
func GetUserID(c *fiber.Ctx) (int, bool) {
	id, ok := c.Locals(utils.ContextUserID).(int)
	return id, ok
}

// GetUserRole gets the user role from the context
func GetUserRole(c *fiber.Ctx) (string, bool) {
	role, ok := c.Locals(utils.ContextUserRole).(string)
	return role, ok
}

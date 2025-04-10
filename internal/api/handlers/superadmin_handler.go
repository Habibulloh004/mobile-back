package handlers

import (
	"mobilka/internal/service"
	"mobilka/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// SuperAdminHandler handles super admin requests
type SuperAdminHandler struct {
	superAdminService *service.SuperAdminService
}

// NewSuperAdminHandler creates a new super admin handler
func NewSuperAdminHandler(superAdminService *service.SuperAdminService) *SuperAdminHandler {
	return &SuperAdminHandler{
		superAdminService: superAdminService,
	}
}

// GetProfile handles retrieving the super admin profile
func (h *SuperAdminHandler) GetProfile(c *fiber.Ctx) error {
	// Get super admin ID from context
	userID, ok := c.Locals(utils.ContextUserID).(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Unauthorized",
		})
	}

	// Get super admin profile
	superAdmin, err := h.superAdminService.GetByID(c.Context(), userID)
	if err != nil {
		if err == utils.ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Super admin not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to retrieve profile",
		})
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": utils.StatusSuccess,
		"data":   superAdmin.ToResponse(),
	})
}

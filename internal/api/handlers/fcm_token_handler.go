package handlers

import (
	"strconv"

	"mobilka/internal/models"
	"mobilka/internal/service"
	"mobilka/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// FCMTokenHandler handles FCM token requests
type FCMTokenHandler struct {
	fcmTokenService *service.FCMTokenService
}

// NewFCMTokenHandler creates a new FCM token handler
func NewFCMTokenHandler(fcmTokenService *service.FCMTokenService) *FCMTokenHandler {
	return &FCMTokenHandler{
		fcmTokenService: fcmTokenService,
	}
}

// Create handles creating a new FCM token
func (h *FCMTokenHandler) Create(c *fiber.Ctx) error {
	// Get admin ID from context
	adminID, ok := c.Locals(utils.ContextUserID).(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Unauthorized",
		})
	}

	var req models.FCMTokenCreateRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid request body",
		})
	}

	// Create FCM token
	fcmToken, err := h.fcmTokenService.Create(c.Context(), adminID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to create FCM token",
		})
	}

	// Return response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": utils.StatusSuccess,
		"data":   fcmToken.ToResponse(),
	})
}

// GetAll handles retrieving all FCM tokens for the current admin
func (h *FCMTokenHandler) GetAll(c *fiber.Ctx) error {
	// Get admin ID from context
	adminID, ok := c.Locals(utils.ContextUserID).(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Unauthorized",
		})
	}

	// Get role from context
	role, _ := c.Locals(utils.ContextUserRole).(string)

	var fcmTokens []*models.FCMToken
	var err error

	// Super admin can see all FCM tokens, admin can only see their own
	if role == utils.RoleSuperAdmin {
		fcmTokens, err = h.fcmTokenService.GetAll(c.Context())
	} else {
		fcmTokens, err = h.fcmTokenService.GetByAdminID(c.Context(), adminID)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to retrieve FCM tokens",
		})
	}

	// Convert to response objects
	var responses []models.FCMTokenResponse
	for _, fcmToken := range fcmTokens {
		responses = append(responses, fcmToken.ToResponse())
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": utils.StatusSuccess,
		"data":   responses,
	})
}

// Delete handles deleting an FCM token
func (h *FCMTokenHandler) Delete(c *fiber.Ctx) error {
	// Get FCM token ID from URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid FCM token ID",
		})
	}

	// Delete FCM token
	err = h.fcmTokenService.Delete(c.Context(), id)
	if err != nil {
		if err == utils.ErrResourceNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "FCM token not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to delete FCM token",
		})
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  utils.StatusSuccess,
		"message": "FCM token deleted successfully",
	})
}

// DeleteByToken handles deleting an FCM token by token string
func (h *FCMTokenHandler) DeleteByToken(c *fiber.Ctx) error {
	// Get FCM token from request body
	var req struct {
		Token string `json:"token" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid request body",
		})
	}

	// Delete FCM token
	err := h.fcmTokenService.DeleteByToken(c.Context(), req.Token)
	if err != nil {
		if err == utils.ErrResourceNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "FCM token not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to delete FCM token",
		})
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  utils.StatusSuccess,
		"message": "FCM token deleted successfully",
	})
}

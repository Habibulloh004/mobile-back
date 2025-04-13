package handlers

import (
	"fmt"
	"strconv"

	"mobilka/internal/models"
	"mobilka/internal/service"
	"mobilka/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// NotificationHandler handles notification requests
type NotificationHandler struct {
	notificationService *service.NotificationService
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(notificationService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// Create handles creating a new notification
func (h *NotificationHandler) Create(c *fiber.Ctx) error {
	// Get admin ID from context
	contextAdminID, ok := c.Locals(utils.ContextUserID).(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Unauthorized",
		})
	}

	// Get role from context
	role, _ := c.Locals(utils.ContextUserRole).(string)

	var req models.NotificationCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid request body",
		})
	}

	// Validate required fields
	if req.Payload == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Payload is required",
		})
	}

	if req.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Title is required",
		})
	}

	if req.Body == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Body is required",
		})
	}

	// Determine which admin ID to use
	adminID := contextAdminID

	// For notification creation, we need to update the model to support admin_id in the request
	if req.AdminID > 0 && role == utils.RoleSuperAdmin {
		// Only super admins can create notifications for other admins
		adminID = req.AdminID
		fmt.Printf("Super admin creating notification for admin ID: %d\n", adminID)
	} else {
		fmt.Printf("Regular admin creating notification with own ID: %d\n", adminID)
	}

	// Create notification
	notification, err := h.notificationService.Create(c.Context(), adminID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to create notification",
		})
	}

	// Return response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": utils.StatusSuccess,
		"data":   notification.ToResponse(),
	})
}

// GetAll handles retrieving all notifications for the current admin
func (h *NotificationHandler) GetAll(c *fiber.Ctx) error {
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

	var notifications []*models.Notification
	var err error

	// Super admin can see all notifications, admin can only see their own
	if role == utils.RoleSuperAdmin {
		notifications, err = h.notificationService.GetAll(c.Context())
	} else {
		notifications, err = h.notificationService.GetByAdminID(c.Context(), adminID)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to retrieve notifications",
		})
	}

	// Convert to response objects
	var responses []models.NotificationResponse
	for _, notification := range notifications {
		responses = append(responses, notification.ToResponse())
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": utils.StatusSuccess,
		"data":   responses,
	})
}

// GetByID handles retrieving a notification by ID
func (h *NotificationHandler) GetByID(c *fiber.Ctx) error {
	// Get notification ID from URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid notification ID",
		})
	}

	// Get notification
	notification, err := h.notificationService.GetByID(c.Context(), id)
	if err != nil {
		if err == utils.ErrResourceNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Notification not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to retrieve notification",
		})
	}

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

	// Check if the user has access to this notification
	if role != utils.RoleSuperAdmin && notification.AdminID != adminID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Access denied",
		})
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": utils.StatusSuccess,
		"data":   notification.ToResponse(),
	})
}

// Update handles updating a notification
func (h *NotificationHandler) Update(c *fiber.Ctx) error {
	// Get notification ID from URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid notification ID",
		})
	}

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

	var req models.NotificationUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid request body",
		})
	}

	// First, get the existing notification to check ownership
	existingNotification, err := h.notificationService.GetByID(c.Context(), id)
	if err != nil {
		if err == utils.ErrResourceNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Notification not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to retrieve notification",
		})
	}

	// Check if user has permission to update this notification
	// Super admins can update any notification
	// Regular admins can only update their own notifications
	if role != utils.RoleSuperAdmin && existingNotification.AdminID != adminID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "You don't have permission to update this notification",
		})
	}

	// Allow changing adminID only for super admins
	targetAdminID := existingNotification.AdminID
	if req.AdminID > 0 && role == utils.RoleSuperAdmin {
		targetAdminID = req.AdminID
		fmt.Printf("Super admin changing notification admin ID from %d to %d\n", existingNotification.AdminID, targetAdminID)
	}

	// Update notification
	notification, err := h.notificationService.Update(c.Context(), id, targetAdminID, &req)
	if err != nil {
		if err == utils.ErrResourceNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Notification not found or access denied",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to update notification",
		})
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": utils.StatusSuccess,
		"data":   notification.ToResponse(),
	})
}

// Delete handles deleting a notification
func (h *NotificationHandler) Delete(c *fiber.Ctx) error {
	// Get notification ID from URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid notification ID",
		})
	}

	// Get admin ID from context
	adminID, ok := c.Locals(utils.ContextUserID).(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Unauthorized",
		})
	}

	// Delete notification
	err = h.notificationService.Delete(c.Context(), id, adminID)
	if err != nil {
		if err == utils.ErrResourceNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Notification not found or access denied",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to delete notification",
		})
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  utils.StatusSuccess,
		"message": "Notification deleted successfully",
	})
}

// GetPublicByAdminID handles retrieving notifications for a specific admin with pagination
func (h *NotificationHandler) GetPublicByAdminID(c *fiber.Ctx) error {
	// Get admin ID from URL
	adminID, err := strconv.Atoi(c.Params("adminID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid admin ID",
		})
	}

	// Parse pagination parameters
	skip, err := strconv.Atoi(c.Query("skip", "0"))
	if err != nil || skip < 0 {
		skip = 0
	}

	step, err := strconv.Atoi(c.Query("step", "10"))
	if err != nil || step <= 0 || step > 100 {
		step = 10 // Default limit is 10, max is 100
	}

	// Get notifications with pagination
	notifications, err := h.notificationService.GetByAdminIDWithPagination(c.Context(), adminID, skip, step)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to retrieve notifications",
		})
	}

	// Convert to response objects
	var responses []models.NotificationResponse
	for _, notification := range notifications {
		responses = append(responses, notification.ToResponse())
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": utils.StatusSuccess,
		"data":   responses,
		"meta": fiber.Map{
			"skip": skip,
			"step": step,
		},
	})
}

package handlers

import (
	"strconv"

	"mobilka/internal/models"
	"mobilka/internal/service"
	"mobilka/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// BannerHandler handles banner requests
type BannerHandler struct {
	bannerService *service.BannerService
}

// GetByID handles retrieving a banner by ID
func (h *BannerHandler) GetByID(c *fiber.Ctx) error {
	// Get banner ID from URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid banner ID",
		})
	}

	// Get banner
	banner, err := h.bannerService.GetByID(c.Context(), id)
	if err != nil {
		if err == utils.ErrResourceNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Banner not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to retrieve banner",
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

	// Check if the user has access to this banner
	if role != utils.RoleSuperAdmin && banner.AdminID != adminID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Access denied",
		})
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": utils.StatusSuccess,
		"data":   banner.ToResponse(),
	})
}

// Update handles updating a banner
func (h *BannerHandler) Update(c *fiber.Ctx) error {
	// Get banner ID from URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid banner ID",
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

	var req models.BannerUpdateRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid request body",
		})
	}

	// Update banner
	banner, err := h.bannerService.Update(c.Context(), id, adminID, &req)
	if err != nil {
		if err == utils.ErrResourceNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Banner not found or access denied",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to update banner",
		})
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": utils.StatusSuccess,
		"data":   banner.ToResponse(),
	})
}

// Delete handles deleting a banner
func (h *BannerHandler) Delete(c *fiber.Ctx) error {
	// Get banner ID from URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid banner ID",
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

	// Delete banner
	err = h.bannerService.Delete(c.Context(), id, adminID)
	if err != nil {
		if err == utils.ErrResourceNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Banner not found or access denied",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to delete banner",
		})
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  utils.StatusSuccess,
		"message": "Banner deleted successfully",
	})
}

// NewBannerHandler creates a new banner handler
func NewBannerHandler(bannerService *service.BannerService) *BannerHandler {
	return &BannerHandler{
		bannerService: bannerService,
	}
}

// Create handles creating a new banner
func (h *BannerHandler) Create(c *fiber.Ctx) error {
	// Get admin ID from context
	adminID, ok := c.Locals(utils.ContextUserID).(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Unauthorized",
		})
	}

	var req models.BannerCreateRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid request body",
		})
	}

	// Create banner
	banner, err := h.bannerService.Create(c.Context(), adminID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to create banner",
		})
	}

	// Return response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": utils.StatusSuccess,
		"data":   banner.ToResponse(),
	})
}

// GetAll handles retrieving all banners for the current admin
func (h *BannerHandler) GetAll(c *fiber.Ctx) error {
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

	var banners []*models.Banner
	var err error

	// Super admin can see all banners, admin can only see their own
	if role == utils.RoleSuperAdmin {
		banners, err = h.bannerService.GetAll(c.Context())
	} else {
		banners, err = h.bannerService.GetByAdminID(c.Context(), adminID)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to retrieve banners",
		})
	}

	// Convert to response objects
	var responses []models.BannerResponse
	for _, banner := range banners {
		responses = append(responses, banner.ToResponse())
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": utils.StatusSuccess,
		"data":   responses,
	})
}

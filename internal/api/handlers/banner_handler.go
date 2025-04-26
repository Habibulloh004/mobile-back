package handlers

import (
	"fmt"
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

    // Get role from context
    role, _ := c.Locals(utils.ContextUserRole).(string)

    var req models.BannerUpdateRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  utils.StatusError,
            "message": "Invalid request body",
        })
    }

    // First, get the existing banner to check ownership
    existingBanner, err := h.bannerService.GetByID(c.Context(), id)
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

    // Check if user has permission to update this banner
    // Super admins can update any banner
    // Regular admins can only update their own banners
    if role != utils.RoleSuperAdmin && existingBanner.AdminID != adminID {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "status":  utils.StatusError,
            "message": "You don't have permission to update this banner",
        })
    }

    // Allow changing adminID only for super admins
    targetAdminID := existingBanner.AdminID
    if req.AdminID > 0 && role == utils.RoleSuperAdmin {
        targetAdminID = req.AdminID
        fmt.Printf("Super admin changing banner admin ID from %d to %d\n", existingBanner.AdminID, targetAdminID)
    }

    // Update banner
    banner, err := h.bannerService.Update(c.Context(), id, targetAdminID, &req)
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
	contextAdminID, ok := c.Locals(utils.ContextUserID).(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Unauthorized",
		})
	}

	// Get role from context
	role, _ := c.Locals(utils.ContextUserRole).(string)
	
	var req models.BannerCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid request body",
		})
	}

	// Validate required fields
	if req.Image == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Image is required",
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
	
	// For banner creation, we need to update the model to support admin_id in the request
	// Add the AdminID field to BannerCreateRequest in internal/models/banner.go
	if req.AdminID > 0 && role == utils.RoleSuperAdmin {
		// Only super admins can create banners for other admins
		adminID = req.AdminID
		fmt.Printf("Super admin creating banner for admin ID: %d\n", adminID)
	} else {
		fmt.Printf("Regular admin creating banner with own ID: %d\n", adminID)
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

func (h *BannerHandler) GetByIDPublicMobile(c *fiber.Ctx) error {
	// Get admin ID from context
	adminID, _ := strconv.Atoi(c.Params("id"))
	fmt.Print(adminID)


	var banners []*models.Banner
	var err error

	banners, err = h.bannerService.GetByAdminID(c.Context(), adminID)

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

// GetPublicByAdminID handles retrieving all banners for a specific admin without authentication

func (h *BannerHandler) GetPublicByAdminID(c *fiber.Ctx) error {
	// Get admin ID from URL
	adminID, err := strconv.Atoi(c.Params("adminID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid admin ID",
		})
	}

	// Get banners
	banners, err := h.bannerService.GetByAdminID(c.Context(), adminID)
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

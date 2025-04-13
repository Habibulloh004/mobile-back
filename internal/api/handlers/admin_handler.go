package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"mobilka/internal/models"
	"mobilka/internal/service"
	"mobilka/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// AdminHandler handles admin requests
type AdminHandler struct {
	adminService *service.AdminService
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(adminService *service.AdminService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
	}
}

// Create handles creating a new admin
func (h *AdminHandler) Create(c *fiber.Ctx) error {
	var req models.AdminCreateRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid request body",
		})
	}

	// Validate required fields
	if req.UserName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Username is required",
		})
	}

	if req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Email is required",
		})
	}

	if req.CompanyName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Company name is required",
		})
	}

	if req.SystemID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "System id is required",
		})
	}

	// Create admin
	admin, err := h.adminService.Create(c.Context(), &req)
	if err != nil {
		// Check if it's a detailed app error
		var appErr *utils.AppError
		if errors.As(err, &appErr) {
			return c.Status(appErr.Code).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": appErr.Message,
			})
		}

		// Default error response
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to create admin: " + err.Error(),
		})
	}

	// Return response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": utils.StatusSuccess,
		"data":   admin.ToResponse(),
	})
}

// GetAll handles retrieving all admins
func (h *AdminHandler) GetAll(c *fiber.Ctx) error {
	// Get all admins
	admins, err := h.adminService.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to retrieve admins",
		})
	}

	// Convert to response objects
	var responses []models.AdminResponse
	for _, admin := range admins {
		responses = append(responses, admin.ToResponse())
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": utils.StatusSuccess,
		"data":   responses,
	})
}

// GetByID handles retrieving an admin by ID
func (h *AdminHandler) GetByID(c *fiber.Ctx) error {
	// Get admin ID from URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid admin ID",
		})
	}

	// Get admin
	admin, err := h.adminService.GetByID(c.Context(), id)
	if err != nil {
		if err == utils.ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Admin not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to retrieve admin",
		})
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": utils.StatusSuccess,
		"data":   admin.ToResponse(),
	})
}

// GetProfile handles retrieving the current admin's profile
func (h *AdminHandler) GetProfile(c *fiber.Ctx) error {
	// Get admin ID from context
	userID, ok := c.Locals(utils.ContextUserID).(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Unauthorized",
		})
	}

	// Get admin
	admin, err := h.adminService.GetByID(c.Context(), userID)
	if err != nil {
		if err == utils.ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Admin not found",
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
		"data":   admin.ToResponse(),
	})
}

// Update handles updating an admin
func (h *AdminHandler) Update(c *fiber.Ctx) error {
    // Get admin ID from URL
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  utils.StatusError,
            "message": "Invalid admin ID",
        })
    }

    // Log the raw request body for debugging
    rawBody := c.Body()
    fmt.Printf("Raw request body for admin ID %d update: %s\n", id, string(rawBody))

    var req models.AdminUpdateRequest
    if err := c.BodyParser(&req); err != nil {
        fmt.Printf("Error parsing request body: %v\n", err)
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  utils.StatusError,
            "message": "Invalid request body: " + err.Error(),
        })
    }

    // Log parsed request for debugging
    fmt.Printf("Parsed request for admin ID %d update: %+v\n", id, req)

    // Check if system_token is in the request (even if parsing didn't set it)
    var requestMap map[string]interface{}
    if err := json.Unmarshal(rawBody, &requestMap); err == nil {
        if systemToken, exists := requestMap["system_token"]; exists {
            fmt.Printf("Found system_token in raw map: %v\n", systemToken)
            if req.SystemToken == "" {
                req.SystemToken = fmt.Sprintf("%v", systemToken)
                fmt.Printf("Setting req.SystemToken to: %s\n", req.SystemToken)
            }
        }
    }

    // Update admin
    admin, err := h.adminService.Update(c.Context(), id, &req)
    if err != nil {
        if err == utils.ErrUserNotFound {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "status":  utils.StatusError,
                "message": "Admin not found",
            })
        }

        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "status":  utils.StatusError,
            "message": "Failed to update admin: " + err.Error(),
        })
    }

    // Return response
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "status": utils.StatusSuccess,
        "data":   admin.ToResponse(),
    })
}

// Delete handles deleting an admin
func (h *AdminHandler) Delete(c *fiber.Ctx) error {
	// Get admin ID from URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid admin ID",
		})
	}

	// Delete admin
	err = h.adminService.Delete(c.Context(), id)
	if err != nil {
		if err == utils.ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Admin not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to delete admin",
		})
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  utils.StatusSuccess,
		"message": "Admin deleted successfully",
	})
}

// GetByIDPublic handles retrieving an admin by ID without authentication
func (h *AdminHandler) GetByIDPublic(c *fiber.Ctx) error {
	// Get admin ID from URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid admin ID",
		})
	}

	// Get admin and increment users count
	admin, err := h.adminService.GetByIDPublic(c.Context(), id)
	if err != nil {
		if err == utils.ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Admin not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to retrieve admin",
		})
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": utils.StatusSuccess,
		"data":   admin.ToResponse(),
	})
}

package handlers

import (
	"fmt"
	"strconv"

	"mobilka/internal/models"
	"mobilka/internal/service"
	"mobilka/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// RestaurantHandler handles restaurant requests
type RestaurantHandler struct {
	restaurantService *service.RestaurantService
}

// NewRestaurantHandler creates a new restaurant handler
func NewRestaurantHandler(restaurantService *service.RestaurantService) *RestaurantHandler {
	return &RestaurantHandler{
		restaurantService: restaurantService,
	}
}

// Create handles creating a new restaurant or updating if one already exists
func (h *RestaurantHandler) Create(c *fiber.Ctx) error {
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

	var req models.RestaurantCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid request body",
		})
	}

	// Validate required fields
	if req.Text == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Text is required",
		})
	}

	// Determine which admin ID to use
	adminID := contextAdminID

	// For restaurant creation, super admin can create restaurants for other admins
	if req.AdminID > 0 && role == utils.RoleSuperAdmin {
		// Only super admins can create restaurants for other admins
		adminID = req.AdminID
		fmt.Printf("Super admin creating/updating restaurant for admin ID: %d\n", adminID)
	} else {
		fmt.Printf("Regular admin creating/updating restaurant with own ID: %d\n", adminID)
	}

	// Create or update restaurant
	restaurant, err := h.restaurantService.Create(c.Context(), adminID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to create/update restaurant",
		})
	}

	// Return response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  utils.StatusSuccess,
		"message": "Restaurant created/updated successfully",
		"data":    restaurant.ToResponse(),
	})
}

// GetAll handles retrieving all restaurants for the current admin
func (h *RestaurantHandler) GetAll(c *fiber.Ctx) error {
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

	var restaurants []*models.Restaurant
	var err error

	// Super admin can see all restaurants, admin can only see their own
	if role == utils.RoleSuperAdmin {
		restaurants, err = h.restaurantService.GetAll(c.Context())
	} else {
		restaurants, err = h.restaurantService.GetByAdminID(c.Context(), adminID)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to retrieve restaurants",
		})
	}

	// Convert to response objects
	var responses []models.RestaurantResponse
	for _, restaurant := range restaurants {
		responses = append(responses, restaurant.ToResponse())
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": utils.StatusSuccess,
		"data":   responses,
	})
}

// GetByID handles retrieving a restaurant by ID
func (h *RestaurantHandler) GetByID(c *fiber.Ctx) error {
	// Get restaurant ID from URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid restaurant ID",
		})
	}

	// Get restaurant
	restaurant, err := h.restaurantService.GetByID(c.Context(), id)
	if err != nil {
		if err == utils.ErrResourceNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Restaurant not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to retrieve restaurant",
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

	// Check if the user has access to this restaurant
	if role != utils.RoleSuperAdmin && restaurant.AdminID != adminID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Access denied",
		})
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": utils.StatusSuccess,
		"data":   restaurant.ToResponse(),
	})
}

// Update handles updating a restaurant
func (h *RestaurantHandler) Update(c *fiber.Ctx) error {
	// Get restaurant ID from URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid restaurant ID",
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

	var req models.RestaurantUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid request body",
		})
	}

	// First, get the existing restaurant to check ownership
	existingRestaurant, err := h.restaurantService.GetByID(c.Context(), id)
	if err != nil {
		if err == utils.ErrResourceNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Restaurant not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to retrieve restaurant",
		})
	}

	// Check if user has permission to update this restaurant
	// Super admins can update any restaurant
	// Regular admins can only update their own restaurants
	if role != utils.RoleSuperAdmin && existingRestaurant.AdminID != adminID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "You don't have permission to update this restaurant",
		})
	}

	// Allow changing adminID only for super admins
	targetAdminID := existingRestaurant.AdminID
	if req.AdminID > 0 && role == utils.RoleSuperAdmin {
		targetAdminID = req.AdminID
		fmt.Printf("Super admin changing restaurant admin ID from %d to %d\n", existingRestaurant.AdminID, targetAdminID)
	}

	// Update restaurant
	restaurant, err := h.restaurantService.Update(c.Context(), id, targetAdminID, &req)
	if err != nil {
		if err == utils.ErrResourceNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Restaurant not found or access denied",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to update restaurant",
		})
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": utils.StatusSuccess,
		"data":   restaurant.ToResponse(),
	})
}

// Delete handles deleting a restaurant
func (h *RestaurantHandler) Delete(c *fiber.Ctx) error {
	// Get restaurant ID from URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid restaurant ID",
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

	// Delete restaurant
	err = h.restaurantService.Delete(c.Context(), id, adminID)
	if err != nil {
		if err == utils.ErrResourceNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Restaurant not found or access denied",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to delete restaurant",
		})
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  utils.StatusSuccess,
		"message": "Restaurant deleted successfully",
	})
}

// GetPublicByAdminID handles retrieving all restaurants for a specific admin without authentication
func (h *RestaurantHandler) GetPublicByAdminID(c *fiber.Ctx) error {
	// Get admin ID from URL
	adminID, err := strconv.Atoi(c.Params("adminID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid admin ID",
		})
	}

	// Get restaurants
	restaurants, err := h.restaurantService.GetByAdminID(c.Context(), adminID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to retrieve restaurants",
		})
	}

	// Convert to response objects
	var responses []models.RestaurantResponse
	for _, restaurant := range restaurants {
		responses = append(responses, restaurant.ToResponse())
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": utils.StatusSuccess,
		"data":   responses,
	})
}

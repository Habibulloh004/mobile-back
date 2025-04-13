package handlers

import (
	"log"
	"strconv"
	"strings"

	"mobilka/internal/models"
	"mobilka/internal/service"
	"mobilka/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// SubscriptionTierHandler handles subscription tier requests
type SubscriptionTierHandler struct {
	subscriptionTierService *service.SubscriptionTierService
}

// NewSubscriptionTierHandler creates a new subscription tier handler
func NewSubscriptionTierHandler(subscriptionTierService *service.SubscriptionTierService) *SubscriptionTierHandler {
	return &SubscriptionTierHandler{
		subscriptionTierService: subscriptionTierService,
	}
}

// Create handles creating a new subscription tier
func (h *SubscriptionTierHandler) Create(c *fiber.Ctx) error {
	var req models.SubscriptionTierCreateRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid request body",
		})
	}

	// Validate required fields
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Name is required",
		})
	}

	// Create subscription tier
	tier, err := h.subscriptionTierService.Create(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to create subscription tier: " + err.Error(),
		})
	}

	// Return response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": utils.StatusSuccess,
		"data":   tier.ToResponse(),
	})
}

// GetAll handles retrieving all subscription tiers
func (h *SubscriptionTierHandler) GetAll(c *fiber.Ctx) error {
	// Get all subscription tiers

	tiers, err := h.subscriptionTierService.GetAll(c.Context())
	if err != nil {
		// Log the detailed error for debugging
		log.Printf("Error retrieving subscription tiers: %v", err)

		// Check if it's a database connection or schema error
		if strings.Contains(err.Error(), "does not exist") {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Subscription system is being initialized. Please try again later.",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to retrieve subscription tiers",
		})
	}

	// If no tiers exist yet, return an empty array instead of null
	if tiers == nil {
		tiers = []*models.SubscriptionTier{}
	}

	// Convert to response objects
	var responses []models.SubscriptionTierResponse
	for _, tier := range tiers {
		responses = append(responses, tier.ToResponse())
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": utils.StatusSuccess,
		"data":   responses,
	})
}

// GetByID handles retrieving a subscription tier by ID
func (h *SubscriptionTierHandler) GetByID(c *fiber.Ctx) error {
	// Get subscription tier ID from URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid subscription tier ID",
		})
	}

	// Get subscription tier
	tier, err := h.subscriptionTierService.GetByID(c.Context(), id)
	if err != nil {
		if err == utils.ErrResourceNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Subscription tier not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to retrieve subscription tier",
		})
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": utils.StatusSuccess,
		"data":   tier.ToResponse(),
	})
}

// Update handles updating a subscription tier
func (h *SubscriptionTierHandler) Update(c *fiber.Ctx) error {
	// Get subscription tier ID from URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid subscription tier ID",
		})
	}

	var req models.SubscriptionTierUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid request body",
		})
	}

	// Update subscription tier
	tier, err := h.subscriptionTierService.Update(c.Context(), id, &req)
	if err != nil {
		if err == utils.ErrResourceNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Subscription tier not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to update subscription tier",
		})
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": utils.StatusSuccess,
		"data":   tier.ToResponse(),
	})
}

// Delete handles deleting a subscription tier
func (h *SubscriptionTierHandler) Delete(c *fiber.Ctx) error {
	// Get subscription tier ID from URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid subscription tier ID",
		})
	}

	// Delete subscription tier
	err = h.subscriptionTierService.Delete(c.Context(), id)
	if err != nil {
		if err == utils.ErrResourceNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Subscription tier not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to delete subscription tier",
		})
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  utils.StatusSuccess,
		"message": "Subscription tier deleted successfully",
	})
}

package handlers

import (
	"strconv"

	"mobilka/internal/models"
	"mobilka/internal/service"
	"mobilka/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// PaymentHandler handles payment requests
type PaymentHandler struct {
	paymentService *service.PaymentService
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(paymentService *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// RecordPayment handles recording a new payment
func (h *PaymentHandler) RecordPayment(c *fiber.Ctx) error {
	// Get admin ID from context
	adminID, ok := c.Locals(utils.ContextUserID).(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Unauthorized",
		})
	}

	var req models.PaymentCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid request body",
		})
	}

	// Validate required fields
	if req.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Amount must be greater than zero",
		})
	}

	if req.PaymentMethod == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Payment method is required",
		})
	}

	// Record payment
	payment, err := h.paymentService.RecordPayment(c.Context(), adminID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to record payment: " + err.Error(),
		})
	}

	// Return response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  utils.StatusSuccess,
		"data":    payment.ToResponse(),
		"message": "Payment recorded successfully. It will be reviewed by an administrator.",
	})
}

// GetAdminPayments handles retrieving all payments for the current admin
func (h *PaymentHandler) GetAdminPayments(c *fiber.Ctx) error {
	// Get admin ID from context
	adminID, ok := c.Locals(utils.ContextUserID).(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Unauthorized",
		})
	}

	// Get admin payments
	payments, err := h.paymentService.GetPaymentsByAdminID(c.Context(), adminID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to retrieve payments",
		})
	}

	// Convert to response objects
	var responses []models.PaymentHistoryResponse
	for _, payment := range payments {
		responses = append(responses, payment.ToResponse())
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": utils.StatusSuccess,
		"data":   responses,
	})
}

// GetAllPayments handles retrieving all payments (super admin only)
func (h *PaymentHandler) GetAllPayments(c *fiber.Ctx) error {
	// Get all payments
	payments, err := h.paymentService.GetAllPayments(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to retrieve payments",
		})
	}

	// Convert to response objects
	var responses []models.PaymentHistoryResponse
	for _, payment := range payments {
		responses = append(responses, payment.ToResponse())
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": utils.StatusSuccess,
		"data":   responses,
	})
}

// GetPendingPayments handles retrieving all pending payments (super admin only)
func (h *PaymentHandler) GetPendingPayments(c *fiber.Ctx) error {
	// Get pending payments
	payments, err := h.paymentService.GetPendingPayments(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to retrieve pending payments",
		})
	}

	// Convert to response objects
	var responses []models.PaymentHistoryResponse
	for _, payment := range payments {
		responses = append(responses, payment.ToResponse())
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": utils.StatusSuccess,
		"data":   responses,
	})
}

// GetPaymentByID handles retrieving a payment by ID
func (h *PaymentHandler) GetPaymentByID(c *fiber.Ctx) error {
	// Get payment ID from URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid payment ID",
		})
	}

	// Get payment
	payment, err := h.paymentService.GetPaymentByID(c.Context(), id)
	if err != nil {
		if err == utils.ErrResourceNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Payment not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to retrieve payment",
		})
	}

	// Check if the request is from super admin or the admin who made the payment
	role, _ := c.Locals(utils.ContextUserRole).(string)
	userID, _ := c.Locals(utils.ContextUserID).(int)

	if role != utils.RoleSuperAdmin && payment.AdminID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Access denied",
		})
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": utils.StatusSuccess,
		"data":   payment.ToResponse(),
	})
}

// VerifyPayment handles verifying a payment (super admin only)
func (h *PaymentHandler) VerifyPayment(c *fiber.Ctx) error {
	// Get super admin ID from context
	superAdminID, ok := c.Locals(utils.ContextUserID).(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Unauthorized",
		})
	}

	// Get payment ID from URL
	paymentID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid payment ID",
		})
	}

	var req models.PaymentVerifyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid request body",
		})
	}

	// Validate status
	if req.Status != "verified" && req.Status != "rejected" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Status must be 'verified' or 'rejected'",
		})
	}

	// Verify payment
	err = h.paymentService.VerifyPayment(c.Context(), paymentID, superAdminID, &req)
	if err != nil {
		if err == utils.ErrResourceNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Payment not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to verify payment: " + err.Error(),
		})
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  utils.StatusSuccess,
		"message": "Payment " + req.Status + " successfully",
	})
}

// GetSubscriptionInfo handles retrieving subscription information for an admin
func (h *PaymentHandler) GetSubscriptionInfo(c *fiber.Ctx) error {
	// Get admin ID from context or params
	var adminID int
	var ok bool

	role, _ := c.Locals(utils.ContextUserRole).(string)

	if role == utils.RoleSuperAdmin && c.Params("id") != "" {
		// Super admin can check any admin's subscription
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Invalid admin ID",
			})
		}
		adminID = id
	} else {
		// Regular admins can only check their own subscription
		adminID, ok = c.Locals(utils.ContextUserID).(int)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Unauthorized",
			})
		}
	}

	// Check and update subscription status if needed
	admin, err := h.paymentService.CheckSubscriptionStatus(c.Context(), adminID)
	if err != nil {
		if err == utils.ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Admin not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to check subscription status",
		})
	}

	// Get detailed subscription info
	admin, tier, latestPayment, err := h.paymentService.GetSubscriptionInfo(c.Context(), adminID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to retrieve subscription information",
		})
	}

	// Calculate subscription fee
	fee, recommendedTier, err := h.paymentService.CalculateMonthlySubscriptionFee(c.Context(), admin.Users)
	if err != nil {
		fee = 0
	}

	// Prepare response
	response := fiber.Map{
		"admin":                admin.ToResponse(),
		"monthly_fee":          fee,
		"subscription_status":  admin.SubscriptionStatus,
		"is_access_restricted": admin.IsAccessRestricted,
	}

	if tier != nil {
		response["current_tier"] = tier.ToResponse()
	}

	if recommendedTier != nil && (tier == nil || recommendedTier.ID != tier.ID) {
		response["recommended_tier"] = recommendedTier.ToResponse()
	}

	if latestPayment != nil {
		response["latest_payment"] = latestPayment.ToResponse()
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": utils.StatusSuccess,
		"data":   response,
	})
}

package handlers

import (
	"mobilka/internal/models"
	"mobilka/internal/service"
	"mobilka/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// SuperAdminLogin handles super admin login requests
func (h *AuthHandler) SuperAdminLogin(c *fiber.Ctx) error {
	var req models.SuperAdminLoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid request body",
		})
	}

	// Validate request
	if req.Login == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Login and password are required",
		})
	}

	// Attempt login
	superAdmin, token, err := h.authService.SuperAdminLogin(c.Context(), req.Login, req.Password)
	if err != nil {
		if err == utils.ErrInvalidCredentials {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Invalid credentials",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Login failed",
		})
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": utils.StatusSuccess,
		"data": fiber.Map{
			"user":  superAdmin.ToResponse(),
			"token": token,
		},
	})
}

// AdminLogin handles admin login requests
func (h *AuthHandler) AdminLogin(c *fiber.Ctx) error {
	var req models.AdminLoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid request body",
		})
	}

	// Validate request
	if req.UserName == "" || req.SystemID == "" || req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Username, system ID, and email are required",
		})
	}

	// Attempt login
	admin, token, err := h.authService.AdminLogin(c.Context(), req.UserName, req.SystemID, req.Email)
	if err != nil {
		if err == utils.ErrInvalidCredentials {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Invalid credentials",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Login failed",
		})
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": utils.StatusSuccess,
		"data": fiber.Map{
			"user":  admin.ToResponse(),
			"token": token,
		},
	})
}

// SuperAdminChangePassword handles password change for super admin
func (h *AuthHandler) SuperAdminChangePassword(c *fiber.Ctx) error {
	// Get super admin ID from context
	userID, ok := c.Locals(utils.ContextUserID).(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Unauthorized",
		})
	}

	// Parse request
	var req struct {
		OldPassword string `json:"old_password" validate:"required"`
		NewPassword string `json:"new_password" validate:"required,min=8"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Invalid request body",
		})
	}

	// Validate request
	if req.OldPassword == "" || req.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Old and new passwords are required",
		})
	}

	// Change password
	err := h.authService.SuperAdminChangePassword(c.Context(), userID, req.OldPassword, req.NewPassword)
	if err != nil {
		if err == utils.ErrInvalidCredentials {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Invalid old password",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Password change failed",
		})
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  utils.StatusSuccess,
		"message": "Password changed successfully",
	})
}

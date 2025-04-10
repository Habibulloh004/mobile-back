package service

import (
	"context"

	"mobilka/internal/models"
	"mobilka/internal/repository"
	"mobilka/internal/utils"
)

// AuthService handles authentication operations
type AuthService struct {
	superAdminRepo *repository.SuperAdminRepository
	adminRepo      *repository.AdminRepository
}

// NewAuthService creates a new authentication service
func NewAuthService(
	superAdminRepo *repository.SuperAdminRepository,
	adminRepo *repository.AdminRepository,
) *AuthService {
	return &AuthService{
		superAdminRepo: superAdminRepo,
		adminRepo:      adminRepo,
	}
}

// SuperAdminLogin handles super admin login
func (s *AuthService) SuperAdminLogin(ctx context.Context, login, password string) (*models.SuperAdmin, string, error) {
	// Get super admin by login
	superAdmin, err := s.superAdminRepo.GetByLogin(ctx, login)
	if err != nil {
		return nil, "", utils.ErrInvalidCredentials
	}

	// Verify password
	if !utils.CheckPassword(password, superAdmin.Password) {
		return nil, "", utils.ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := utils.GenerateSuperAdminToken(superAdmin)
	if err != nil {
		return nil, "", err
	}

	return superAdmin, token, nil
}

// AdminLogin handles admin login
func (s *AuthService) AdminLogin(ctx context.Context, userName, systemID, email string) (*models.Admin, string, error) {
	// Get admin by username, system ID, and email
	admin, err := s.adminRepo.GetByCredentials(ctx, userName, systemID, email)
	if err != nil {
		return nil, "", utils.ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := utils.GenerateAdminToken(admin)
	if err != nil {
		return nil, "", err
	}

	return admin, token, nil
}

// SuperAdminChangePassword handles super admin password change
func (s *AuthService) SuperAdminChangePassword(ctx context.Context, id int, oldPassword, newPassword string) error {
	// Get super admin by ID
	superAdmin, err := s.superAdminRepo.GetByID(ctx, id)
	if err != nil {
		return utils.ErrUserNotFound
	}

	// Verify old password
	if !utils.CheckPassword(oldPassword, superAdmin.Password) {
		return utils.ErrInvalidCredentials
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update password
	return s.superAdminRepo.UpdatePassword(ctx, id, hashedPassword)
}

// SetupDefaultSuperAdmin sets up the default super admin account
func (s *AuthService) SetupDefaultSuperAdmin(ctx context.Context) (string, error) {
	// Generate password (now fixed to "helloworld")
	plainPassword, hashedPassword, err := utils.GenerateSuperAdminPassword()
	if err != nil {
		return "", err
	}

	// Setup default super admin
	err = s.superAdminRepo.SetupDefaultSuperAdmin(ctx, "superadmin", hashedPassword)
	if err != nil {
		return "", err
	}

	return plainPassword, nil
}

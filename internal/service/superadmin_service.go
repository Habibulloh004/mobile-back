package service

import (
	"context"

	"mobilka/internal/models"
	"mobilka/internal/repository"
)

// SuperAdminService handles super admin operations
type SuperAdminService struct {
	superAdminRepo *repository.SuperAdminRepository
}

// NewSuperAdminService creates a new super admin service
func NewSuperAdminService(
	superAdminRepo *repository.SuperAdminRepository,
) *SuperAdminService {
	return &SuperAdminService{
		superAdminRepo: superAdminRepo,
	}
}

// GetByID retrieves a super admin by ID
func (s *SuperAdminService) GetByID(ctx context.Context, id int) (*models.SuperAdmin, error) {
	return s.superAdminRepo.GetByID(ctx, id)
}

// GetByLogin retrieves a super admin by login
func (s *SuperAdminService) GetByLogin(ctx context.Context, login string) (*models.SuperAdmin, error) {
	return s.superAdminRepo.GetByLogin(ctx, login)
}

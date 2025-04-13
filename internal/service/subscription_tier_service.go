package service

import (
	"context"

	"mobilka/internal/models"
	"mobilka/internal/repository"
)

// SubscriptionTierService handles subscription tier operations
type SubscriptionTierService struct {
	subscriptionTierRepo *repository.SubscriptionTierRepository
}

// NewSubscriptionTierService creates a new subscription tier service
func NewSubscriptionTierService(subscriptionTierRepo *repository.SubscriptionTierRepository) *SubscriptionTierService {
	return &SubscriptionTierService{
		subscriptionTierRepo: subscriptionTierRepo,
	}
}

// Create creates a new subscription tier
func (s *SubscriptionTierService) Create(ctx context.Context, req *models.SubscriptionTierCreateRequest) (*models.SubscriptionTier, error) {
	tier := &models.SubscriptionTier{
		Name:        req.Name,
		MinUsers:    req.MinUsers,
		MaxUsers:    req.MaxUsers,
		Price:       req.Price,
		Description: req.Description,
	}

	err := s.subscriptionTierRepo.Create(ctx, tier)
	if err != nil {
		return nil, err
	}

	return tier, nil
}

// GetByID retrieves a subscription tier by ID
func (s *SubscriptionTierService) GetByID(ctx context.Context, id int) (*models.SubscriptionTier, error) {
	return s.subscriptionTierRepo.GetByID(ctx, id)
}

// GetAll retrieves all subscription tiers
func (s *SubscriptionTierService) GetAll(ctx context.Context) ([]*models.SubscriptionTier, error) {
	return s.subscriptionTierRepo.GetAll(ctx)
}

// Update updates a subscription tier
func (s *SubscriptionTierService) Update(ctx context.Context, id int, req *models.SubscriptionTierUpdateRequest) (*models.SubscriptionTier, error) {
	// First get the current tier
	tier, err := s.subscriptionTierRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" {
		tier.Name = req.Name
	}

	if req.MinUsers != 0 {
		tier.MinUsers = req.MinUsers
	}

	// Only update MaxUsers if explicitly provided in the request
	// The check is explicit because MaxUsers is a pointer and can be legitimately set to nil
	if req.MaxUsers != nil {
		tier.MaxUsers = req.MaxUsers
	}

	if req.Price != 0 {
		tier.Price = req.Price
	}

	if req.Description != "" {
		tier.Description = req.Description
	}

	// Update in database
	err = s.subscriptionTierRepo.Update(ctx, id, tier)
	if err != nil {
		return nil, err
	}

	return tier, nil
}

// Delete deletes a subscription tier
func (s *SubscriptionTierService) Delete(ctx context.Context, id int) error {
	return s.subscriptionTierRepo.Delete(ctx, id)
}

// GetTierForUserCount retrieves the appropriate subscription tier for a given user count
func (s *SubscriptionTierService) GetTierForUserCount(ctx context.Context, userCount int) (*models.SubscriptionTier, error) {
	return s.subscriptionTierRepo.GetTierForUserCount(ctx, userCount)
}

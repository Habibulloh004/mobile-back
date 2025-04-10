package service

import (
	"context"

	"mobilka/internal/models"
	"mobilka/internal/repository"
)

// BannerService handles banner operations
type BannerService struct {
	bannerRepo *repository.BannerRepository
}

// NewBannerService creates a new banner service
func NewBannerService(bannerRepo *repository.BannerRepository) *BannerService {
	return &BannerService{
		bannerRepo: bannerRepo,
	}
}

// Create creates a new banner
func (s *BannerService) Create(ctx context.Context, adminID int, req *models.BannerCreateRequest) (*models.Banner, error) {
	banner := &models.Banner{
		AdminID: adminID,
		Image:   req.Image,
		Title:   req.Title,
		Body:    req.Body,
	}

	err := s.bannerRepo.Create(ctx, banner)
	if err != nil {
		return nil, err
	}

	return banner, nil
}

// GetByID retrieves a banner by ID
func (s *BannerService) GetByID(ctx context.Context, id int) (*models.Banner, error) {
	return s.bannerRepo.GetByID(ctx, id)
}

// GetByAdminID retrieves all banners for a specific admin
func (s *BannerService) GetByAdminID(ctx context.Context, adminID int) ([]*models.Banner, error) {
	return s.bannerRepo.GetByAdminID(ctx, adminID)
}

// GetAll retrieves all banners
func (s *BannerService) GetAll(ctx context.Context) ([]*models.Banner, error) {
	return s.bannerRepo.GetAll(ctx)
}

// Update updates a banner
func (s *BannerService) Update(ctx context.Context, id int, adminID int, req *models.BannerUpdateRequest) (*models.Banner, error) {
	// Get existing banner
	banner, err := s.bannerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Image != "" {
		banner.Image = req.Image
	}

	if req.Title != "" {
		banner.Title = req.Title
	}

	if req.Body != "" {
		banner.Body = req.Body
	}

	// Ensure the admin can only update their own banners
	banner.AdminID = adminID

	// Update in database
	err = s.bannerRepo.Update(ctx, id, banner)
	if err != nil {
		return nil, err
	}

	return banner, nil
}

// Delete deletes a banner
func (s *BannerService) Delete(ctx context.Context, id int, adminID int) error {
	return s.bannerRepo.Delete(ctx, id, adminID)
}

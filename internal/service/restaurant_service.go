package service

import (
	"context"
	"fmt"

	"mobilka/internal/models"
	"mobilka/internal/repository"
)

// RestaurantService handles restaurant operations
type RestaurantService struct {
	restaurantRepo *repository.RestaurantRepository
}

// NewRestaurantService creates a new restaurant service
func NewRestaurantService(restaurantRepo *repository.RestaurantRepository) *RestaurantService {
	return &RestaurantService{
		restaurantRepo: restaurantRepo,
	}
}

// Create creates a new restaurant or updates an existing one if the admin already has one
func (s *RestaurantService) Create(ctx context.Context, adminID int, req *models.RestaurantCreateRequest) (*models.Restaurant, error) {
	// Create a new restaurant with the provided admin ID
	restaurant := &models.Restaurant{
		AdminID:     adminID,
		Text:        req.Text,
		Contacts:    req.Contacts,
		SocialMedia: req.SocialMedia,
	}

	// Log the restaurant object before saving (for debugging)
	fmt.Printf("Creating/updating restaurant with admin ID: %d\n", restaurant.AdminID)

	// Save the restaurant to the database - will update if one already exists
	err := s.restaurantRepo.Create(ctx, restaurant)
	if err != nil {
		return nil, err
	}

	return restaurant, nil
}

// GetByID retrieves a restaurant by ID
func (s *RestaurantService) GetByID(ctx context.Context, id int) (*models.Restaurant, error) {
	return s.restaurantRepo.GetByID(ctx, id)
}

// GetByAdminID retrieves all restaurants for a specific admin
func (s *RestaurantService) GetByAdminID(ctx context.Context, adminID int) ([]*models.Restaurant, error) {
	return s.restaurantRepo.GetByAdminID(ctx, adminID)
}

// GetAll retrieves all restaurants
func (s *RestaurantService) GetAll(ctx context.Context) ([]*models.Restaurant, error) {
	return s.restaurantRepo.GetAll(ctx)
}

// Update updates a restaurant
func (s *RestaurantService) Update(ctx context.Context, id int, adminID int, req *models.RestaurantUpdateRequest) (*models.Restaurant, error) {
	// First get the current restaurant
	restaurant, err := s.restaurantRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Text != "" {
		restaurant.Text = req.Text
	}

	// Update nested objects if provided
	if req.Contacts != nil {
		// Update only non-empty fields to preserve existing values
		if req.Contacts.Phone != "" {
			restaurant.Contacts.Phone = req.Contacts.Phone
		}
		if req.Contacts.Gmail != "" {
			restaurant.Contacts.Gmail = req.Contacts.Gmail
		}
		if req.Contacts.Location != "" {
			restaurant.Contacts.Location = req.Contacts.Location
		}
	}

	if req.SocialMedia != nil {
		// Update only non-empty fields to preserve existing values
		if req.SocialMedia.Instagram != "" {
			restaurant.SocialMedia.Instagram = req.SocialMedia.Instagram
		}
		if req.SocialMedia.Telegram != "" {
			restaurant.SocialMedia.Telegram = req.SocialMedia.Telegram
		}
		if req.SocialMedia.Facebook != "" {
			restaurant.SocialMedia.Facebook = req.SocialMedia.Facebook
		}
	}

	// Update admin ID if specified (for super admin)
	restaurant.AdminID = adminID

	// Log update operation for debugging
	fmt.Printf("Updating restaurant ID %d with admin ID %d\n", id, restaurant.AdminID)

	// Update in database
	err = s.restaurantRepo.Update(ctx, id, restaurant)
	if err != nil {
		return nil, err
	}

	return restaurant, nil
}

// Delete deletes a restaurant
func (s *RestaurantService) Delete(ctx context.Context, id int, adminID int) error {
	return s.restaurantRepo.Delete(ctx, id, adminID)
}

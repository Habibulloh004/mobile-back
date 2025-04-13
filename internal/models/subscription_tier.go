package models

import (
	"time"
)

// SubscriptionTier represents a pricing tier for admin subscriptions
type SubscriptionTier struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	MinUsers    int       `json:"min_users"`
	MaxUsers    *int      `json:"max_users"` // Pointer to allow NULL for unlimited users
	Price       float64   `json:"price"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SubscriptionTierCreateRequest represents the request to create a subscription tier
type SubscriptionTierCreateRequest struct {
	Name        string  `json:"name" validate:"required"`
	MinUsers    int     `json:"min_users" validate:"required,min=0"`
	MaxUsers    *int    `json:"max_users"`
	Price       float64 `json:"price" validate:"required,min=0"`
	Description string  `json:"description"`
}

// SubscriptionTierUpdateRequest represents the request to update a subscription tier
type SubscriptionTierUpdateRequest struct {
	Name        string  `json:"name"`
	MinUsers    int     `json:"min_users" validate:"min=0"`
	MaxUsers    *int    `json:"max_users"`
	Price       float64 `json:"price" validate:"min=0"`
	Description string  `json:"description"`
}

// SubscriptionTierResponse represents the response for a subscription tier
type SubscriptionTierResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	MinUsers    int       `json:"min_users"`
	MaxUsers    *int      `json:"max_users"`
	Price       float64   `json:"price"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToResponse converts SubscriptionTier to SubscriptionTierResponse
func (s *SubscriptionTier) ToResponse() SubscriptionTierResponse {
	return SubscriptionTierResponse{
		ID:          s.ID,
		Name:        s.Name,
		MinUsers:    s.MinUsers,
		MaxUsers:    s.MaxUsers,
		Price:       s.Price,
		Description: s.Description,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

package models

import (
	"time"
)

// Restaurant model represents the restaurant entity
type Restaurant struct {
	ID          int         `json:"id"`
	AdminID     int         `json:"admin_id"`
	Text        string      `json:"text"`
	Contacts    Contacts    `json:"contacts"`
	SocialMedia SocialMedia `json:"social_media"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// Contacts represents contact information for a restaurant
type Contacts struct {
	Phone    string `json:"phone"`
	Gmail    string `json:"gmail"`
	Location string `json:"location"`
}

// SocialMedia represents social media links for a restaurant
type SocialMedia struct {
	Instagram string `json:"instagram"`
	Telegram  string `json:"telegram"`
	Facebook  string `json:"facebook"`
}

// RestaurantCreateRequest represents the creation request for a restaurant
type RestaurantCreateRequest struct {
	AdminID     int         `json:"admin_id"`
	Text        string      `json:"text" validate:"required"`
	Contacts    Contacts    `json:"contacts"`
	SocialMedia SocialMedia `json:"social_media"`
}

// RestaurantUpdateRequest represents the update request for a restaurant
type RestaurantUpdateRequest struct {
	AdminID     int          `json:"admin_id"`
	Text        string       `json:"text"`
	Contacts    *Contacts    `json:"contacts"`
	SocialMedia *SocialMedia `json:"social_media"`
}

// RestaurantResponse represents the response for restaurant
type RestaurantResponse struct {
	ID          int         `json:"id"`
	AdminID     int         `json:"admin_id"`
	Text        string      `json:"text"`
	Contacts    Contacts    `json:"contacts"`
	SocialMedia SocialMedia `json:"social_media"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// ToResponse converts Restaurant model to RestaurantResponse
func (r *Restaurant) ToResponse() RestaurantResponse {
	return RestaurantResponse{
		ID:          r.ID,
		AdminID:     r.AdminID,
		Text:        r.Text,
		Contacts:    r.Contacts,
		SocialMedia: r.SocialMedia,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

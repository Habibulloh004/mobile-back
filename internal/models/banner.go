package models

import (
	"time"
)

// Banner model represents the banner entity
type Banner struct {
	ID        int       `json:"id"`
	AdminID   int       `json:"admin_id"`
	Image     string    `json:"image"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BannerCreateRequest represents the creation request for a banner
type BannerCreateRequest struct {
	Image string `json:"image" validate:"required"`
	Title string `json:"title" validate:"required"`
	Body  string `json:"body" validate:"required"`
}

// BannerUpdateRequest represents the update request for a banner
type BannerUpdateRequest struct {
	Image string `json:"image"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

// BannerResponse represents the response for banner
type BannerResponse struct {
	ID        int       `json:"id"`
	AdminID   int       `json:"admin_id"`
	Image     string    `json:"image"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse converts Banner model to BannerResponse
func (b *Banner) ToResponse() BannerResponse {
	return BannerResponse{
		ID:        b.ID,
		AdminID:   b.AdminID,
		Image:     b.Image,
		Title:     b.Title,
		Body:      b.Body,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}

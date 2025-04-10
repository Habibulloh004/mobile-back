package models

import (
	"time"
)

// FCMToken model represents the FCM token entity
type FCMToken struct {
	ID        int       `json:"id"`
	AdminID   int       `json:"admin_id"`
	FCMToken  string    `json:"fcm_token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FCMTokenCreateRequest represents the creation request for an FCM token
type FCMTokenCreateRequest struct {
	FCMToken string `json:"fcm_token" validate:"required"`
}

// FCMTokenResponse represents the response for FCM token
type FCMTokenResponse struct {
	ID        int       `json:"id"`
	AdminID   int       `json:"admin_id"`
	FCMToken  string    `json:"fcm_token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse converts FCMToken model to FCMTokenResponse
func (ft *FCMToken) ToResponse() FCMTokenResponse {
	return FCMTokenResponse{
		ID:        ft.ID,
		AdminID:   ft.AdminID,
		FCMToken:  ft.FCMToken,
		CreatedAt: ft.CreatedAt,
		UpdatedAt: ft.UpdatedAt,
	}
}

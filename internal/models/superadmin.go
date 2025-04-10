package models

import (
	"time"
)

// SuperAdmin model represents the super admin entity
type SuperAdmin struct {
	ID        int       `json:"id"`
	Login     string    `json:"login"`
	Password  string    `json:"-"` // Password is not exposed in JSON responses
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SuperAdminLoginRequest represents the login request for super admin
type SuperAdminLoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// SuperAdminResponse represents the response for super admin without sensitive data
type SuperAdminResponse struct {
	ID        int       `json:"id"`
	Login     string    `json:"login"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse converts SuperAdmin model to SuperAdminResponse
func (sa *SuperAdmin) ToResponse() SuperAdminResponse {
	return SuperAdminResponse{
		ID:        sa.ID,
		Login:     sa.Login,
		CreatedAt: sa.CreatedAt,
		UpdatedAt: sa.UpdatedAt,
	}
}

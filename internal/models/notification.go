package models

import (
	"time"
)

// Notification model represents the notification entity
type Notification struct {
	ID        int       `json:"id"`
	AdminID   int       `json:"admin_id"`
	Payload   string    `json:"payload"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NotificationCreateRequest represents the creation request for a notification
type NotificationCreateRequest struct {
	Payload string `json:"payload" validate:"required"`
	Title   string `json:"title" validate:"required"`
	Body    string `json:"body" validate:"required"`
}

// NotificationUpdateRequest represents the update request for a notification
type NotificationUpdateRequest struct {
	Payload string `json:"payload"`
	Title   string `json:"title"`
	Body    string `json:"body"`
}

// NotificationResponse represents the response for notification
type NotificationResponse struct {
	ID        int       `json:"id"`
	AdminID   int       `json:"admin_id"`
	Payload   string    `json:"payload"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse converts Notification model to NotificationResponse
func (n *Notification) ToResponse() NotificationResponse {
	return NotificationResponse{
		ID:        n.ID,
		AdminID:   n.AdminID,
		Payload:   n.Payload,
		Title:     n.Title,
		Body:      n.Body,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}
}

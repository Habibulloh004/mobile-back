package service

import (
	"context"
	"fmt"

	"mobilka/internal/models"
	"mobilka/internal/repository"
)

// NotificationService handles notification operations
type NotificationService struct {
	notificationRepo *repository.NotificationRepository
	fcmTokenRepo     *repository.FCMTokenRepository
}

// NewNotificationService creates a new notification service
func NewNotificationService(
	notificationRepo *repository.NotificationRepository,
	fcmTokenRepo *repository.FCMTokenRepository,
) *NotificationService {
	return &NotificationService{
		notificationRepo: notificationRepo,
		fcmTokenRepo:     fcmTokenRepo,
	}
}

// Create creates a new notification
func (s *NotificationService) Create(ctx context.Context, adminID int, req *models.NotificationCreateRequest) (*models.Notification, error) {
	notification := &models.Notification{
		AdminID: adminID,
		Payload: req.Payload,
		Title:   req.Title,
		Body:    req.Body,
	}

	err := s.notificationRepo.Create(ctx, notification)
	if err != nil {
		return nil, err
	}

	return notification, nil
}

// GetByID retrieves a notification by ID
func (s *NotificationService) GetByID(ctx context.Context, id int) (*models.Notification, error) {
	return s.notificationRepo.GetByID(ctx, id)
}

// GetByAdminID retrieves all notifications for a specific admin
func (s *NotificationService) GetByAdminID(ctx context.Context, adminID int) ([]*models.Notification, error) {
	return s.notificationRepo.GetByAdminID(ctx, adminID)
}

// GetAll retrieves all notifications
func (s *NotificationService) GetAll(ctx context.Context) ([]*models.Notification, error) {
	return s.notificationRepo.GetAll(ctx)
}

// Update updates a notification
func (s *NotificationService) Update(ctx context.Context, id int, adminID int, req *models.NotificationUpdateRequest) (*models.Notification, error) {
	// First get the current notification
	notification, err := s.notificationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Payload != "" {
		notification.Payload = req.Payload
	}

	if req.Title != "" {
		notification.Title = req.Title
	}

	if req.Body != "" {
		notification.Body = req.Body
	}

	// Update admin ID if specified (for super admin)
	notification.AdminID = adminID

	// Log update operation for debugging
	fmt.Printf("Updating notification ID %d with admin ID %d\n", id, notification.AdminID)

	// Update in database
	err = s.notificationRepo.Update(ctx, id, notification)
	if err != nil {
		return nil, err
	}

	return notification, nil
}

// Delete deletes a notification
func (s *NotificationService) Delete(ctx context.Context, id int, adminID int) error {
	return s.notificationRepo.Delete(ctx, id, adminID)
}

// GetByAdminIDWithPagination retrieves notifications for a specific admin with pagination
func (s *NotificationService) GetByAdminIDWithPagination(ctx context.Context, adminID, skip, step int) ([]*models.Notification, error) {
	return s.notificationRepo.GetByAdminIDWithPagination(ctx, adminID, skip, step)
}

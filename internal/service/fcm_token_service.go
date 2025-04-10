package service

import (
	"context"

	"mobilka/internal/models"
	"mobilka/internal/repository"
)

// FCMTokenService handles FCM token operations
type FCMTokenService struct {
	fcmTokenRepo *repository.FCMTokenRepository
}

// NewFCMTokenService creates a new FCM token service
func NewFCMTokenService(
	fcmTokenRepo *repository.FCMTokenRepository,
) *FCMTokenService {
	return &FCMTokenService{
		fcmTokenRepo: fcmTokenRepo,
	}
}

// Create creates a new FCM token
func (s *FCMTokenService) Create(ctx context.Context, adminID int, req *models.FCMTokenCreateRequest) (*models.FCMToken, error) {
	fcmToken := &models.FCMToken{
		AdminID:  adminID,
		FCMToken: req.FCMToken,
	}

	err := s.fcmTokenRepo.Create(ctx, fcmToken)
	if err != nil {
		return nil, err
	}

	return fcmToken, nil
}

// GetByID retrieves an FCM token by ID
func (s *FCMTokenService) GetByID(ctx context.Context, id int) (*models.FCMToken, error) {
	return s.fcmTokenRepo.GetByID(ctx, id)
}

// GetByToken retrieves an FCM token by token string
func (s *FCMTokenService) GetByToken(ctx context.Context, token string) (*models.FCMToken, error) {
	return s.fcmTokenRepo.GetByToken(ctx, token)
}

// GetByAdminID retrieves all FCM tokens for a specific admin
func (s *FCMTokenService) GetByAdminID(ctx context.Context, adminID int) ([]*models.FCMToken, error) {
	return s.fcmTokenRepo.GetByAdminID(ctx, adminID)
}

// GetAll retrieves all FCM tokens
func (s *FCMTokenService) GetAll(ctx context.Context) ([]*models.FCMToken, error) {
	return s.fcmTokenRepo.GetAll(ctx)
}

// Delete deletes an FCM token
func (s *FCMTokenService) Delete(ctx context.Context, id int) error {
	return s.fcmTokenRepo.Delete(ctx, id)
}

// DeleteByToken deletes an FCM token by token string
func (s *FCMTokenService) DeleteByToken(ctx context.Context, token string) error {
	return s.fcmTokenRepo.DeleteByToken(ctx, token)
}

// DeleteByAdminID deletes all FCM tokens for a specific admin
func (s *FCMTokenService) DeleteByAdminID(ctx context.Context, adminID int) error {
	return s.fcmTokenRepo.DeleteByAdminID(ctx, adminID)
}

package service

import (
	"context"
	// "errors"
	// "fmt"
	"time"

	"mobilka/internal/models"
	"mobilka/internal/repository"
	"mobilka/internal/utils"
	// "github.com/google/uuid"
)

// AdminService handles admin operations
type AdminService struct {
	adminRepo *repository.AdminRepository
}

// NewAdminService creates a new admin service
func NewAdminService(
	adminRepo *repository.AdminRepository,
) *AdminService {
	return &AdminService{
		adminRepo: adminRepo,
	}
}

// Create creates a new admin
func (s *AdminService) Create(ctx context.Context, req *models.AdminCreateRequest) (*models.Admin, error) {

	// Hash passwords
	var smsPasswordHash string
	var paymentPasswordHash string
	var err error

	if req.SmsPassword != "" {
		smsPasswordHash, err = utils.HashPassword(req.SmsPassword)
		if err != nil {
			return nil, err
		}
	}

	if req.PaymentPassword != "" {
		paymentPasswordHash, err = utils.HashPassword(req.PaymentPassword)
		if err != nil {
			return nil, err
		}
	}

	// Create admin
	admin := &models.Admin{
		UserName:               req.UserName,
		Email:                  req.Email,
		CompanyName:            req.CompanyName,
		SystemID:               req.SystemID,
		SystemToken:            req.SystemToken,
		SystemTokenUpdatedTime: time.Now(),
		SmsToken:               req.SmsToken,
		SmsTokenUpdatedTime:    time.Now(),
		SmsEmail:               req.SmsEmail,
		SmsPassword:            smsPasswordHash,
		SmsMessage:             req.SmsMessage,
		PaymentUsername:        req.PaymentUsername,
		PaymentPassword:        paymentPasswordHash,
	}

	// Save to database
	err = s.adminRepo.Create(ctx, admin)
	if err != nil {
		return nil, err
	}

	return admin, nil
}

// GetByID retrieves an admin by ID
func (s *AdminService) GetByID(ctx context.Context, id int) (*models.Admin, error) {
	return s.adminRepo.GetByID(ctx, id)
}

// GetByEmail retrieves an admin by email
func (s *AdminService) GetByEmail(ctx context.Context, email string) (*models.Admin, error) {
	return s.adminRepo.GetByEmail(ctx, email)
}

// GetAll retrieves all admins
func (s *AdminService) GetAll(ctx context.Context) ([]*models.Admin, error) {
	return s.adminRepo.GetAll(ctx)
}

// Update updates an admin
func (s *AdminService) Update(ctx context.Context, id int, req *models.AdminUpdateRequest) (*models.Admin, error) {
	// Get existing admin
	admin, err := s.adminRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.UserName != "" {
		admin.UserName = req.UserName
	}

	if req.Email != "" {
		admin.Email = req.Email
	}

	if req.CompanyName != "" {
		admin.CompanyName = req.CompanyName
	}

	if req.SystemID != "" {
		admin.SystemID = req.SystemID
	}

	// Update system token if provided
	if req.SystemToken != "" {
		admin.SystemToken = req.SystemToken
		admin.SystemTokenUpdatedTime = time.Now()
	}

	// Update SMS token if provided
	if req.SmsToken != "" {
		admin.SmsToken = req.SmsToken
		admin.SmsTokenUpdatedTime = time.Now()
	}

	if req.SmsEmail != "" {
		admin.SmsEmail = req.SmsEmail
	}

	if req.SmsPassword != "" {
		smsPasswordHash, err := utils.HashPassword(req.SmsPassword)
		if err != nil {
			return nil, err
		}
		admin.SmsPassword = smsPasswordHash
	}

	if req.SmsMessage != "" {
		admin.SmsMessage = req.SmsMessage
	}

	if req.PaymentUsername != "" {
		admin.PaymentUsername = req.PaymentUsername
	}

	if req.PaymentPassword != "" {
		paymentPasswordHash, err := utils.HashPassword(req.PaymentPassword)
		if err != nil {
			return nil, err
		}
		admin.PaymentPassword = paymentPasswordHash
	}

	// Update in database
	err = s.adminRepo.Update(ctx, id, admin)
	if err != nil {
		return nil, err
	}

	return admin, nil
}

// Delete deletes an admin
func (s *AdminService) Delete(ctx context.Context, id int) error {
	return s.adminRepo.Delete(ctx, id)
}

// GetByIDPublic retrieves an admin by ID and increments the users count
func (s *AdminService) GetByIDPublic(ctx context.Context, id int) (*models.Admin, error) {
	// Get admin
	admin, err := s.adminRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Increment users count
	err = s.adminRepo.IncrementUsersCount(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update local admin object to reflect the incremented count
	admin.Users++

	return admin, nil
}

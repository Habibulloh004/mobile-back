package service

import (
	"context"
	"time"

	"mobilka/internal/models"
	"mobilka/internal/repository"
	"mobilka/internal/utils"
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
	smsPasswordHash, err := utils.HashPassword(req.SmsPassword)
	if err != nil {
		return nil, err
	}

	paymentPasswordHash, err := utils.HashPassword(req.PaymentPassword)
	if err != nil {
		return nil, err
	}

	// Generate tokens
	systemToken, err := utils.GenerateSystemToken()
	if err != nil {
		return nil, err
	}

	smsToken, err := utils.GenerateSmsToken()
	if err != nil {
		return nil, err
	}

	// Create admin
	admin := &models.Admin{
		UserName:               req.UserName,
		Email:                  req.Email,
		CompanyName:            req.CompanyName,
		SystemID:               req.SystemID,
		SystemToken:            systemToken,
		SystemTokenUpdatedTime: time.Now(),
		SmsToken:               smsToken,
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

// RegenerateSystemToken regenerates the system token for an admin
func (s *AdminService) RegenerateSystemToken(ctx context.Context, id int) (string, error) {
	// Generate new token
	systemToken, err := utils.GenerateSystemToken()
	if err != nil {
		return "", err
	}

	// Update in database
	err = s.adminRepo.UpdateSystemToken(ctx, id, systemToken)
	if err != nil {
		return "", err
	}

	return systemToken, nil
}

// RegenerateSmsToken regenerates the SMS token for an admin
func (s *AdminService) RegenerateSmsToken(ctx context.Context, id int) (string, error) {
	// Generate new token
	smsToken, err := utils.GenerateSmsToken()
	if err != nil {
		return "", err
	}

	// Update in database
	err = s.adminRepo.UpdateSmsToken(ctx, id, smsToken)
	if err != nil {
		return "", err
	}

	return smsToken, nil
}

// Delete deletes an admin
func (s *AdminService) Delete(ctx context.Context, id int) error {
	return s.adminRepo.Delete(ctx, id)
}

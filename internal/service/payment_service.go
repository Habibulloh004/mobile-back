package service

import (
	"context"
	"time"

	"mobilka/internal/models"
	"mobilka/internal/repository"
	"mobilka/internal/utils"
)

// PaymentService handles payment operations
type PaymentService struct {
	paymentRepo          *repository.PaymentHistoryRepository
	adminRepo            *repository.AdminRepository
	subscriptionTierRepo *repository.SubscriptionTierRepository
}

// NewPaymentService creates a new payment service
func NewPaymentService(
	paymentRepo *repository.PaymentHistoryRepository,
	adminRepo *repository.AdminRepository,
	subscriptionTierRepo *repository.SubscriptionTierRepository,
) *PaymentService {
	return &PaymentService{
		paymentRepo:          paymentRepo,
		adminRepo:            adminRepo,
		subscriptionTierRepo: subscriptionTierRepo,
	}
}

// RecordPayment records a new payment from an admin
func (s *PaymentService) RecordPayment(ctx context.Context, adminID int, req *models.PaymentCreateRequest) (*models.PaymentHistory, error) {
	// Get admin to check if they exist
	admin, err := s.adminRepo.GetByID(ctx, adminID)
	if err != nil {
		return nil, err
	}

	// Determine appropriate subscription tier based on user count
	tier, err := s.subscriptionTierRepo.GetTierForUserCount(ctx, admin.Users)
	if err != nil {
		// If no specific tier found, don't associate payment with a tier
		tier = nil
	}

	// Create payment record
	payment := &models.PaymentHistory{
		AdminID:       adminID,
		Amount:        req.Amount,
		PaymentDate:   time.Now(),
		PaymentMethod: req.PaymentMethod,
		TransactionID: req.TransactionID,
		Status:        "pending",
		Notes:         req.Notes,
	}

	// Associate with subscription tier if found
	if tier != nil {
		payment.SubscriptionTierID = &tier.ID
	}

	// Save payment to database
	err = s.paymentRepo.Create(ctx, payment)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

// GetPaymentByID retrieves a payment by ID
func (s *PaymentService) GetPaymentByID(ctx context.Context, id int) (*models.PaymentHistory, error) {
	return s.paymentRepo.GetByID(ctx, id)
}

// GetPaymentsByAdminID retrieves all payments for a specific admin
func (s *PaymentService) GetPaymentsByAdminID(ctx context.Context, adminID int) ([]*models.PaymentHistory, error) {
	// Verify admin exists
	_, err := s.adminRepo.GetByID(ctx, adminID)
	if err != nil {
		return nil, err
	}

	return s.paymentRepo.GetByAdminID(ctx, adminID)
}

// GetAllPayments retrieves all payments
func (s *PaymentService) GetAllPayments(ctx context.Context) ([]*models.PaymentHistory, error) {
	return s.paymentRepo.GetAll(ctx)
}

// GetPendingPayments retrieves all pending payments
func (s *PaymentService) GetPendingPayments(ctx context.Context) ([]*models.PaymentHistory, error) {
	return s.paymentRepo.GetPendingPayments(ctx)
}

// VerifyPayment verifies a payment and updates admin's subscription status
func (s *PaymentService) VerifyPayment(ctx context.Context, paymentID int, superAdminID int, req *models.PaymentVerifyRequest) error {
	// Get payment
	payment, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return err
	}

	// Get admin
	admin, err := s.adminRepo.GetByID(ctx, payment.AdminID)
	if err != nil {
		return err
	}

	// Update payment status
	err = s.paymentRepo.VerifyPayment(
		ctx,
		paymentID,
		superAdminID,
		req.Status,
		req.Notes,
		req.PeriodStart,
		req.PeriodEnd,
	)
	if err != nil {
		return err
	}

	// If payment is verified, update admin subscription status
	if req.Status == "verified" {
		// Default to 1 month subscription period if not specified
		periodEnd := req.PeriodEnd
		if periodEnd == nil {
			defaultEnd := time.Now().AddDate(0, 1, 0) // 1 month from now
			periodEnd = &defaultEnd
		}

		// Update admin subscription status
		err = s.adminRepo.UpdateSubscriptionStatus(
			ctx,
			admin.ID,
			payment.SubscriptionTierID,
			"active",
			periodEnd,
			false, // remove access restriction
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// RejectPayment rejects a payment without updating subscription
func (s *PaymentService) RejectPayment(ctx context.Context, paymentID int, superAdminID int, notes string) error {
	// Get payment
	_, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return err
	}

	// Update payment status
	return s.paymentRepo.VerifyPayment(
		ctx,
		paymentID,
		superAdminID,
		"rejected",
		notes,
		nil,
		nil,
	)
}

// CheckSubscriptionStatus checks admin's subscription status and updates if needed
func (s *PaymentService) CheckSubscriptionStatus(ctx context.Context, adminID int) (*models.Admin, error) {
	// Get admin with subscription info
	admin, tier, err := s.adminRepo.GetByIDWithSubscriptionInfo(ctx, adminID)
	if err != nil {
		return nil, err
	}

	// If admin has no subscription tier, determine appropriate tier
	if tier == nil && admin.SubscriptionTierID == nil {
		newTier, err := s.subscriptionTierRepo.GetTierForUserCount(ctx, admin.Users)
		if err == nil {
			// Update admin with appropriate tier
			admin.SubscriptionTierID = &newTier.ID

			// Don't change subscription status - just associate the tier
			err = s.adminRepo.UpdateSubscriptionStatus(
				ctx,
				admin.ID,
				admin.SubscriptionTierID,
				admin.SubscriptionStatus,
				admin.SubscriptionExpiresAt,
				admin.IsAccessRestricted,
			)
			if err != nil {
				return nil, err
			}
		}
	}

	// Check if subscription has expired
	if admin.SubscriptionStatus == "active" &&
		admin.SubscriptionExpiresAt != nil &&
		admin.SubscriptionExpiresAt.Before(time.Now()) {

		// Update to expired status
		admin.SubscriptionStatus = "expired"
		admin.IsAccessRestricted = true

		err = s.adminRepo.UpdateSubscriptionStatus(
			ctx,
			admin.ID,
			admin.SubscriptionTierID,
			admin.SubscriptionStatus,
			admin.SubscriptionExpiresAt,
			admin.IsAccessRestricted,
		)
		if err != nil {
			return nil, err
		}
	}

	return admin, nil
}

// GetSubscriptionInfo retrieves subscription information for an admin
func (s *PaymentService) GetSubscriptionInfo(ctx context.Context, adminID int) (
	*models.Admin,
	*models.SubscriptionTier,
	*models.PaymentHistory,
	error) {

	// Get admin with subscription tier
	admin, tier, err := s.adminRepo.GetByIDWithSubscriptionInfo(ctx, adminID)
	if err != nil {
		return nil, nil, nil, err
	}

	// Get latest verified payment
	var latestPayment *models.PaymentHistory
	latestPayment, err = s.paymentRepo.GetLatestVerifiedPayment(ctx, adminID)
	if err != nil && err != utils.ErrResourceNotFound {
		return nil, nil, nil, err
	}

	return admin, tier, latestPayment, nil
}

// ExpireSubscriptions checks and expires all overdue subscriptions
func (s *PaymentService) ExpireSubscriptions(ctx context.Context) (int, error) {
	return s.adminRepo.ExpireSubscriptions(ctx)
}

// CalculateMonthlySubscriptionFee calculates the monthly subscription fee based on user count
func (s *PaymentService) CalculateMonthlySubscriptionFee(ctx context.Context, userCount int) (float64, *models.SubscriptionTier, error) {
	tier, err := s.subscriptionTierRepo.GetTierForUserCount(ctx, userCount)
	if err != nil {
		return 0, nil, err
	}

	return tier.Price, tier, nil
}

// CheckAdminAccess checks if an admin has access to features based on payment status
func (s *PaymentService) CheckAdminAccess(ctx context.Context, adminID int) (bool, error) {
	// First check cache if available (for better performance)

	// Check database
	return s.adminRepo.CheckAdminAccess(ctx, adminID)
}

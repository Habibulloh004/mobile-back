package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"mobilka/internal/models"
	"mobilka/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PaymentHistoryRepository handles database operations for payment history
type PaymentHistoryRepository struct {
	db *pgxpool.Pool
}

// NewPaymentHistoryRepository creates a new payment history repository
func NewPaymentHistoryRepository(db *pgxpool.Pool) *PaymentHistoryRepository {
	return &PaymentHistoryRepository{
		db: db,
	}
}

// Create creates a new payment history record
func (r *PaymentHistoryRepository) Create(ctx context.Context, payment *models.PaymentHistory) error {
	query := `
		INSERT INTO payment_history (
			admin_id, amount, payment_date, payment_method, transaction_id, 
			subscription_tier_id, status, notes
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		payment.AdminID,
		payment.Amount,
		payment.PaymentDate,
		payment.PaymentMethod,
		payment.TransactionID,
		payment.SubscriptionTierID,
		payment.Status,
		payment.Notes,
	).Scan(
		&payment.ID,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	return err
}

// GetByID retrieves a payment history record by ID
func (r *PaymentHistoryRepository) GetByID(ctx context.Context, id int) (*models.PaymentHistory, error) {
	query := `
		SELECT id, admin_id, amount, payment_date, payment_method, transaction_id, 
		       subscription_tier_id, period_start, period_end, status, notes, 
		       verified_by, verified_at, created_at, updated_at
		FROM payment_history
		WHERE id = $1
	`

	var payment models.PaymentHistory
	var subscriptionTierID sql.NullInt32
	var periodStart sql.NullTime
	var periodEnd sql.NullTime
	var verifiedBy sql.NullInt32
	var verifiedAt sql.NullTime

	err := r.db.QueryRow(ctx, query, id).Scan(
		&payment.ID,
		&payment.AdminID,
		&payment.Amount,
		&payment.PaymentDate,
		&payment.PaymentMethod,
		&payment.TransactionID,
		&subscriptionTierID,
		&periodStart,
		&periodEnd,
		&payment.Status,
		&payment.Notes,
		&verifiedBy,
		&verifiedAt,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrResourceNotFound
		}
		return nil, err
	}

	if subscriptionTierID.Valid {
		val := int(subscriptionTierID.Int32)
		payment.SubscriptionTierID = &val
	}

	if periodStart.Valid {
		payment.PeriodStart = &periodStart.Time
	}

	if periodEnd.Valid {
		payment.PeriodEnd = &periodEnd.Time
	}

	if verifiedBy.Valid {
		val := int(verifiedBy.Int32)
		payment.VerifiedBy = &val
	}

	if verifiedAt.Valid {
		payment.VerifiedAt = &verifiedAt.Time
	}

	return &payment, nil
}

// GetByAdminID retrieves all payment history records for a specific admin
func (r *PaymentHistoryRepository) GetByAdminID(ctx context.Context, adminID int) ([]*models.PaymentHistory, error) {
	query := `
		SELECT id, admin_id, amount, payment_date, payment_method, transaction_id, 
		       subscription_tier_id, period_start, period_end, status, notes, 
		       verified_by, verified_at, created_at, updated_at
		FROM payment_history
		WHERE admin_id = $1
		ORDER BY payment_date DESC
	`

	rows, err := r.db.Query(ctx, query, adminID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*models.PaymentHistory
	for rows.Next() {
		var payment models.PaymentHistory
		var subscriptionTierID sql.NullInt32
		var periodStart sql.NullTime
		var periodEnd sql.NullTime
		var verifiedBy sql.NullInt32
		var verifiedAt sql.NullTime

		err := rows.Scan(
			&payment.ID,
			&payment.AdminID,
			&payment.Amount,
			&payment.PaymentDate,
			&payment.PaymentMethod,
			&payment.TransactionID,
			&subscriptionTierID,
			&periodStart,
			&periodEnd,
			&payment.Status,
			&payment.Notes,
			&verifiedBy,
			&verifiedAt,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if subscriptionTierID.Valid {
			val := int(subscriptionTierID.Int32)
			payment.SubscriptionTierID = &val
		}

		if periodStart.Valid {
			payment.PeriodStart = &periodStart.Time
		}

		if periodEnd.Valid {
			payment.PeriodEnd = &periodEnd.Time
		}

		if verifiedBy.Valid {
			val := int(verifiedBy.Int32)
			payment.VerifiedBy = &val
		}

		if verifiedAt.Valid {
			payment.VerifiedAt = &verifiedAt.Time
		}

		payments = append(payments, &payment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}

// GetAll retrieves all payment history records
func (r *PaymentHistoryRepository) GetAll(ctx context.Context) ([]*models.PaymentHistory, error) {
	query := `
		SELECT id, admin_id, amount, payment_date, payment_method, transaction_id, 
		       subscription_tier_id, period_start, period_end, status, notes, 
		       verified_by, verified_at, created_at, updated_at
		FROM payment_history
		ORDER BY payment_date DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*models.PaymentHistory
	for rows.Next() {
		var payment models.PaymentHistory
		var subscriptionTierID sql.NullInt32
		var periodStart sql.NullTime
		var periodEnd sql.NullTime
		var verifiedBy sql.NullInt32
		var verifiedAt sql.NullTime

		err := rows.Scan(
			&payment.ID,
			&payment.AdminID,
			&payment.Amount,
			&payment.PaymentDate,
			&payment.PaymentMethod,
			&payment.TransactionID,
			&subscriptionTierID,
			&periodStart,
			&periodEnd,
			&payment.Status,
			&payment.Notes,
			&verifiedBy,
			&verifiedAt,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if subscriptionTierID.Valid {
			val := int(subscriptionTierID.Int32)
			payment.SubscriptionTierID = &val
		}

		if periodStart.Valid {
			payment.PeriodStart = &periodStart.Time
		}

		if periodEnd.Valid {
			payment.PeriodEnd = &periodEnd.Time
		}

		if verifiedBy.Valid {
			val := int(verifiedBy.Int32)
			payment.VerifiedBy = &val
		}

		if verifiedAt.Valid {
			payment.VerifiedAt = &verifiedAt.Time
		}

		payments = append(payments, &payment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}

// GetPendingPayments retrieves all pending payment history records
func (r *PaymentHistoryRepository) GetPendingPayments(ctx context.Context) ([]*models.PaymentHistory, error) {
	query := `
		SELECT id, admin_id, amount, payment_date, payment_method, transaction_id, 
		       subscription_tier_id, period_start, period_end, status, notes, 
		       verified_by, verified_at, created_at, updated_at
		FROM payment_history
		WHERE status = 'pending'
		ORDER BY payment_date ASC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*models.PaymentHistory
	for rows.Next() {
		var payment models.PaymentHistory
		var subscriptionTierID sql.NullInt32
		var periodStart sql.NullTime
		var periodEnd sql.NullTime
		var verifiedBy sql.NullInt32
		var verifiedAt sql.NullTime

		err := rows.Scan(
			&payment.ID,
			&payment.AdminID,
			&payment.Amount,
			&payment.PaymentDate,
			&payment.PaymentMethod,
			&payment.TransactionID,
			&subscriptionTierID,
			&periodStart,
			&periodEnd,
			&payment.Status,
			&payment.Notes,
			&verifiedBy,
			&verifiedAt,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if subscriptionTierID.Valid {
			val := int(subscriptionTierID.Int32)
			payment.SubscriptionTierID = &val
		}

		if periodStart.Valid {
			payment.PeriodStart = &periodStart.Time
		}

		if periodEnd.Valid {
			payment.PeriodEnd = &periodEnd.Time
		}

		if verifiedBy.Valid {
			val := int(verifiedBy.Int32)
			payment.VerifiedBy = &val
		}

		if verifiedAt.Valid {
			payment.VerifiedAt = &verifiedAt.Time
		}

		payments = append(payments, &payment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}

// VerifyPayment updates a payment record status to verified or rejected
func (r *PaymentHistoryRepository) VerifyPayment(ctx context.Context, id int, superAdminID int, status string, notes string, periodStart *time.Time, periodEnd *time.Time) error {
	query := `
		UPDATE payment_history
		SET status = $2, notes = $3, verified_by = $4, verified_at = CURRENT_TIMESTAMP, 
		    period_start = $5, period_end = $6
		WHERE id = $1
		RETURNING verified_at
	`

	var verifiedAt time.Time
	err := r.db.QueryRow(ctx, query, id, status, notes, superAdminID, periodStart, periodEnd).Scan(&verifiedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.ErrResourceNotFound
		}
		return err
	}

	return nil
}

// UpdateAdminSubscription updates an admin's subscription status based on payment verification
func (r *PaymentHistoryRepository) UpdateAdminSubscription(ctx context.Context, adminID int, subscriptionTierID *int, expiresAt *time.Time, status string, isRestricted bool) error {
	query := `
		UPDATE admin
		SET subscription_tier_id = $2, subscription_expires_at = $3, 
		    subscription_status = $4, is_access_restricted = $5
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, adminID, subscriptionTierID, expiresAt, status, isRestricted)
	return err
}

// GetLatestVerifiedPayment gets the most recent verified payment for an admin
func (r *PaymentHistoryRepository) GetLatestVerifiedPayment(ctx context.Context, adminID int) (*models.PaymentHistory, error) {
	query := `
		SELECT id, admin_id, amount, payment_date, payment_method, transaction_id, 
		       subscription_tier_id, period_start, period_end, status, notes, 
		       verified_by, verified_at, created_at, updated_at
		FROM payment_history
		WHERE admin_id = $1 AND status = 'verified'
		ORDER BY verified_at DESC
		LIMIT 1
	`

	var payment models.PaymentHistory
	var subscriptionTierID sql.NullInt32
	var periodStart sql.NullTime
	var periodEnd sql.NullTime
	var verifiedBy sql.NullInt32
	var verifiedAt sql.NullTime

	err := r.db.QueryRow(ctx, query, adminID).Scan(
		&payment.ID,
		&payment.AdminID,
		&payment.Amount,
		&payment.PaymentDate,
		&payment.PaymentMethod,
		&payment.TransactionID,
		&subscriptionTierID,
		&periodStart,
		&periodEnd,
		&payment.Status,
		&payment.Notes,
		&verifiedBy,
		&verifiedAt,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrResourceNotFound
		}
		return nil, err
	}

	if subscriptionTierID.Valid {
		val := int(subscriptionTierID.Int32)
		payment.SubscriptionTierID = &val
	}

	if periodStart.Valid {
		payment.PeriodStart = &periodStart.Time
	}

	if periodEnd.Valid {
		payment.PeriodEnd = &periodEnd.Time
	}

	if verifiedBy.Valid {
		val := int(verifiedBy.Int32)
		payment.VerifiedBy = &val
	}

	if verifiedAt.Valid {
		payment.VerifiedAt = &verifiedAt.Time
	}

	return &payment, nil
}

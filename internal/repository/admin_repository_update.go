package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"mobilka/internal/models"
	"mobilka/internal/utils"
)

// These are additional methods to be added to the existing AdminRepository

// UpdateSubscriptionStatus updates an admin's subscription status
func (r *AdminRepository) UpdateSubscriptionStatus(ctx context.Context, id int, subscriptionTierID *int, status string, expiresAt *time.Time, isRestricted bool) error {
	query := `
		UPDATE admin
		SET subscription_tier_id = $2, 
		    subscription_status = $3,
		    subscription_expires_at = $4,
		    is_access_restricted = $5
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query, id, subscriptionTierID, status, expiresAt, isRestricted)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return utils.ErrUserNotFound
	}

	return nil
}

// GetAllWithExpiringSubscriptions retrieves all admins with subscriptions about to expire
func (r *AdminRepository) GetAllWithExpiringSubscriptions(ctx context.Context, daysToExpiry int) ([]*models.Admin, error) {
	query := `
		SELECT 
			id, user_name, email, company_name, system_id, system_token, 
			system_token_updated_time, sms_token, sms_token_updated_time, sms_email, 
			sms_password, sms_message, payment_username, payment_password, 
			users, subscription_tier_id, subscription_status, subscription_expires_at,
			is_access_restricted, created_at, updated_at
		FROM admin
		WHERE subscription_status = 'active' 
		  AND subscription_expires_at IS NOT NULL
		  AND subscription_expires_at <= (CURRENT_TIMESTAMP + INTERVAL '${daysToExpiry} days')
		ORDER BY subscription_expires_at ASC
	`
	// .replace("${daysToExpiry}", daysToExpiry.toString())

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var admins []*models.Admin
	for rows.Next() {
		var admin models.Admin
		var subscriptionTierID sql.NullInt32
		var subscriptionExpiresAt sql.NullTime

		err := rows.Scan(
			&admin.ID,
			&admin.UserName,
			&admin.Email,
			&admin.CompanyName,
			&admin.SystemID,
			&admin.SystemToken,
			&admin.SystemTokenUpdatedTime,
			&admin.SmsToken,
			&admin.SmsTokenUpdatedTime,
			&admin.SmsEmail,
			&admin.SmsPassword,
			&admin.SmsMessage,
			&admin.PaymentUsername,
			&admin.PaymentPassword,
			&admin.Users,
			&subscriptionTierID,
			&admin.SubscriptionStatus,
			&subscriptionExpiresAt,
			&admin.IsAccessRestricted,
			&admin.CreatedAt,
			&admin.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if subscriptionTierID.Valid {
			val := int(subscriptionTierID.Int32)
			admin.SubscriptionTierID = &val
		}

		if subscriptionExpiresAt.Valid {
			admin.SubscriptionExpiresAt = &subscriptionExpiresAt.Time
		}

		admins = append(admins, &admin)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return admins, nil
}

// GetActiveCount gets the count of active admins (with active subscriptions)
func (r *AdminRepository) GetActiveCount(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM admin
		WHERE subscription_status = 'active' AND is_access_restricted = false
	`

	var count int
	err := r.db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// ExpireSubscriptions expires all subscriptions that have passed their expiration date
func (r *AdminRepository) ExpireSubscriptions(ctx context.Context) (int, error) {
	query := `
		UPDATE admin
		SET subscription_status = 'expired', is_access_restricted = true
		WHERE subscription_status = 'active' 
		  AND subscription_expires_at IS NOT NULL
		  AND subscription_expires_at < CURRENT_TIMESTAMP
	`

	result, err := r.db.Exec(ctx, query)
	if err != nil {
		return 0, err
	}

	return int(result.RowsAffected()), nil
}

// GetByIDWithSubscriptionInfo retrieves an admin by ID with subscription information
func (r *AdminRepository) GetByIDWithSubscriptionInfo(ctx context.Context, id int) (*models.Admin, *models.SubscriptionTier, error) {
	query := `
		SELECT 
			a.id, a.user_name, a.email, a.company_name, a.system_id, a.system_token, 
			a.system_token_updated_time, a.sms_token, a.sms_token_updated_time, a.sms_email, 
			a.sms_password, a.sms_message, a.payment_username, a.payment_password, 
			a.users, a.subscription_tier_id, a.subscription_status, a.subscription_expires_at,
			a.is_access_restricted, a.created_at, a.updated_at,
			st.id, st.name, st.min_users, st.max_users, st.price, st.description, 
			st.created_at, st.updated_at
		FROM admin a
		LEFT JOIN subscription_tier st ON a.subscription_tier_id = st.id
		WHERE a.id = $1
	`

	var admin models.Admin
	var subscriptionTier models.SubscriptionTier
	var subscriptionTierID sql.NullInt32
	var subscriptionExpiresAt sql.NullTime
	var tierID sql.NullInt32
	var tierName sql.NullString
	var tierMinUsers sql.NullInt32
	var tierMaxUsers sql.NullInt32
	var tierPrice sql.NullFloat64
	var tierDescription sql.NullString
	var tierCreatedAt sql.NullTime
	var tierUpdatedAt sql.NullTime

	err := r.db.QueryRow(ctx, query, id).Scan(
		&admin.ID,
		&admin.UserName,
		&admin.Email,
		&admin.CompanyName,
		&admin.SystemID,
		&admin.SystemToken,
		&admin.SystemTokenUpdatedTime,
		&admin.SmsToken,
		&admin.SmsTokenUpdatedTime,
		&admin.SmsEmail,
		&admin.SmsPassword,
		&admin.SmsMessage,
		&admin.PaymentUsername,
		&admin.PaymentPassword,
		&admin.Users,
		&subscriptionTierID,
		&admin.SubscriptionStatus,
		&subscriptionExpiresAt,
		&admin.IsAccessRestricted,
		&admin.CreatedAt,
		&admin.UpdatedAt,
		&tierID,
		&tierName,
		&tierMinUsers,
		&tierMaxUsers,
		&tierPrice,
		&tierDescription,
		&tierCreatedAt,
		&tierUpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, utils.ErrUserNotFound
		}
		return nil, nil, err
	}

	if subscriptionTierID.Valid {
		val := int(subscriptionTierID.Int32)
		admin.SubscriptionTierID = &val
	}

	if subscriptionExpiresAt.Valid {
		admin.SubscriptionExpiresAt = &subscriptionExpiresAt.Time
	}

	// If there's no subscription tier, return just the admin
	if !tierID.Valid {
		return &admin, nil, nil
	}

	// Otherwise populate the subscription tier
	subscriptionTier.ID = int(tierID.Int32)

	if tierName.Valid {
		subscriptionTier.Name = tierName.String
	}

	if tierMinUsers.Valid {
		subscriptionTier.MinUsers = int(tierMinUsers.Int32)
	}

	if tierMaxUsers.Valid {
		val := int(tierMaxUsers.Int32)
		subscriptionTier.MaxUsers = &val
	}

	if tierPrice.Valid {
		subscriptionTier.Price = tierPrice.Float64
	}

	if tierDescription.Valid {
		subscriptionTier.Description = tierDescription.String
	}

	if tierCreatedAt.Valid {
		subscriptionTier.CreatedAt = tierCreatedAt.Time
	}

	if tierUpdatedAt.Valid {
		subscriptionTier.UpdatedAt = tierUpdatedAt.Time
	}

	return &admin, &subscriptionTier, nil
}

// CheckAdminAccess checks if an admin has valid access based on subscription status
func (r *AdminRepository) CheckAdminAccess(ctx context.Context, adminID int) (bool, error) {
	query := `
		SELECT NOT is_access_restricted
		FROM admin
		WHERE id = $1
	`

	var hasAccess bool
	err := r.db.QueryRow(ctx, query, adminID).Scan(&hasAccess)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, utils.ErrUserNotFound
		}
		return false, err
	}

	return hasAccess, nil
}

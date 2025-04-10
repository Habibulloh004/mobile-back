package repository

import (
	"context"
	"database/sql"
	"errors"

	"mobilka/internal/models"
	"mobilka/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

// FCMTokenRepository handles database operations for FCM tokens
type FCMTokenRepository struct {
	db *pgxpool.Pool
}

// NewFCMTokenRepository creates a new FCM token repository
func NewFCMTokenRepository(db *pgxpool.Pool) *FCMTokenRepository {
	return &FCMTokenRepository{
		db: db,
	}
}

// Create creates a new FCM token
func (r *FCMTokenRepository) Create(ctx context.Context, fcmToken *models.FCMToken) error {
	query := `
		INSERT INTO fcm_token (admin_id, fcm_token)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		fcmToken.AdminID,
		fcmToken.FCMToken,
	).Scan(
		&fcmToken.ID,
		&fcmToken.CreatedAt,
		&fcmToken.UpdatedAt,
	)

	return err
}

// GetByID retrieves an FCM token by ID
func (r *FCMTokenRepository) GetByID(ctx context.Context, id int) (*models.FCMToken, error) {
	query := `
		SELECT id, admin_id, fcm_token, created_at, updated_at
		FROM fcm_token
		WHERE id = $1
	`

	var fcmToken models.FCMToken
	err := r.db.QueryRow(ctx, query, id).Scan(
		&fcmToken.ID,
		&fcmToken.AdminID,
		&fcmToken.FCMToken,
		&fcmToken.CreatedAt,
		&fcmToken.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrResourceNotFound
		}
		return nil, err
	}

	return &fcmToken, nil
}

// GetByToken retrieves an FCM token by the token string
func (r *FCMTokenRepository) GetByToken(ctx context.Context, token string) (*models.FCMToken, error) {
	query := `
		SELECT id, admin_id, fcm_token, created_at, updated_at
		FROM fcm_token
		WHERE fcm_token = $1
	`

	var fcmToken models.FCMToken
	err := r.db.QueryRow(ctx, query, token).Scan(
		&fcmToken.ID,
		&fcmToken.AdminID,
		&fcmToken.FCMToken,
		&fcmToken.CreatedAt,
		&fcmToken.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrResourceNotFound
		}
		return nil, err
	}

	return &fcmToken, nil
}

// GetByAdminID retrieves all FCM tokens for a specific admin
func (r *FCMTokenRepository) GetByAdminID(ctx context.Context, adminID int) ([]*models.FCMToken, error) {
	query := `
		SELECT id, admin_id, fcm_token, created_at, updated_at
		FROM fcm_token
		WHERE admin_id = $1
		ORDER BY id
	`

	rows, err := r.db.Query(ctx, query, adminID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fcmTokens []*models.FCMToken
	for rows.Next() {
		var fcmToken models.FCMToken
		err := rows.Scan(
			&fcmToken.ID,
			&fcmToken.AdminID,
			&fcmToken.FCMToken,
			&fcmToken.CreatedAt,
			&fcmToken.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		fcmTokens = append(fcmTokens, &fcmToken)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return fcmTokens, nil
}

// GetAll retrieves all FCM tokens
func (r *FCMTokenRepository) GetAll(ctx context.Context) ([]*models.FCMToken, error) {
	query := `
		SELECT id, admin_id, fcm_token, created_at, updated_at
		FROM fcm_token
		ORDER BY id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fcmTokens []*models.FCMToken
	for rows.Next() {
		var fcmToken models.FCMToken
		err := rows.Scan(
			&fcmToken.ID,
			&fcmToken.AdminID,
			&fcmToken.FCMToken,
			&fcmToken.CreatedAt,
			&fcmToken.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		fcmTokens = append(fcmTokens, &fcmToken)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return fcmTokens, nil
}

// Delete deletes an FCM token
func (r *FCMTokenRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM fcm_token WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return utils.ErrResourceNotFound
	}

	return nil
}

// DeleteByToken deletes an FCM token by the token string
func (r *FCMTokenRepository) DeleteByToken(ctx context.Context, token string) error {
	query := `DELETE FROM fcm_token WHERE fcm_token = $1`

	result, err := r.db.Exec(ctx, query, token)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return utils.ErrResourceNotFound
	}

	return nil
}

// DeleteByAdminID deletes all FCM tokens for a specific admin
func (r *FCMTokenRepository) DeleteByAdminID(ctx context.Context, adminID int) error {
	query := `DELETE FROM fcm_token WHERE admin_id = $1`

	_, err := r.db.Exec(ctx, query, adminID)
	return err
}

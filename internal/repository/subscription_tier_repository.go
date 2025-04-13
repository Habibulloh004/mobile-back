package repository

import (
	"context"
	"database/sql"
	"errors"

	"mobilka/internal/models"
	"mobilka/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

// SubscriptionTierRepository handles database operations for subscription tiers
type SubscriptionTierRepository struct {
	db *pgxpool.Pool
}

// NewSubscriptionTierRepository creates a new subscription tier repository
func NewSubscriptionTierRepository(db *pgxpool.Pool) *SubscriptionTierRepository {
	return &SubscriptionTierRepository{
		db: db,
	}
}

// Create creates a new subscription tier
func (r *SubscriptionTierRepository) Create(ctx context.Context, tier *models.SubscriptionTier) error {
	query := `
		INSERT INTO subscription_tier (name, min_users, max_users, price, description)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		tier.Name,
		tier.MinUsers,
		tier.MaxUsers,
		tier.Price,
		tier.Description,
	).Scan(
		&tier.ID,
		&tier.CreatedAt,
		&tier.UpdatedAt,
	)

	return err
}

// GetByID retrieves a subscription tier by ID
func (r *SubscriptionTierRepository) GetByID(ctx context.Context, id int) (*models.SubscriptionTier, error) {
	query := `
		SELECT id, name, min_users, max_users, price, description, created_at, updated_at
		FROM subscription_tier
		WHERE id = $1
	`

	var tier models.SubscriptionTier
	var maxUsers sql.NullInt32

	err := r.db.QueryRow(ctx, query, id).Scan(
		&tier.ID,
		&tier.Name,
		&tier.MinUsers,
		&maxUsers,
		&tier.Price,
		&tier.Description,
		&tier.CreatedAt,
		&tier.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrResourceNotFound
		}
		return nil, err
	}

	if maxUsers.Valid {
		val := int(maxUsers.Int32)
		tier.MaxUsers = &val
	}

	return &tier, nil
}

// GetAll retrieves all subscription tiers
func (r *SubscriptionTierRepository) GetAll(ctx context.Context) ([]*models.SubscriptionTier, error) {
	query := `
		SELECT id, name, min_users, max_users, price, description, created_at, updated_at
		FROM subscription_tier
		ORDER BY min_users
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tiers []*models.SubscriptionTier
	for rows.Next() {
		var tier models.SubscriptionTier
		var maxUsers sql.NullInt32

		err := rows.Scan(
			&tier.ID,
			&tier.Name,
			&tier.MinUsers,
			&maxUsers,
			&tier.Price,
			&tier.Description,
			&tier.CreatedAt,
			&tier.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if maxUsers.Valid {
			val := int(maxUsers.Int32)
			tier.MaxUsers = &val
		}

		tiers = append(tiers, &tier)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tiers, nil
}

// Update updates a subscription tier
func (r *SubscriptionTierRepository) Update(ctx context.Context, id int, tier *models.SubscriptionTier) error {
	query := `
		UPDATE subscription_tier
		SET name = $2, min_users = $3, max_users = $4, price = $5, description = $6
		WHERE id = $1
		RETURNING updated_at
	`

	err := r.db.QueryRow(ctx, query,
		id,
		tier.Name,
		tier.MinUsers,
		tier.MaxUsers,
		tier.Price,
		tier.Description,
	).Scan(&tier.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.ErrResourceNotFound
		}
		return err
	}

	return nil
}

// Delete deletes a subscription tier
func (r *SubscriptionTierRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM subscription_tier WHERE id = $1`

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

// GetTierForUserCount retrieves the appropriate subscription tier for a given user count
func (r *SubscriptionTierRepository) GetTierForUserCount(ctx context.Context, userCount int) (*models.SubscriptionTier, error) {
	query := `
		SELECT id, name, min_users, max_users, price, description, created_at, updated_at
		FROM subscription_tier
		WHERE min_users <= $1 AND (max_users IS NULL OR max_users >= $1)
		ORDER BY price DESC
		LIMIT 1
	`

	var tier models.SubscriptionTier
	var maxUsers sql.NullInt32

	err := r.db.QueryRow(ctx, query, userCount).Scan(
		&tier.ID,
		&tier.Name,
		&tier.MinUsers,
		&maxUsers,
		&tier.Price,
		&tier.Description,
		&tier.CreatedAt,
		&tier.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrResourceNotFound
		}
		return nil, err
	}

	if maxUsers.Valid {
		val := int(maxUsers.Int32)
		tier.MaxUsers = &val
	}

	return &tier, nil
}

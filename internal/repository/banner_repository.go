package repository

import (
	"context"
	"database/sql"
	"errors"

	"mobilka/internal/models"
	"mobilka/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

// BannerRepository handles database operations for banners
type BannerRepository struct {
	db *pgxpool.Pool
}

// NewBannerRepository creates a new banner repository
func NewBannerRepository(db *pgxpool.Pool) *BannerRepository {
	return &BannerRepository{
		db: db,
	}
}

// Create creates a new banner
func (r *BannerRepository) Create(ctx context.Context, banner *models.Banner) error {
	query := `
		INSERT INTO banner (admin_id, image, title, body)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		banner.AdminID,
		banner.Image,
		banner.Title,
		banner.Body,
	).Scan(
		&banner.ID,
		&banner.CreatedAt,
		&banner.UpdatedAt,
	)

	return err
}

// GetByID retrieves a banner by ID
func (r *BannerRepository) GetByID(ctx context.Context, id int) (*models.Banner, error) {
	query := `
		SELECT id, admin_id, image, title, body, created_at, updated_at
		FROM banner
		WHERE id = $1
	`

	var banner models.Banner
	err := r.db.QueryRow(ctx, query, id).Scan(
		&banner.ID,
		&banner.AdminID,
		&banner.Image,
		&banner.Title,
		&banner.Body,
		&banner.CreatedAt,
		&banner.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrResourceNotFound
		}
		return nil, err
	}

	return &banner, nil
}

// GetByAdminID retrieves all banners for a specific admin
func (r *BannerRepository) GetByAdminID(ctx context.Context, adminID int) ([]*models.Banner, error) {
	query := `
		SELECT id, admin_id, image, title, body, created_at, updated_at
		FROM banner
		WHERE admin_id = $1
		ORDER BY id
	`

	rows, err := r.db.Query(ctx, query, adminID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var banners []*models.Banner
	for rows.Next() {
		var banner models.Banner
		err := rows.Scan(
			&banner.ID,
			&banner.AdminID,
			&banner.Image,
			&banner.Title,
			&banner.Body,
			&banner.CreatedAt,
			&banner.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		banners = append(banners, &banner)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return banners, nil
}

// GetAll retrieves all banners
func (r *BannerRepository) GetAll(ctx context.Context) ([]*models.Banner, error) {
	query := `
		SELECT id, admin_id, image, title, body, created_at, updated_at
		FROM banner
		ORDER BY id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var banners []*models.Banner
	for rows.Next() {
		var banner models.Banner
		err := rows.Scan(
			&banner.ID,
			&banner.AdminID,
			&banner.Image,
			&banner.Title,
			&banner.Body,
			&banner.CreatedAt,
			&banner.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		banners = append(banners, &banner)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return banners, nil
}

// Update updates a banner
func (r *BannerRepository) Update(ctx context.Context, id int, banner *models.Banner) error {
	query := `
		UPDATE banner
		SET image = $2, title = $3, body = $4
		WHERE id = $1 AND admin_id = $5
		RETURNING updated_at
	`

	err := r.db.QueryRow(ctx, query,
		id,
		banner.Image,
		banner.Title,
		banner.Body,
		banner.AdminID,
	).Scan(&banner.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.ErrResourceNotFound
		}
		return err
	}

	return nil
}

// Delete deletes a banner
func (r *BannerRepository) Delete(ctx context.Context, id int, adminID int) error {
	query := `DELETE FROM banner WHERE id = $1 AND admin_id = $2`

	result, err := r.db.Exec(ctx, query, id, adminID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return utils.ErrResourceNotFound
	}

	return nil
}

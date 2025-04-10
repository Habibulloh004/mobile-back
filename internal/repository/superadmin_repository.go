package repository

import (
	"context"
	"database/sql"
	"errors"

	"mobilka/internal/models"
	"mobilka/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

// SuperAdminRepository handles database operations for the super admin
type SuperAdminRepository struct {
	db *pgxpool.Pool
}

// NewSuperAdminRepository creates a new super admin repository
func NewSuperAdminRepository(db *pgxpool.Pool) *SuperAdminRepository {
	return &SuperAdminRepository{
		db: db,
	}
}

// GetByLogin retrieves a super admin by login
func (r *SuperAdminRepository) GetByLogin(ctx context.Context, login string) (*models.SuperAdmin, error) {
	query := `
		SELECT id, login, password, created_at, updated_at
		FROM super_admin
		WHERE login = $1
	`

	var superAdmin models.SuperAdmin
	err := r.db.QueryRow(ctx, query, login).Scan(
		&superAdmin.ID,
		&superAdmin.Login,
		&superAdmin.Password,
		&superAdmin.CreatedAt,
		&superAdmin.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrUserNotFound
		}
		return nil, err
	}

	return &superAdmin, nil
}

// GetByID retrieves a super admin by ID
func (r *SuperAdminRepository) GetByID(ctx context.Context, id int) (*models.SuperAdmin, error) {
	query := `
		SELECT id, login, password, created_at, updated_at
		FROM super_admin
		WHERE id = $1
	`

	var superAdmin models.SuperAdmin
	err := r.db.QueryRow(ctx, query, id).Scan(
		&superAdmin.ID,
		&superAdmin.Login,
		&superAdmin.Password,
		&superAdmin.CreatedAt,
		&superAdmin.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrUserNotFound
		}
		return nil, err
	}

	return &superAdmin, nil
}

// UpdatePassword updates the super admin password
func (r *SuperAdminRepository) UpdatePassword(ctx context.Context, id int, hashedPassword string) error {
	query := `
		UPDATE super_admin
		SET password = $2
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query, id, hashedPassword)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return utils.ErrUserNotFound
	}

	return nil
}

// SetupDefaultSuperAdmin creates the default super admin account with the provided credentials
func (r *SuperAdminRepository) SetupDefaultSuperAdmin(ctx context.Context, login, hashedPassword string) error {
	query := `
		INSERT INTO super_admin (login, password)
		VALUES ($1, $2)
		ON CONFLICT (login) DO UPDATE
		SET password = $2
	`

	_, err := r.db.Exec(ctx, query, login, hashedPassword)
	return err
}

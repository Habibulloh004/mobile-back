package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"mobilka/internal/models"
	"mobilka/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RestaurantRepository handles database operations for restaurants
type RestaurantRepository struct {
	db *pgxpool.Pool
}

// NewRestaurantRepository creates a new restaurant repository
func NewRestaurantRepository(db *pgxpool.Pool) *RestaurantRepository {
	return &RestaurantRepository{
		db: db,
	}
}

// Create creates a new restaurant or updates an existing one if the admin already has one
func (r *RestaurantRepository) Create(ctx context.Context, restaurant *models.Restaurant) error {
	// First check if a restaurant already exists for this admin
	var exists bool
	var existingID int
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM restaurant WHERE admin_id = $1), 
		COALESCE((SELECT id FROM restaurant WHERE admin_id = $1 LIMIT 1), 0)
	`, restaurant.AdminID).Scan(&exists, &existingID)
	if err != nil {
		return fmt.Errorf("failed to check existing restaurant: %w", err)
	}

	// Convert nested structs to JSON
	contactsJSON, err := json.Marshal(restaurant.Contacts)
	if err != nil {
		return fmt.Errorf("failed to marshal contacts: %w", err)
	}

	socialMediaJSON, err := json.Marshal(restaurant.SocialMedia)
	if err != nil {
		return fmt.Errorf("failed to marshal social media: %w", err)
	}

	// If restaurant already exists, update it instead of creating a new one
	if exists {
		query := `
			UPDATE restaurant
			SET text = $2, contacts = $3, social_media = $4
			WHERE admin_id = $1
			RETURNING id, created_at, updated_at
		`

		err = r.db.QueryRow(ctx, query,
			restaurant.AdminID,
			restaurant.Text,
			contactsJSON,
			socialMediaJSON,
		).Scan(
			&restaurant.ID,
			&restaurant.CreatedAt,
			&restaurant.UpdatedAt,
		)

		if err != nil {
			return fmt.Errorf("failed to update existing restaurant: %w", err)
		}

		return nil
	}

	// If no restaurant exists for this admin, create a new one
	query := `
		INSERT INTO restaurant (admin_id, text, contacts, social_media)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err = r.db.QueryRow(ctx, query,
		restaurant.AdminID,
		restaurant.Text,
		contactsJSON,
		socialMediaJSON,
	).Scan(
		&restaurant.ID,
		&restaurant.CreatedAt,
		&restaurant.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create restaurant: %w", err)
	}

	return nil
}

// GetByID retrieves a restaurant by ID
func (r *RestaurantRepository) GetByID(ctx context.Context, id int) (*models.Restaurant, error) {
	query := `
		SELECT id, admin_id, text, contacts, social_media, created_at, updated_at
		FROM restaurant
		WHERE id = $1
	`

	var restaurant models.Restaurant
	var contactsJSON, socialMediaJSON []byte

	err := r.db.QueryRow(ctx, query, id).Scan(
		&restaurant.ID,
		&restaurant.AdminID,
		&restaurant.Text,
		&contactsJSON,
		&socialMediaJSON,
		&restaurant.CreatedAt,
		&restaurant.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrResourceNotFound
		}
		return nil, err
	}

	// Unmarshal JSON to structs
	if err := json.Unmarshal(contactsJSON, &restaurant.Contacts); err != nil {
		return nil, fmt.Errorf("failed to unmarshal contacts: %w", err)
	}

	if err := json.Unmarshal(socialMediaJSON, &restaurant.SocialMedia); err != nil {
		return nil, fmt.Errorf("failed to unmarshal social media: %w", err)
	}

	return &restaurant, nil
}

// GetByAdminID retrieves all restaurants for a specific admin
func (r *RestaurantRepository) GetByAdminID(ctx context.Context, adminID int) ([]*models.Restaurant, error) {
	query := `
		SELECT id, admin_id, text, contacts, social_media, created_at, updated_at
		FROM restaurant
		WHERE admin_id = $1
		ORDER BY id
	`

	rows, err := r.db.Query(ctx, query, adminID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var restaurants []*models.Restaurant
	for rows.Next() {
		var restaurant models.Restaurant
		var contactsJSON, socialMediaJSON []byte

		err := rows.Scan(
			&restaurant.ID,
			&restaurant.AdminID,
			&restaurant.Text,
			&contactsJSON,
			&socialMediaJSON,
			&restaurant.CreatedAt,
			&restaurant.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Unmarshal JSON to structs
		if err := json.Unmarshal(contactsJSON, &restaurant.Contacts); err != nil {
			return nil, fmt.Errorf("failed to unmarshal contacts: %w", err)
		}

		if err := json.Unmarshal(socialMediaJSON, &restaurant.SocialMedia); err != nil {
			return nil, fmt.Errorf("failed to unmarshal social media: %w", err)
		}

		restaurants = append(restaurants, &restaurant)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return restaurants, nil
}

// GetAll retrieves all restaurants
func (r *RestaurantRepository) GetAll(ctx context.Context) ([]*models.Restaurant, error) {
	query := `
		SELECT id, admin_id, text, contacts, social_media, created_at, updated_at
		FROM restaurant
		ORDER BY id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var restaurants []*models.Restaurant
	for rows.Next() {
		var restaurant models.Restaurant
		var contactsJSON, socialMediaJSON []byte

		err := rows.Scan(
			&restaurant.ID,
			&restaurant.AdminID,
			&restaurant.Text,
			&contactsJSON,
			&socialMediaJSON,
			&restaurant.CreatedAt,
			&restaurant.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Unmarshal JSON to structs
		if err := json.Unmarshal(contactsJSON, &restaurant.Contacts); err != nil {
			return nil, fmt.Errorf("failed to unmarshal contacts: %w", err)
		}

		if err := json.Unmarshal(socialMediaJSON, &restaurant.SocialMedia); err != nil {
			return nil, fmt.Errorf("failed to unmarshal social media: %w", err)
		}

		restaurants = append(restaurants, &restaurant)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return restaurants, nil
}

// Update updates a restaurant
func (r *RestaurantRepository) Update(ctx context.Context, id int, restaurant *models.Restaurant) error {
	// Convert nested structs to JSON
	contactsJSON, err := json.Marshal(restaurant.Contacts)
	if err != nil {
		return fmt.Errorf("failed to marshal contacts: %w", err)
	}

	socialMediaJSON, err := json.Marshal(restaurant.SocialMedia)
	if err != nil {
		return fmt.Errorf("failed to marshal social media: %w", err)
	}

	query := `
		UPDATE restaurant
		SET admin_id = $2, text = $3, contacts = $4, social_media = $5
		WHERE id = $1
		RETURNING updated_at
	`

	err = r.db.QueryRow(ctx, query,
		id,
		restaurant.AdminID,
		restaurant.Text,
		contactsJSON,
		socialMediaJSON,
	).Scan(&restaurant.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.ErrResourceNotFound
		}
		return fmt.Errorf("failed to update restaurant: %w", err)
	}

	return nil
}

// Delete deletes a restaurant
func (r *RestaurantRepository) Delete(ctx context.Context, id int, adminID int) error {
	query := `DELETE FROM restaurant WHERE id = $1 AND admin_id = $2`

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

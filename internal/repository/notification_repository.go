package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"mobilka/internal/models"
	"mobilka/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NotificationRepository handles database operations for notifications
type NotificationRepository struct {
	db *pgxpool.Pool
}

// NewNotificationRepository creates a new notification repository
func NewNotificationRepository(db *pgxpool.Pool) *NotificationRepository {
	return &NotificationRepository{
		db: db,
	}
}

// Create creates a new notification
func (r *NotificationRepository) Create(ctx context.Context, notification *models.Notification) error {
	query := `
		INSERT INTO notification (admin_id, payload, title, body)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		notification.AdminID,
		notification.Payload,
		notification.Title,
		notification.Body,
	).Scan(
		&notification.ID,
		&notification.CreatedAt,
		&notification.UpdatedAt,
	)

	return err
}

// GetByID retrieves a notification by ID
func (r *NotificationRepository) GetByID(ctx context.Context, id int) (*models.Notification, error) {
	query := `
		SELECT id, admin_id, payload, title, body, created_at, updated_at
		FROM notification
		WHERE id = $1
	`

	var notification models.Notification
	err := r.db.QueryRow(ctx, query, id).Scan(
		&notification.ID,
		&notification.AdminID,
		&notification.Payload,
		&notification.Title,
		&notification.Body,
		&notification.CreatedAt,
		&notification.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrResourceNotFound
		}
		return nil, err
	}

	return &notification, nil
}

// GetByAdminID retrieves all notifications for a specific admin
func (r *NotificationRepository) GetByAdminID(ctx context.Context, adminID int) ([]*models.Notification, error) {
	query := `
		SELECT id, admin_id, payload, title, body, created_at, updated_at
		FROM notification
		WHERE admin_id = $1
		ORDER BY id DESC
	`

	rows, err := r.db.Query(ctx, query, adminID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*models.Notification
	for rows.Next() {
		var notification models.Notification
		err := rows.Scan(
			&notification.ID,
			&notification.AdminID,
			&notification.Payload,
			&notification.Title,
			&notification.Body,
			&notification.CreatedAt,
			&notification.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, &notification)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notifications, nil
}

// GetAll retrieves all notifications
func (r *NotificationRepository) GetAll(ctx context.Context) ([]*models.Notification, error) {
	query := `
		SELECT id, admin_id, payload, title, body, created_at, updated_at
		FROM notification
		ORDER BY id DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*models.Notification
	for rows.Next() {
		var notification models.Notification
		err := rows.Scan(
			&notification.ID,
			&notification.AdminID,
			&notification.Payload,
			&notification.Title,
			&notification.Body,
			&notification.CreatedAt,
			&notification.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, &notification)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notifications, nil
}

// Update updates a notification
func (r *NotificationRepository) Update(ctx context.Context, id int, notification *models.Notification) error {
	// Update all fields including admin_id
	query := `
        UPDATE notification
        SET admin_id = $2, payload = $3, title = $4, body = $5
        WHERE id = $1
        RETURNING updated_at
    `

	// Log the update query parameters for debugging
	fmt.Printf("Repository: Updating notification %d with admin_id: %d, title: %s\n",
		id, notification.AdminID, notification.Title)

	err := r.db.QueryRow(ctx, query,
		id,
		notification.AdminID,
		notification.Payload,
		notification.Title,
		notification.Body,
	).Scan(&notification.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.ErrResourceNotFound
		}
		return fmt.Errorf("failed to update notification: %w", err)
	}

	return nil
}

// Delete deletes a notification
func (r *NotificationRepository) Delete(ctx context.Context, id int, adminID int) error {
	query := `DELETE FROM notification WHERE id = $1 AND admin_id = $2`

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

// GetByAdminIDWithPagination retrieves notifications for a specific admin with pagination
func (r *NotificationRepository) GetByAdminIDWithPagination(ctx context.Context, adminID, skip, step int) ([]*models.Notification, error) {
	query := `
		SELECT id, admin_id, payload, title, body, created_at, updated_at
		FROM notification
		WHERE admin_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, adminID, step, skip)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*models.Notification
	for rows.Next() {
		var notification models.Notification
		err := rows.Scan(
			&notification.ID,
			&notification.AdminID,
			&notification.Payload,
			&notification.Title,
			&notification.Body,
			&notification.CreatedAt,
			&notification.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, &notification)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notifications, nil
}

package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"mobilka/internal/models"
	"mobilka/internal/utils"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AdminRepository handles database operations for admins
type AdminRepository struct {
	db *pgxpool.Pool
}

// NewAdminRepository creates a new admin repository
func NewAdminRepository(db *pgxpool.Pool) *AdminRepository {
	return &AdminRepository{
		db: db,
	}
}

// Create creates a new admin
func (r *AdminRepository) Create(ctx context.Context, admin *models.Admin) error {
	// Start a transaction to control ID sequence
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	// Set up a defer to handle transaction rollback if needed
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// First check if admin with this email already exists
	var exists bool
	err = tx.QueryRow(ctx, `
        SELECT EXISTS(SELECT 1 FROM admin WHERE email = $1)
    `, admin.Email).Scan(&exists)

	if err != nil {
		return err
	}

	if exists {
		return utils.NewAppError(utils.ErrResourceAlreadyExists, "Admin with this email already exists", 409)
	}

	// Also check if username
	err = tx.QueryRow(ctx, `
        SELECT EXISTS(SELECT 1 FROM admin WHERE user_name = $1)
    `, admin.UserName).Scan(&exists)

	if err != nil {
		return err
	}

	if exists {
		return utils.NewAppError(utils.ErrResourceAlreadyExists,
			"Admin with this username already exists", 409)
	}

	// Now insert the new admin
	query := `
        INSERT INTO admin (
            user_name, email, company_name, system_id, system_token, 
            system_token_updated_time, sms_token, sms_token_updated_time, sms_email, 
            sms_password, sms_message, payment_username, payment_password, bot_token,
            bot_chat_id, delivery
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
        ) RETURNING id, created_at, updated_at
    `

	err = tx.QueryRow(ctx, query,
		admin.UserName,
		admin.Email,
		admin.CompanyName,
		admin.SystemID,
		admin.SystemToken,
		admin.SystemTokenUpdatedTime,
		admin.SmsToken,
		admin.SmsTokenUpdatedTime,
		admin.SmsEmail,
		admin.SmsPassword,
		admin.SmsMessage,
		admin.PaymentUsername,
		admin.PaymentPassword,
		admin.BotToken,
		admin.BotChatID,
		admin.Delivery,
	).Scan(
		&admin.ID,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)

	if err != nil {
		return err
	}

	// Commit the transaction
	return tx.Commit(ctx)
}

// GetByID retrieves an admin by ID
func (r *AdminRepository) GetByID(ctx context.Context, id int) (*models.Admin, error) {
	query := `
		SELECT 
			id, user_name, email, company_name, system_id, system_token, 
			system_token_updated_time, sms_token, sms_token_updated_time, sms_email, 
			sms_password, sms_message, payment_username, payment_password, bot_token,
			bot_chat_id, delivery, users, created_at, updated_at
		FROM admin
		WHERE id = $1
	`

	var admin models.Admin
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
		&admin.BotToken,
		&admin.BotChatID,
		&admin.Delivery,
		&admin.Users,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrUserNotFound
		}
		return nil, err
	}

	return &admin, nil
}

// GetAll retrieves all admins
func (r *AdminRepository) GetAll(ctx context.Context) ([]*models.Admin, error) {
	query := `
		SELECT 
			id, user_name, email, company_name, system_id, system_token, 
			system_token_updated_time, sms_token, sms_token_updated_time, sms_email, 
			sms_password, sms_message, payment_username, payment_password, bot_token,
			bot_chat_id, delivery, users, created_at, updated_at
		FROM admin
		ORDER BY id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var admins []*models.Admin
	for rows.Next() {
		var admin models.Admin
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
			&admin.BotToken,
			&admin.BotChatID,
			&admin.Delivery,
			&admin.Users,
			&admin.CreatedAt,
			&admin.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		admins = append(admins, &admin)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return admins, nil
}

// GetByEmail retrieves an admin by email
func (r *AdminRepository) GetByEmail(ctx context.Context, email string) (*models.Admin, error) {
	query := `
		SELECT 
			id, user_name, email, company_name, system_id, system_token, 
			system_token_updated_time, sms_token, sms_token_updated_time, sms_email, 
			sms_password, sms_message, payment_username, payment_password, bot_token,
			bot_chat_id, delivery, users, created_at, updated_at
		FROM admin
		WHERE email = $1
	`

	var admin models.Admin
	err := r.db.QueryRow(ctx, query, email).Scan(
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
		&admin.BotToken,
		&admin.BotChatID,
		&admin.Delivery,
		&admin.Users,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrUserNotFound
		}
		return nil, err
	}

	return &admin, nil
}

// GetByUserNameAndSystemID retrieves an admin by username and system ID
func (r *AdminRepository) GetByUserNameAndSystemID(ctx context.Context, userName, systemID string) (*models.Admin, error) {
	query := `
		SELECT 
			id, user_name, email, company_name, system_id, system_token, 
			system_token_updated_time, sms_token, sms_token_updated_time, sms_email, 
			sms_password, sms_message, payment_username, payment_password, bot_token,
			bot_chat_id, delivery, users, created_at, updated_at
		FROM admin
		WHERE user_name = $1 AND system_id = $2
	`

	var admin models.Admin
	err := r.db.QueryRow(ctx, query, userName, systemID).Scan(
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
		&admin.BotToken,
		&admin.BotChatID,
		&admin.Delivery,
		&admin.Users,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrUserNotFound
		}
		return nil, err
	}

	return &admin, nil
}

// GetByCredentials retrieves an admin by username, system ID, and email
func (r *AdminRepository) GetByCredentials(ctx context.Context, userName, systemID, email string) (*models.Admin, error) {
	query := `
		SELECT 
			id, user_name, email, company_name, system_id, system_token, 
			system_token_updated_time, sms_token, sms_token_updated_time, sms_email, 
			sms_password, sms_message, payment_username, payment_password, bot_token,
			bot_chat_id, delivery, users, created_at, updated_at
		FROM admin
		WHERE user_name = $1 AND system_id = $2 AND email = $3
	`

	var admin models.Admin
	err := r.db.QueryRow(ctx, query, userName, systemID, email).Scan(
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
		&admin.BotToken,
		&admin.BotChatID,
		&admin.Delivery,
		&admin.Users,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrUserNotFound
		}
		return nil, err
	}

	return &admin, nil
}

// UpdateSystemToken updates an admin's system token
func (r *AdminRepository) UpdateSystemToken(ctx context.Context, id int, token string) error {
	query := `
        UPDATE admin
        SET 
            system_token = $2,
            system_token_updated_time = CURRENT_TIMESTAMP
        WHERE id = $1
    `

	result, err := r.db.Exec(ctx, query, id, token)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return utils.ErrUserNotFound
	}

	return nil
}

// UpdateSmsToken updates an admin's SMS token
func (r *AdminRepository) UpdateSmsToken(ctx context.Context, id int, token string) error {
	query := `
        UPDATE admin
        SET 
            sms_token = $2,
            sms_token_updated_time = CURRENT_TIMESTAMP
        WHERE id = $1
    `

	result, err := r.db.Exec(ctx, query, id, token)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return utils.ErrUserNotFound
	}

	return nil
}

// Update updates an admin - including delivery and bot fields
func (r *AdminRepository) Update(ctx context.Context, id int, admin *models.Admin) error {
	query := `
        UPDATE admin
        SET 
            user_name = $2,
            email = $3,
            company_name = $4,
            system_id = $5,
            sms_email = $6,
            sms_password = $7,
            sms_message = $8,
            payment_username = $9,
            payment_password = $10,
            bot_token = $11,
            bot_chat_id = $12,
            delivery = $13
        WHERE id = $1
        RETURNING updated_at
    `

	err := r.db.QueryRow(ctx, query,
		id,
		admin.UserName,
		admin.Email,
		admin.CompanyName,
		admin.SystemID,
		admin.SmsEmail,
		admin.SmsPassword,
		admin.SmsMessage,
		admin.PaymentUsername,
		admin.PaymentPassword,
		admin.BotToken,
		admin.BotChatID,
		admin.Delivery,
	).Scan(&admin.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.ErrUserNotFound
		}

		// Check for unique constraint violations
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // unique_violation
				if strings.Contains(pgErr.Message, "admin_email_unique") {
					return utils.NewAppError(utils.ErrResourceAlreadyExists, "Email already exists", 409)
				} else if strings.Contains(pgErr.Message, "admin_username_systemid_unique") {
					return utils.NewAppError(utils.ErrResourceAlreadyExists, "Username and system ID combination already exists", 409)
				}
			}
		}

		return err
	}

	return nil
}

// Delete deletes an admin
func (r *AdminRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM admin WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return utils.ErrUserNotFound
	}

	return nil
}

// IncrementUsersCount increments the users count for an admin
func (r *AdminRepository) IncrementUsersCount(ctx context.Context, id int) error {
	query := `
		UPDATE admin
		SET users = users + 1
		WHERE id = $1
		RETURNING users
	`

	var users int
	err := r.db.QueryRow(ctx, query, id).Scan(&users)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.ErrUserNotFound
		}
		return err
	}

	return nil
}

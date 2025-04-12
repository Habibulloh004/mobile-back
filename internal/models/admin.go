package models

import (
	"time"
)

// Admin model represents the admin entity
type Admin struct {
	ID                     int       `json:"id"`
	UserName               string    `json:"user_name"`
	Email                  string    `json:"email"`
	CompanyName            string    `json:"company_name"`
	SystemID               string    `json:"system_id"`
	SystemToken            string    `json:"system_token"`
	SystemTokenUpdatedTime time.Time `json:"system_token_updated_time"`
	SmsToken               string    `json:"sms_token"`
	SmsTokenUpdatedTime    time.Time `json:"sms_token_updated_time"` // Added field
	SmsEmail               string    `json:"sms_email"`
	SmsPassword            string    `json:"sms_password"` // Password is not exposed in JSON responses
	SmsMessage             string    `json:"sms_message"`
	PaymentUsername        string    `json:"payment_username"`
	PaymentPassword        string    `json:"payment_password"` // Password is not exposed in JSON responses
	Users                  int       `json:"users"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

// AdminCreateRequest represents the creation request for an admin

type AdminCreateRequest struct {
	UserName        string `json:"user_name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	CompanyName     string `json:"company_name" validate:"required"`
	SystemID        string `json:"system_id"`
	SystemToken     string `json:"system_token"`
	SmsToken        string `json:"sms_token"`
	SmsEmail        string `json:"sms_email"`
	SmsPassword     string `json:"sms_password"`
	SmsMessage      string `json:"sms_message"`
	PaymentUsername string `json:"payment_username"`
	PaymentPassword string `json:"payment_password"`
}

// AdminUpdateRequest represents the update request for an admin
type AdminUpdateRequest struct {
	UserName        string `json:"user_name"`
	Email           string `json:"email" validate:"omitempty,email"`
	CompanyName     string `json:"company_name"`
	SystemID        string `json:"system_id"`
	SystemToken     string `json:"system_token"`
	SmsEmail        string `json:"sms_email" validate:"omitempty,email"`
	SmsPassword     string `json:"sms_password"`
	SmsMessage      string `json:"sms_message"`
	SmsToken        string `json:"sms_token"`
	PaymentUsername string `json:"payment_username"`
	PaymentPassword string `json:"payment_password"`
}

// AdminLoginRequest represents the login request for admin
type AdminLoginRequest struct {
	UserName string `json:"user_name" validate:"required"`
	SystemID string `json:"system_id" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

// AdminResponse represents the response for admin without sensitive data
type AdminResponse struct {
	ID                     int       `json:"id"`
	UserName               string    `json:"user_name"`
	Email                  string    `json:"email"`
	CompanyName            string    `json:"company_name"`
	SystemID               string    `json:"system_id"`
	SystemToken            string    `json:"system_token"`
	SystemTokenUpdatedTime time.Time `json:"system_token_updated_time"`
	SmsToken               string    `json:"sms_token"`
	SmsEmail               string    `json:"sms_email"`
	SmsMessage             string    `json:"sms_message"`
	SmsPassword            string    `json:"sms_password"`
	SmsTokenUpdatedTime    time.Time `json:"sms_token_updated_time"`
	PaymentUsername        string    `json:"payment_username"`
	PaymentPassword        string    `json:"payment_password"`
	Users                  int       `json:"users"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

// ToResponse converts Admin model to AdminResponse
func (a *Admin) ToResponse() AdminResponse {
	return AdminResponse{
		ID:                     a.ID,
		UserName:               a.UserName,
		Email:                  a.Email,
		CompanyName:            a.CompanyName,
		SystemID:               a.SystemID,
		SystemToken:            a.SystemToken,
		SystemTokenUpdatedTime: a.SystemTokenUpdatedTime,
		SmsToken:               a.SmsToken,
		SmsEmail:               a.SmsEmail,
		SmsMessage:             a.SmsMessage,
		SmsPassword:            a.SmsPassword,
		SmsTokenUpdatedTime:    a.SmsTokenUpdatedTime,
		PaymentUsername:        a.PaymentUsername,
		PaymentPassword:        a.PaymentPassword,
		Users:                  a.Users,
		CreatedAt:              a.CreatedAt,
		UpdatedAt:              a.UpdatedAt,
	}
}

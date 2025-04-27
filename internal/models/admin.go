package models

import (
	"time"
)

// Admin model represents the admin entity
type Admin struct {
	ID                     int        `json:"id"`
	UserName               string     `json:"user_name"`
	Email                  string     `json:"email"`
	CompanyName            string     `json:"company_name"`
	SystemID               string     `json:"system_id"`
	SystemToken            string     `json:"system_token"`
	SystemTokenUpdatedTime time.Time  `json:"system_token_updated_time"`
	SmsToken               string     `json:"sms_token"`
	SmsTokenUpdatedTime    time.Time  `json:"sms_token_updated_time"` // Added field
	SmsEmail               string     `json:"sms_email"`
	SmsPassword            string     `json:"sms_password"` // Password is not exposed in JSON responses
	SmsMessage             string     `json:"sms_message"`
	PaymentUsername        string     `json:"payment_username"`
	PaymentPassword        string     `json:"payment_password"` // Password is not exposed in JSON responses
	Delivery               int        `json:"delivery"`
	Users                  int        `json:"users"`
	SubscriptionTierID     *int       `json:"subscription_tier_id"`
	SubscriptionStatus     string     `json:"subscription_status"`
	SubscriptionExpiresAt  *time.Time `json:"subscription_expires_at"`
	IsAccessRestricted     bool       `json:"is_access_restricted"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
}

// AdminCreateRequest represents the creation request for an admin

type AdminCreateRequest struct {
	UserName           string `json:"user_name" validate:"required"`
	Email              string `json:"email" validate:"required,email"`
	CompanyName        string `json:"company_name" validate:"required"`
	Delivery           int    `json:"delivery"`
	SystemID           string `json:"system_id"`
	SystemToken        string `json:"system_token"`
	SmsToken           string `json:"sms_token"`
	SmsEmail           string `json:"sms_email"`
	SmsPassword        string `json:"sms_password"`
	SmsMessage         string `json:"sms_message"`
	PaymentUsername    string `json:"payment_username"`
	PaymentPassword    string `json:"payment_password"`
	SubscriptionTierID *int   `json:"subscription_tier_id"`
}

// AdminUpdateRequest represents the update request for an admin
type AdminUpdateRequest struct {
	UserName           string `json:"user_name"`
	Email              string `json:"email" validate:"omitempty,email"`
	Delivery           int    `json:"delivery"`
	CompanyName        string `json:"company_name"`
	SystemID           string `json:"system_id"`
	SystemToken        string `json:"system_token"`
	SmsEmail           string `json:"sms_email" validate:"omitempty,email"`
	SmsPassword        string `json:"sms_password"`
	SmsMessage         string `json:"sms_message"`
	SmsToken           string `json:"sms_token"`
	PaymentUsername    string `json:"payment_username"`
	PaymentPassword    string `json:"payment_password"`
	AdminID            int    `json:"admin_id,omitempty"`
	SubscriptionTierID *int   `json:"subscription_tier_id"`
}

// AdminLoginRequest represents the login request for admin
type AdminLoginRequest struct {
	UserName string `json:"user_name" validate:"required"`
	SystemID string `json:"system_id" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

// AdminResponse represents the response for admin without sensitive data
type AdminResponse struct {
	ID                     int        `json:"id"`
	UserName               string     `json:"user_name"`
	Email                  string     `json:"email"`
	CompanyName            string     `json:"company_name"`
	Delivery               int        `json:"delivery"`
	SystemID               string     `json:"system_id"`
	SystemToken            string     `json:"system_token"`
	SystemTokenUpdatedTime time.Time  `json:"system_token_updated_time"`
	SmsToken               string     `json:"sms_token"`
	SmsEmail               string     `json:"sms_email"`
	SmsMessage             string     `json:"sms_message"`
	SmsPassword            string     `json:"sms_password"`
	SmsTokenUpdatedTime    time.Time  `json:"sms_token_updated_time"`
	PaymentUsername        string     `json:"payment_username"`
	PaymentPassword        string     `json:"payment_password"`
	Users                  int        `json:"users"`
	SubscriptionTierID     *int       `json:"subscription_tier_id"`
	SubscriptionTierName   string     `json:"subscription_tier_name,omitempty"`
	SubscriptionStatus     string     `json:"subscription_status"`
	SubscriptionExpiresAt  *time.Time `json:"subscription_expires_at"`
	IsAccessRestricted     bool       `json:"is_access_restricted"`
	MonthlySubscriptionFee float64    `json:"monthly_subscription_fee"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
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
		Delivery:               a.Delivery,
		SystemTokenUpdatedTime: a.SystemTokenUpdatedTime,
		SmsToken:               a.SmsToken,
		SmsEmail:               a.SmsEmail,
		SmsMessage:             a.SmsMessage,
		SmsPassword:            a.SmsPassword,
		SmsTokenUpdatedTime:    a.SmsTokenUpdatedTime,
		PaymentUsername:        a.PaymentUsername,
		PaymentPassword:        a.PaymentPassword,
		Users:                  a.Users,
		SubscriptionTierID:     a.SubscriptionTierID,
		SubscriptionStatus:     a.SubscriptionStatus,
		SubscriptionExpiresAt:  a.SubscriptionExpiresAt,
		IsAccessRestricted:     a.IsAccessRestricted,
		CreatedAt:              a.CreatedAt,
		UpdatedAt:              a.UpdatedAt,
	}
}

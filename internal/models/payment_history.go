package models

import (
	"time"
)

// PaymentHistory represents a payment record for an admin
type PaymentHistory struct {
	ID                 int        `json:"id"`
	AdminID            int        `json:"admin_id"`
	Amount             float64    `json:"amount"`
	PaymentDate        time.Time  `json:"payment_date"`
	PaymentMethod      string     `json:"payment_method"`
	TransactionID      string     `json:"transaction_id"`
	SubscriptionTierID *int       `json:"subscription_tier_id"`
	PeriodStart        *time.Time `json:"period_start"`
	PeriodEnd          *time.Time `json:"period_end"`
	Status             string     `json:"status"` // pending, verified, rejected
	Notes              string     `json:"notes"`
	VerifiedBy         *int       `json:"verified_by"`
	VerifiedAt         *time.Time `json:"verified_at"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// PaymentCreateRequest represents the request to record a payment
type PaymentCreateRequest struct {
	Amount        float64 `json:"amount" validate:"required,min=0"`
	PaymentMethod string  `json:"payment_method" validate:"required"`
	TransactionID string  `json:"transaction_id"`
	Notes         string  `json:"notes"`
}

// PaymentVerifyRequest represents the request to verify a payment
type PaymentVerifyRequest struct {
	Status      string     `json:"status" validate:"required,oneof=verified rejected"`
	Notes       string     `json:"notes"`
	PeriodStart *time.Time `json:"period_start"`
	PeriodEnd   *time.Time `json:"period_end"`
}

// PaymentHistoryResponse represents the response for a payment record
type PaymentHistoryResponse struct {
	ID                   int        `json:"id"`
	AdminID              int        `json:"admin_id"`
	AdminName            string     `json:"admin_name,omitempty"` // Added for convenience in listing
	Amount               float64    `json:"amount"`
	PaymentDate          time.Time  `json:"payment_date"`
	PaymentMethod        string     `json:"payment_method"`
	TransactionID        string     `json:"transaction_id"`
	SubscriptionTierID   *int       `json:"subscription_tier_id"`
	SubscriptionTierName string     `json:"subscription_tier_name,omitempty"`
	PeriodStart          *time.Time `json:"period_start"`
	PeriodEnd            *time.Time `json:"period_end"`
	Status               string     `json:"status"`
	Notes                string     `json:"notes"`
	VerifiedBy           *int       `json:"verified_by"`
	VerifiedByName       string     `json:"verified_by_name,omitempty"`
	VerifiedAt           *time.Time `json:"verified_at"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

// ToResponse converts PaymentHistory to PaymentHistoryResponse
func (p *PaymentHistory) ToResponse() PaymentHistoryResponse {
	return PaymentHistoryResponse{
		ID:                 p.ID,
		AdminID:            p.AdminID,
		Amount:             p.Amount,
		PaymentDate:        p.PaymentDate,
		PaymentMethod:      p.PaymentMethod,
		TransactionID:      p.TransactionID,
		SubscriptionTierID: p.SubscriptionTierID,
		PeriodStart:        p.PeriodStart,
		PeriodEnd:          p.PeriodEnd,
		Status:             p.Status,
		Notes:              p.Notes,
		VerifiedBy:         p.VerifiedBy,
		VerifiedAt:         p.VerifiedAt,
		CreatedAt:          p.CreatedAt,
		UpdatedAt:          p.UpdatedAt,
	}
}

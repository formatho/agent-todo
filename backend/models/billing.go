package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserSubscription represents a user's subscription to a pricing tier
type UserSubscription struct {
	ID                   string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID               string         `gorm:"type:varchar(36);not null;index" json:"user_id"`
	Tier                 string         `gorm:"type:varchar(50);not null;default:'free'" json:"tier"` // free, starter, pro, enterprise
	Status               string         `gorm:"type:varchar(50);not null;default:'active'" json:"status"` // active, pending, canceled, past_due
	StripeCustomerID     string         `gorm:"type:varchar(255)" json:"stripe_customer_id,omitempty"`
	StripeSubscriptionID string         `gorm:"type:varchar(255)" json:"stripe_subscription_id,omitempty"`
	StripeSessionID      string         `gorm:"type:varchar(255)" json:"stripe_session_id,omitempty"`
	AgentsLimit          int            `gorm:"not null;default:3" json:"agents_limit"`
	OrganizationsLimit   int            `gorm:"not null;default:1" json:"organizations_limit"`
	CurrentPeriodStart   time.Time      `json:"current_period_start,omitempty"`
	CurrentPeriodEnd     time.Time      `json:"current_period_end,omitempty"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// BeforeCreate generates UUID before creating record
func (s *UserSubscription) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}

// PaymentHistory represents a record of payments
type PaymentHistory struct {
	ID                   string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID               string         `gorm:"type:varchar(36);not null;index" json:"user_id"`
	SubscriptionID       string         `gorm:"type:varchar(36);not null;index" json:"subscription_id"`
	StripePaymentIntentID string        `gorm:"type:varchar(255)" json:"stripe_payment_intent_id,omitempty"`
	Amount               float64        `gorm:"not null" json:"amount"`
	Currency             string         `gorm:"type:varchar(10);not null;default:'usd'" json:"currency"`
	Status               string         `gorm:"type:varchar(50);not null" json:"status"` // succeeded, failed, pending
	Description          string         `gorm:"type:text" json:"description,omitempty"`
	CreatedAt            time.Time      `json:"created_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// BeforeCreate generates UUID before creating record
func (p *PaymentHistory) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

// UsageRecord tracks feature usage for billing purposes
type UsageRecord struct {
	ID        string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID    string    `gorm:"type:varchar(36);not null;index" json:"user_id"`
	Feature   string    `gorm:"type:varchar(100);not null" json:"feature"` // e.g., "agents_created", "tasks_executed"
	Count     int       `gorm:"not null;default:1" json:"count"`
	Period    string    `gorm:"type:varchar(20);not null;index" json:"period"` // e.g., "2026-03"
	CreatedAt time.Time `json:"created_at"`
}

// BeforeCreate generates UUID before creating record
func (u *UsageRecord) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

// SubscriptionFeature represents features available per tier
type SubscriptionFeature struct {
	ID          string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Tier        string    `gorm:"type:varchar(50);not null;index" json:"tier"` // free, starter, pro, enterprise
	Feature     string    `gorm:"type:varchar(100);not null" json:"feature"`
	Description string    `gorm:"type:text" json:"description"`
	Enabled     bool      `gorm:"default:true" json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BeforeCreate generates UUID before creating record
func (f *SubscriptionFeature) BeforeCreate(tx *gorm.DB) error {
	if f.ID == "" {
		f.ID = uuid.New().String()
	}
	return nil
}

package models

import (
	"time"

	"github.com/google/uuid"
)

// SubscriptionPlan represents the subscription plan type
type SubscriptionPlan string

const (
	SubscriptionPlanFree       SubscriptionPlan = "free"
	SubscriptionPlanTrial      SubscriptionPlan = "trial"
	SubscriptionPlanStarter    SubscriptionPlan = "starter"
	SubscriptionPlanPro        SubscriptionPlan = "pro"
	SubscriptionPlanEnterprise SubscriptionPlan = "enterprise"
)

// SubscriptionStatus represents the status of a subscription
type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "active"
	SubscriptionStatusTrialing  SubscriptionStatus = "trialing"
	SubscriptionStatusPastDue   SubscriptionStatus = "past_due"
	SubscriptionStatusCanceled  SubscriptionStatus = "canceled"
	SubscriptionStatusExpired   SubscriptionStatus = "expired"
)

// Subscription represents an organisation's subscription
type Subscription struct {
	Base
	OrganisationID     uuid.UUID          `gorm:"type:uuid;uniqueIndex;not null" json:"organisation_id"`
	Organisation       *Organisation      `gorm:"foreignKey:OrganisationID" json:"organisation,omitempty"`
	Plan               SubscriptionPlan   `gorm:"not null;default:'free'" json:"plan"`
	Status             SubscriptionStatus `gorm:"not null;default:'active'" json:"status"`
	TrialStartsAt      *time.Time         `json:"trial_starts_at"`
	TrialEndsAt        *time.Time         `json:"trial_ends_at"`
	CurrentPeriodStart *time.Time         `json:"current_period_start"`
	CurrentPeriodEnd   *time.Time         `json:"current_period_end"`
	CanceledAt         *time.Time         `json:"canceled_at"`
	StripeCustomerID   string             `json:"stripe_customer_id"`
	StripeSubscriptionID string           `json:"stripe_subscription_id"`
}

// EmailTemplateType represents the type of email template
type EmailTemplateType string

const (
	EmailTemplateTypeWelcome       EmailTemplateType = "welcome"
	EmailTemplateTypeValueTips     EmailTemplateType = "value_tips"
	EmailTemplateTypeCaseStudy     EmailTemplateType = "case_study"
	EmailTemplateTypeLimitedOffer  EmailTemplateType = "limited_offer"
	EmailTemplateTypeUpgradeReminder EmailTemplateType = "upgrade_reminder"
	EmailTemplateTypeCustom        EmailTemplateType = "custom"
)

// EmailTemplate represents a reusable email template
type EmailTemplate struct {
	Base
	Name         string            `gorm:"not null" json:"name"`
	Subject      string            `gorm:"not null" json:"subject"`
	TemplateType EmailTemplateType `gorm:"not null" json:"template_type"`
	BodyHTML     string            `gorm:"not null" json:"body_html"`
	BodyText     string            `gorm:"not null" json:"body_text"`
	Variables    string            `json:"variables"` // JSON array of variable names
	IsActive     bool              `gorm:"default:true" json:"is_active"`
}

// EmailSequenceStatus represents the status of an email sequence
type EmailSequenceStatus string

const (
	EmailSequenceStatusActive   EmailSequenceStatus = "active"
	EmailSequenceStatusPaused   EmailSequenceStatus = "paused"
	EmailSequenceStatusArchived EmailSequenceStatus = "archived"
)

// EmailSequence represents a nurture email sequence
type EmailSequence struct {
	Base
	Name        string               `gorm:"not null" json:"name"`
	Description string               `json:"description"`
	Trigger     string               `gorm:"not null" json:"trigger"` // e.g., "trial_started", "user_signup"
	Status      EmailSequenceStatus  `gorm:"not null;default:'active'" json:"status"`
	Steps       []EmailSequenceStep  `gorm:"foreignKey:SequenceID" json:"steps,omitempty"`
}

// EmailSequenceStep represents a single email in a sequence
type EmailSequenceStep struct {
	Base
	SequenceID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"sequence_id"`
	Sequence     *EmailSequence `gorm:"foreignKey:SequenceID" json:"sequence,omitempty"`
	TemplateID   uuid.UUID      `gorm:"type:uuid;not null" json:"template_id"`
	Template     *EmailTemplate `gorm:"foreignKey:TemplateID" json:"template,omitempty"`
	DelayDays    int            `gorm:"not null;default:0" json:"delay_days"` // Days after sequence trigger
	Order        int            `gorm:"not null;default:0" json:"order"`
}

// EmailQueueStatus represents the status of a queued email
type EmailQueueStatus string

const (
	EmailQueueStatusPending   EmailQueueStatus = "pending"
	EmailQueueStatusProcessing EmailQueueStatus = "processing"
	EmailQueueStatusSent      EmailQueueStatus = "sent"
	EmailQueueStatusFailed    EmailQueueStatus = "failed"
)

// EmailQueue represents an email waiting to be sent
type EmailQueue struct {
	Base
	SequenceID      *uuid.UUID       `gorm:"type:uuid;index" json:"sequence_id"`
	Sequence        *EmailSequence   `gorm:"foreignKey:SequenceID" json:"sequence,omitempty"`
	SequenceStepID  *uuid.UUID       `gorm:"type:uuid;index" json:"sequence_step_id"`
	SequenceStep    *EmailSequenceStep `gorm:"foreignKey:SequenceStepID" json:"sequence_step,omitempty"`
	UserID          uuid.UUID        `gorm:"type:uuid;not null;index" json:"user_id"`
	User            *User            `gorm:"foreignKey:UserID" json:"user,omitempty"`
	OrganisationID  uuid.UUID        `gorm:"type:uuid;not null;index" json:"organisation_id"`
	Organisation    *Organisation    `gorm:"foreignKey:OrganisationID" json:"organisation,omitempty"`
	To              string           `gorm:"not null" json:"to"`
	Subject         string           `gorm:"not null" json:"subject"`
	BodyHTML        string           `gorm:"not null" json:"body_html"`
	BodyText        string           `gorm:"not null" json:"body_text"`
	Status          EmailQueueStatus `gorm:"not null;default:'pending';index" json:"status"`
	ScheduledAt     time.Time        `gorm:"not null;index" json:"scheduled_at"`
	SentAt          *time.Time       `json:"sent_at"`
	FailedAt        *time.Time       `json:"failed_at"`
	ErrorMessage    string           `json:"error_message"`
	RetryCount      int              `gorm:"default:0" json:"retry_count"`
	MaxRetries      int              `gorm:"default:3" json:"max_retries"`
}

// EmailLog represents a sent email record
type EmailLog struct {
	Base
	QueueID        *uuid.UUID    `gorm:"type:uuid;index" json:"queue_id"`
	Queue          *EmailQueue   `gorm:"foreignKey:QueueID" json:"queue,omitempty"`
	UserID         uuid.UUID     `gorm:"type:uuid;not null;index" json:"user_id"`
	User           *User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	OrganisationID uuid.UUID     `gorm:"type:uuid;not null;index" json:"organisation_id"`
	Organisation   *Organisation `gorm:"foreignKey:OrganisationID" json:"organisation,omitempty"`
	To             string        `gorm:"not null" json:"to"`
	Subject        string        `gorm:"not null" json:"subject"`
	BodyHTML       string        `json:"body_html"`
	BodyText       string        `json:"body_text"`
	SentAt         time.Time     `gorm:"not null;index" json:"sent_at"`
	OpenedAt       *time.Time    `json:"opened_at"`
	ClickedAt      *time.Time    `json:"clicked_at"`
	MessageID      string        `json:"message_id"` // SMTP message ID for tracking
}

package models

// FeedbackType represents the type of feedback
type FeedbackType string

const (
	FeedbackTypeBug            FeedbackType = "bug"
	FeedbackTypeFeatureRequest FeedbackType = "feature_request"
	FeedbackTypeGeneral        FeedbackType = "general"
	FeedbackTypeImprovement    FeedbackType = "improvement"
)

// FeedbackStatus represents the status of feedback
type FeedbackStatus string

const (
	FeedbackStatusNew         FeedbackStatus = "new"
	FeedbackStatusAcknowledged FeedbackStatus = "acknowledged"
	FeedbackStatusInProgress  FeedbackStatus = "in_progress"
	FeedbackStatusResolved    FeedbackStatus = "resolved"
	FeedbackStatusClosed      FeedbackStatus = "closed"
)

// BetaFeedback represents feedback from beta testers
type BetaFeedback struct {
	Base
	TesterEmail  string         `gorm:"type:varchar(255)" json:"tester_email"`
	TesterName   string         `gorm:"type:varchar(255)" json:"tester_name"`
	FeedbackType FeedbackType   `gorm:"type:varchar(50);not null" json:"feedback_type"`
	Title        string         `gorm:"type:varchar(500);not null" json:"title"`
	Description  string         `gorm:"type:text;not null" json:"description"`
	Priority     TaskPriority   `gorm:"type:varchar(20)" json:"priority"`
	Page         string         `gorm:"type:varchar(500)" json:"page"`
	UserAgent    string         `gorm:"type:text" json:"user_agent"`
	Status       FeedbackStatus `gorm:"type:varchar(20);not null;default:'new'" json:"status"`
	AdminNotes   string         `gorm:"type:text" json:"admin_notes"`
	Rating       int            `gorm:"type:integer" json:"rating"`
}

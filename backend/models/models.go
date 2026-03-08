package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base contains common fields for all models
type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate hook for GORM
func (b *Base) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// User represents a human user
type User struct {
	Base
	Email        string `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"not null" json:"-"`
}

// Agent represents an AI agent
type Agent struct {
	Base
	Name        string `gorm:"not null" json:"name"`
	APIKey      string `gorm:"uniqueIndex;not null" json:"api_key"`
	Description string `json:"description"`
}

// ProjectStatus represents the status of a project
type ProjectStatus string

const (
	ProjectStatusActive    ProjectStatus = "active"
	ProjectStatusArchived  ProjectStatus = "archived"
	ProjectStatusCompleted ProjectStatus = "completed"
)

// Project represents a grouping of tasks
type Project struct {
	Base
	Name            string        `gorm:"not null" json:"name"`
	Description     string        `json:"description"`
	Status          ProjectStatus `gorm:"not null;default:'active'" json:"status"`
	CreatedByUserID uuid.UUID     `gorm:"type:uuid;not null" json:"created_by_user_id"`
	CreatedBy       *User         `gorm:"foreignKey:CreatedByUserID" json:"created_by,omitempty"`
	Tasks           []Task        `gorm:"foreignKey:ProjectID" json:"tasks,omitempty"`
}

// TaskStatus represents the status of a task
type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusFailed     TaskStatus = "failed"
)

// TaskPriority represents the priority of a task
type TaskPriority string

const (
	TaskPriorityLow      TaskPriority = "low"
	TaskPriorityMedium   TaskPriority = "medium"
	TaskPriorityHigh     TaskPriority = "high"
	TaskPriorityCritical TaskPriority = "critical"
)

// Task represents a todo item
type Task struct {
	Base
	Title           string        `gorm:"not null" json:"title"`
	Description     string        `json:"description"`
	Status          TaskStatus    `gorm:"not null;default:'pending'" json:"status"`
	Priority        TaskPriority  `gorm:"not null;default:'medium'" json:"priority"`
	DueDate         *time.Time    `json:"due_date"`
	ProjectID       *uuid.UUID    `gorm:"type:uuid" json:"project_id"`
	CreatedByUserID uuid.UUID     `gorm:"type:uuid;not null" json:"created_by_user_id"`
	AssignedAgentID *uuid.UUID    `gorm:"type:uuid" json:"assigned_agent_id"`
	Project         *Project      `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	CreatedBy       *User         `gorm:"foreignKey:CreatedByUserID" json:"created_by,omitempty"`
	AssignedAgent   *Agent        `gorm:"foreignKey:AssignedAgentID" json:"assigned_agent,omitempty"`
	Comments        []TaskComment `gorm:"foreignKey:TaskID" json:"comments,omitempty"`
	Events          []TaskEvent   `gorm:"foreignKey:TaskID" json:"events,omitempty"`
}

// TaskEventType represents the type of event
type TaskEventType string

const (
	TaskEventCreated       TaskEventType = "created"
	TaskEventUpdated       TaskEventType = "updated"
	TaskEventStatusChanged TaskEventType = "status_changed"
	TaskEventAssigned      TaskEventType = "assigned"
	TaskEventUnassigned    TaskEventType = "unassigned"
)

// TaskEvent represents an audit log entry for a task
type TaskEvent struct {
	Base
	TaskID        uuid.UUID     `gorm:"type:uuid;not null" json:"task_id"`
	Task          *Task         `gorm:"foreignKey:TaskID" json:"task,omitempty"`
	EventType     TaskEventType `gorm:"not null" json:"event_type"`
	PreviousState string        `json:"previous_state"`
	NewState      string        `json:"new_state"`
	ChangedBy     string        `json:"changed_by"` // Can be user email or agent name
}

// TaskComment represents a comment on a task
type TaskComment struct {
	Base
	TaskID     uuid.UUID `gorm:"type:uuid;not null" json:"task_id"`
	Task       *Task     `gorm:"foreignKey:TaskID" json:"task,omitempty"`
	Content    string    `gorm:"not null" json:"content"`
	AuthorID   uuid.UUID `gorm:"type:uuid;not null" json:"author_id"`
	AuthorType string    `gorm:"not null" json:"author_type"` // "user" or "agent"
	AuthorName string    `gorm:"not null" json:"author_name"`
}

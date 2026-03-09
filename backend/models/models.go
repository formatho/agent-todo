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
	Email        string        `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string        `gorm:"not null" json:"-"`
	CurrentOrgID *uuid.UUID    `gorm:"type:uuid;index" json:"current_org_id"`
	CurrentOrg   *Organisation `gorm:"foreignKey:CurrentOrgID" json:"current_org,omitempty"`
}

// AgentRole represents the permission level of an agent
type AgentRole string

const (
	AgentRoleRegular    AgentRole = "regular"    // Can only update own tasks
	AgentRoleSupervisor AgentRole = "supervisor" // Can update any task, create agents
	AgentRoleAdmin      AgentRole = "admin"      // Full permissions
)

// Agent represents an AI agent
type Agent struct {
	Base
	Name           string        `gorm:"not null" json:"name"`
	APIKey         string        `gorm:"uniqueIndex;not null" json:"api_key"`
	Description    string        `json:"description"`
	Role           AgentRole     `gorm:"not null;default:'regular'" json:"role"`
	Enabled        bool          `gorm:"not null;default:true" json:"enabled"`
	OrganisationID *uuid.UUID    `gorm:"type:uuid;index" json:"organisation_id"`
	Organisation   *Organisation `gorm:"foreignKey:OrganisationID" json:"organisation,omitempty"`
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
	Name             string        `gorm:"not null" json:"name"`
	Description      string        `json:"description"`
	Status           ProjectStatus `gorm:"not null;default:'active'" json:"status"`
	RepositoryURL    string        `json:"repository_url"`    // GitHub/GitLab URL
	DeployedURL      string        `json:"deployed_url"`      // Production/staging URL
	DocumentationURL string        `json:"documentation_url"` // Docs URL
	LLMContext       string        `json:"llm_context"`       // Instructions, guidelines, goals for AI agents
	CreatedByUserID  uuid.UUID     `gorm:"type:uuid;not null" json:"created_by_user_id"`
	CreatedBy        *User         `gorm:"foreignKey:CreatedByUserID" json:"created_by,omitempty"`
	OrganisationID   *uuid.UUID    `gorm:"type:uuid;index" json:"organisation_id"`
	Organisation     *Organisation `gorm:"foreignKey:OrganisationID" json:"organisation,omitempty"`
	Tasks            []Task        `gorm:"foreignKey:ProjectID" json:"tasks,omitempty"`
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
	Title            string         `gorm:"not null" json:"title"`
	Description      string         `json:"description"`
	Status           TaskStatus     `gorm:"not null;default:'pending'" json:"status"`
	Priority         TaskPriority   `gorm:"not null;default:'medium'" json:"priority"`
	DueDate          *time.Time     `json:"due_date"`
	CommitURL        string         `json:"commit_url"`
	ProjectID        *uuid.UUID     `gorm:"type:uuid" json:"project_id"`
	OrganisationID   *uuid.UUID     `gorm:"type:uuid;index" json:"organisation_id"`
	CreatedByUserID  *uuid.UUID     `gorm:"type:uuid" json:"created_by_user_id"`  // Nullable for agent-created tasks
	CreatedByAgentID *uuid.UUID     `gorm:"type:uuid" json:"created_by_agent_id"` // Nullable for user-created tasks
	AssignedAgentID  *uuid.UUID     `gorm:"type:uuid" json:"assigned_agent_id"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Project          *Project       `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Organisation     *Organisation  `gorm:"foreignKey:OrganisationID" json:"organisation,omitempty"`
	CreatedBy        *User          `gorm:"foreignKey:CreatedByUserID" json:"created_by,omitempty"`
	CreatedByAgent   *Agent         `gorm:"foreignKey:CreatedByAgentID" json:"created_by_agent,omitempty"`
	AssignedAgent    *Agent         `gorm:"foreignKey:AssignedAgentID" json:"assigned_agent,omitempty"`
	Comments         []TaskComment  `gorm:"foreignKey:TaskID" json:"comments,omitempty"`
	Events           []TaskEvent    `gorm:"foreignKey:TaskID" json:"events,omitempty"`
	Subtasks         []Subtask      `gorm:"foreignKey:TaskID" json:"subtasks,omitempty"`
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

// SubtaskStatus represents the status of a subtask
type SubtaskStatus string

const (
	SubtaskStatusPending   SubtaskStatus = "pending"
	SubtaskStatusCompleted SubtaskStatus = "completed"
)

// Subtask represents a subtask within a task
type Subtask struct {
	Base
	Title     string         `gorm:"not null" json:"title"`
	Status    SubtaskStatus  `gorm:"not null;default:'pending'" json:"status"`
	TaskID    uuid.UUID      `gorm:"type:uuid;not null" json:"task_id"`
	Task      *Task          `gorm:"foreignKey:TaskID" json:"task,omitempty"`
	Position  int            `gorm:"not null;default:0" json:"position"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// OrganisationStatus represents the status of an organisation
type OrganisationStatus string

const (
	OrganisationStatusActive    OrganisationStatus = "active"
	OrganisationStatusSuspended OrganisationStatus = "suspended"
	OrganisationStatusArchived  OrganisationStatus = "archived"
)

// Organisation represents a multi-tenant grouping
type Organisation struct {
	Base
	Name            string               `gorm:"not null" json:"name"`
	Slug            string               `gorm:"uniqueIndex;not null" json:"slug"`
	Description     string               `json:"description"`
	Status          OrganisationStatus   `gorm:"not null;default:'active'" json:"status"`
	CreatedByUserID uuid.UUID            `gorm:"type:uuid;not null" json:"created_by_user_id"`
	CreatedBy       *User                `gorm:"foreignKey:CreatedByUserID" json:"created_by,omitempty"`
	Members         []OrganisationMember `gorm:"foreignKey:OrganisationID" json:"members,omitempty"`
	Projects        []Project            `gorm:"foreignKey:OrganisationID" json:"projects,omitempty"`
	Agents          []Agent              `gorm:"foreignKey:OrganisationID" json:"agents,omitempty"`
}

// OrganisationMemberRole represents the role of a member in an organisation
type OrganisationMemberRole string

const (
	OrganisationMemberRoleOwner  OrganisationMemberRole = "owner"
	OrganisationMemberRoleAdmin  OrganisationMemberRole = "admin"
	OrganisationMemberRoleMember OrganisationMemberRole = "member"
)

// OrganisationMember represents a user's membership in an organisation
type OrganisationMember struct {
	Base
	OrganisationID uuid.UUID              `gorm:"type:uuid;not null" json:"organisation_id"`
	Organisation   *Organisation          `gorm:"foreignKey:OrganisationID" json:"organisation,omitempty"`
	UserID         uuid.UUID              `gorm:"type:uuid;not null" json:"user_id"`
	User           *User                  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Role           OrganisationMemberRole `gorm:"not null;default:'member'" json:"role"`
	JoinedAt       time.Time              `json:"joined_at"`
}

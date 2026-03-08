package services

import (
	"time"

	"github.com/formatho/agent-todo/db"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ActivityService struct {
	db *gorm.DB
}

func NewActivityService() *ActivityService {
	return &ActivityService{
		db: db.GetDB(),
	}
}

// ActivityEvent represents a formatted activity event for the feed
type ActivityEvent struct {
	ID            string    `json:"id"`
	TaskID        string    `json:"task_id"`
	TaskTitle     string    `json:"task_title"`
	EventType     string    `json:"event_type"`
	Description   string    `json:"description"`
	ActorName     string    `json:"actor_name"`
	ActorType     string    `json:"actor_type"` // "user" or "agent"
	PreviousState string    `json:"previous_state,omitempty"`
	NewState      string    `json:"new_state,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// GetRecentActivity fetches recent activity events across all tasks
func (s *ActivityService) GetRecentActivity(limit int) ([]ActivityEvent, error) {
	var events []ActivityEvent

	query := `
		SELECT 
			te.id,
			te.task_id,
			t.title as task_title,
			te.event_type,
			te.previous_state,
			te.new_state,
			te.changed_by,
			te.created_at,
			COALESCE(u.email, a.name) as actor_name,
			CASE 
				WHEN u.id IS NOT NULL THEN 'user'
				WHEN a.id IS NOT NULL THEN 'agent'
				ELSE 'system'
			END as actor_type
		FROM task_events te
		JOIN tasks t ON te.task_id = t.id
		LEFT JOIN users u ON te.changed_by = u.id::text
		LEFT JOIN agents a ON te.changed_by = a.id::text
		ORDER BY te.created_at DESC
		LIMIT ?
	`

	rows, err := s.db.Raw(query, limit).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event ActivityEvent
		var changedBy string

		err := rows.Scan(
			&event.ID,
			&event.TaskID,
			&event.TaskTitle,
			&event.EventType,
			&event.PreviousState,
			&event.NewState,
			&changedBy,
			&event.CreatedAt,
			&event.ActorName,
			&event.ActorType,
		)
		if err != nil {
			continue
		}

		// Generate human-readable description
		event.Description = generateDescription(event.EventType, event.TaskTitle, event.ActorName, event.PreviousState, event.NewState)

		events = append(events, event)
	}

	return events, nil
}

func generateDescription(eventType, taskTitle, actorName, previousState, newState string) string {
	switch eventType {
	case "created":
		return actorName + " created task \"" + taskTitle + "\""
	case "status_changed":
		return actorName + " changed status from " + previousState + " to " + newState + " for \"" + taskTitle + "\""
	case "assigned":
		return actorName + " assigned \"" + taskTitle + "\""
	case "unassigned":
		return actorName + " unassigned \"" + taskTitle + "\""
	case "comment_added":
		return actorName + " added a comment to \"" + taskTitle + "\""
	default:
		return actorName + " " + eventType + " on \"" + taskTitle + "\""
	}
}

// CreateEvent creates a new task event
func (s *ActivityService) CreateEvent(taskID uuid.UUID, eventType, previousState, newState, changedBy string) error {
	event := struct {
		ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
		TaskID        uuid.UUID `gorm:"type:uuid;not null"`
		EventType     string    `gorm:"not null"`
		PreviousState string
		NewState      string
		ChangedBy     string
		CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP"`
		UpdatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	}{
		TaskID:        taskID,
		EventType:     eventType,
		PreviousState: previousState,
		NewState:      newState,
		ChangedBy:     changedBy,
	}

	return s.db.Table("task_events").Create(&event).Error
}

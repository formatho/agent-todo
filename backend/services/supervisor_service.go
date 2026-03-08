package services

import (
	"time"

	"github.com/formatho/agent-todo/db"
	"github.com/formatho/agent-todo/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SupervisorService struct {
	db *gorm.DB
}

func NewSupervisorService() *SupervisorService {
	return &SupervisorService{
		db: db.GetDB(),
	}
}

// AgentWithTasks represents an agent with their active tasks
type AgentWithTasks struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Role        models.AgentRole `json:"role"`
	Enabled     bool           `json:"enabled"`
	CreatedAt   time.Time      `json:"created_at"`
	ActiveTasks []TaskSummary  `json:"active_tasks"`
}

// TaskSummary represents a brief task summary
type TaskSummary struct {
	ID       uuid.UUID         `json:"id"`
	Title    string            `json:"title"`
	Status   models.TaskStatus `json:"status"`
	Priority models.TaskPriority `json:"priority"`
}

// GetAgentsWithActiveTasks fetches all agents with their in-progress tasks
func (s *SupervisorService) GetAgentsWithActiveTasks() ([]AgentWithTasks, error) {
	var agents []models.Agent

	// Fetch all agents
	if err := s.db.Find(&agents).Error; err != nil {
		return nil, err
	}

	// For each agent, fetch their active tasks
	result := make([]AgentWithTasks, len(agents))
	for i, agent := range agents {
		var tasks []models.Task

		// Fetch in_progress tasks for this agent
		if err := s.db.Where("assigned_agent_id = ? AND status = ?", agent.ID, "in_progress").Find(&tasks).Error; err != nil {
			return nil, err
		}

		// Convert tasks to summaries
		taskSummaries := make([]TaskSummary, len(tasks))
		for j, task := range tasks {
			taskSummaries[j] = TaskSummary{
				ID:       task.ID,
				Title:    task.Title,
				Status:   task.Status,
				Priority: task.Priority,
			}
		}

		result[i] = AgentWithTasks{
			ID:          agent.ID,
			Name:        agent.Name,
			Description: agent.Description,
			Role:        agent.Role,
			Enabled:     agent.Enabled,
			CreatedAt:   agent.CreatedAt,
			ActiveTasks: taskSummaries,
		}
	}

	return result, nil
}

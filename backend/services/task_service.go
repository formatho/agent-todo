package services

import (
	"time"

	"github.com/formatho/agent-todo/db"
	"github.com/formatho/agent-todo/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskService struct {
	db *gorm.DB
}

func NewTaskService() *TaskService {
	return &TaskService{
		db: db.GetDB(),
	}
}

// TaskFilter represents filters for listing tasks
type TaskFilter struct {
	Status          models.TaskStatus   `json:"status"`
	AssignedAgentID string              `json:"agent_id"`
	Priority        models.TaskPriority `json:"priority"`
	CreatedByUserID string              `json:"created_by_user_id"`
	ProjectID       string              `json:"project_id"`
	SearchTerm      string              `json:"search_term"`
	DueDateFrom     *time.Time          `json:"due_date_from"`
	DueDateTo       *time.Time          `json:"due_date_to"`
}

// Create creates a new task
// createdByUserID: UUID of the user creating the task (for human-created tasks)
// createdByAgentID: UUID of the agent creating the task (for agent-created tasks)
// Exactly one of createdByUserID or createdByAgentID must be provided
func (s *TaskService) Create(title, description string, priority models.TaskPriority, dueDate *time.Time, projectID string, createdByUserID *string, assignedAgentID *string) (*models.Task, error) {
	task := &models.Task{
		Title:       title,
		Description: description,
		Status:      models.TaskStatusPending,
		Priority:    priority,
		DueDate:     dueDate,
	}

	// Set creator - exactly one must be provided
	if createdByUserID != nil && *createdByUserID != "" {
		parsedID := uuid.MustParse(*createdByUserID)
		task.CreatedByUserID = &parsedID
	} else {
		task.CreatedByUserID = nil
	}

	// Note: createdByAgentID will be set by a separate method or by the handler
	// We'll add a CreateByAgent method for clarity

	if projectID != "" {
		parsedProjectID := uuid.MustParse(projectID)
		task.ProjectID = &parsedProjectID
	}

	if assignedAgentID != nil {
		parsedID := uuid.MustParse(*assignedAgentID)
		task.AssignedAgentID = &parsedID
	}

	if err := s.db.Create(task).Error; err != nil {
		return nil, err
	}

	// Create creation event
	s.createEvent(task.ID, models.TaskEventCreated, "", string(task.Status), "system")

	return task, nil
}

// CreateByAgent creates a new task on behalf of an agent
func (s *TaskService) CreateByAgent(title, description string, priority models.TaskPriority, dueDate *time.Time, projectID, createdByAgentID string, assignedAgentID *string) (*models.Task, error) {
	task := &models.Task{
		Title:       title,
		Description: description,
		Status:      models.TaskStatusPending,
		Priority:    priority,
		DueDate:     dueDate,
	}

	// Set agent as creator
	if createdByAgentID != "" {
		parsedID := uuid.MustParse(createdByAgentID)
		task.CreatedByAgentID = &parsedID
	}

	if projectID != "" {
		parsedProjectID := uuid.MustParse(projectID)
		task.ProjectID = &parsedProjectID
	}

	if assignedAgentID != nil {
		parsedID := uuid.MustParse(*assignedAgentID)
		task.AssignedAgentID = &parsedID
	}

	if err := s.db.Create(task).Error; err != nil {
		return nil, err
	}

	// Create creation event
	s.createEvent(task.ID, models.TaskEventCreated, "", string(task.Status), "system")

	return task, nil
}

// GetByID retrieves a task with relationships
func (s *TaskService) GetByID(id string) (*models.Task, error) {
	var task models.Task
	err := s.db.Preload("Project").Preload("CreatedBy").Preload("CreatedByAgent").Preload("AssignedAgent").Preload("Comments").Preload("Events").
		Where("id = ?", id).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// List retrieves tasks with filters
func (s *TaskService) List(filter TaskFilter) ([]models.Task, error) {
	var tasks []models.Task
	query := s.db.Preload("Project").Preload("CreatedBy").Preload("CreatedByAgent").Preload("AssignedAgent")

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if filter.AssignedAgentID != "" {
		query = query.Where("assigned_agent_id = ?", filter.AssignedAgentID)
	}

	if filter.Priority != "" {
		query = query.Where("priority = ?", filter.Priority)
	}

	if filter.CreatedByUserID != "" {
		query = query.Where("created_by_user_id = ?", filter.CreatedByUserID)
	}

	if filter.ProjectID != "" {
		query = query.Where("project_id = ?", filter.ProjectID)
	}

	if filter.SearchTerm != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?",
			"%"+filter.SearchTerm+"%", "%"+filter.SearchTerm+"%")
	}

	if filter.DueDateFrom != nil {
		query = query.Where("due_date >= ?", *filter.DueDateFrom)
	}

	if filter.DueDateTo != nil {
		query = query.Where("due_date <= ?", *filter.DueDateTo)
	}

	err := query.Order("created_at DESC").Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// Update updates a task
func (s *TaskService) Update(id string, title, description *string, priority *models.TaskPriority, dueDate **time.Time, assignedAgentID *string) (*models.Task, error) {
	var task models.Task
	if err := s.db.Where("id = ?", id).First(&task).Error; err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})

	if title != nil {
		updates["title"] = *title
	}
	if description != nil {
		updates["description"] = *description
	}
	if priority != nil {
		updates["priority"] = *priority
	}
	if dueDate != nil {
		if *dueDate == nil {
			updates["due_date"] = nil
		} else {
			updates["due_date"] = **dueDate
		}
	}
	if assignedAgentID != nil {
		if *assignedAgentID == "" {
			updates["assigned_agent_id"] = nil
		} else {
			parsedID := uuid.MustParse(*assignedAgentID)
			updates["assigned_agent_id"] = parsedID
		}
	}

	if err := s.db.Model(&task).Updates(updates).Error; err != nil {
		return nil, err
	}

	// Create update event
	s.createEvent(task.ID, models.TaskEventUpdated, "", "", "system")

	return s.GetByID(id)
}

// UpdateStatus updates the status of a task
func (s *TaskService) UpdateStatus(id string, status models.TaskStatus, changedBy string) (*models.Task, error) {
	var task models.Task
	if err := s.db.Where("id = ?", id).First(&task).Error; err != nil {
		return nil, err
	}

	oldStatus := task.Status
	task.Status = status

	if err := s.db.Save(&task).Error; err != nil {
		return nil, err
	}

	// Create status change event
	s.createEvent(task.ID, models.TaskEventStatusChanged, string(oldStatus), string(status), changedBy)

	return s.GetByID(id)
}

// AssignAgent assigns an agent to a task
func (s *TaskService) AssignAgent(taskID, agentID string) (*models.Task, error) {
	var task models.Task
	if err := s.db.Where("id = ?", taskID).First(&task).Error; err != nil {
		return nil, err
	}

	var agent models.Agent
	if err := s.db.Where("id = ?", agentID).First(&agent).Error; err != nil {
		return nil, err
	}

	parsedAgentID := uuid.MustParse(agentID)
	task.AssignedAgentID = &parsedAgentID

	if err := s.db.Save(&task).Error; err != nil {
		return nil, err
	}

	// Create assignment event
	s.createEvent(task.ID, models.TaskEventAssigned, "", agent.Name, "system")

	return s.GetByID(taskID)
}

// UnassignAgent removes agent assignment from a task
func (s *TaskService) UnassignAgent(taskID string) (*models.Task, error) {
	var task models.Task
	if err := s.db.Where("id = ?", taskID).First(&task).Error; err != nil {
		return nil, err
	}

	oldAgentName := ""
	if task.AssignedAgent != nil {
		oldAgentName = task.AssignedAgent.Name
	}

	task.AssignedAgentID = nil

	if err := s.db.Save(&task).Error; err != nil {
		return nil, err
	}

	// Create unassignment event
	s.createEvent(task.ID, models.TaskEventUnassigned, oldAgentName, "", "system")

	return s.GetByID(taskID)
}

// Delete deletes a task
func (s *TaskService) Delete(id string) error {
	return s.db.Where("id = ?", id).Delete(&models.Task{}).Error
}

// AddComment adds a comment to a task
func (s *TaskService) AddComment(taskID, content string, authorID uuid.UUID, authorType, authorName string) (*models.TaskComment, error) {
	comment := &models.TaskComment{
		TaskID:     uuid.MustParse(taskID),
		Content:    content,
		AuthorID:   authorID,
		AuthorType: authorType,
		AuthorName: authorName,
	}

	if err := s.db.Create(comment).Error; err != nil {
		return nil, err
	}

	return comment, nil
}

// GetComments retrieves comments for a task
func (s *TaskService) GetComments(taskID string) ([]models.TaskComment, error) {
	var comments []models.TaskComment
	err := s.db.Where("task_id = ?", taskID).Order("created_at ASC").Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// createEvent creates a task event (helper method)
func (s *TaskService) createEvent(taskID uuid.UUID, eventType models.TaskEventType, previousState, newState, changedBy string) error {
	event := &models.TaskEvent{
		TaskID:        taskID,
		EventType:     eventType,
		PreviousState: previousState,
		NewState:      newState,
		ChangedBy:     changedBy,
	}
	return s.db.Create(event).Error
}

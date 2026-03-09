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
	OrganisationID  string              `json:"organisation_id"` // Optional: filter by organisation
}

// Create creates a new task
// createdByUserID: UUID of the user creating the task (for human-created tasks)
// createdByAgentID: UUID of the agent creating the task (for agent-created tasks)
// Exactly one of createdByUserID or createdByAgentID must be provided
// organisationID: Optional UUID of the organisation the task belongs to
func (s *TaskService) Create(title, description string, priority models.TaskPriority, dueDate *time.Time, commitURL, projectID string, createdByUserID *string, assignedAgentID *string, organisationID *string) (*models.Task, error) {
	task := &models.Task{
		Title:       title,
		Description: description,
		Status:      models.TaskStatusPending,
		Priority:    priority,
		DueDate:     dueDate,
		CommitURL:   commitURL,
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

	// Set organisation if provided
	if organisationID != nil && *organisationID != "" {
		parsedOrgID := uuid.MustParse(*organisationID)
		task.OrganisationID = &parsedOrgID
	}

	if err := s.db.Create(task).Error; err != nil {
		return nil, err
	}

	// Create creation event
	s.createEvent(task.ID, models.TaskEventCreated, "", string(task.Status), "system")

	return task, nil
}

// CreateByAgent creates a new task on behalf of an agent
// createdByAgentName is used for activity feed attribution
// organisationID: Optional UUID of the organisation the task belongs to
func (s *TaskService) CreateByAgent(title, description string, priority models.TaskPriority, dueDate *time.Time, commitURL, projectID, createdByAgentID, createdByAgentName string, assignedAgentID *string, organisationID *string) (*models.Task, error) {
	task := &models.Task{
		Title:       title,
		Description: description,
		Status:      models.TaskStatusPending,
		Priority:    priority,
		DueDate:     dueDate,
		CommitURL:   commitURL,
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

	// Set organisation if provided
	if organisationID != nil && *organisationID != "" {
		parsedOrgID := uuid.MustParse(*organisationID)
		task.OrganisationID = &parsedOrgID
	}

	if err := s.db.Create(task).Error; err != nil {
		return nil, err
	}

	// Create creation event with agent name for proper attribution
	actorName := createdByAgentName
	if actorName == "" {
		actorName = "system"
	}
	s.createEvent(task.ID, models.TaskEventCreated, "", string(task.Status), actorName)

	return task, nil
}

// GetByID retrieves a task with relationships
func (s *TaskService) GetByID(id string) (*models.Task, error) {
	var task models.Task
	err := s.db.Preload("Project").Preload("CreatedBy").Preload("CreatedByAgent").Preload("AssignedAgent").Preload("Comments").Preload("Events").Preload("Subtasks").
		Where("id = ?", id).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// GetByIDAndOrganisation retrieves a task by ID, ensuring it belongs to the specified organisation
func (s *TaskService) GetByIDAndOrganisation(id, organisationID string) (*models.Task, error) {
	var task models.Task
	err := s.db.Preload("Project").Preload("CreatedBy").Preload("CreatedByAgent").Preload("AssignedAgent").Preload("Comments").Preload("Events").Preload("Subtasks").
		Where("id = ? AND organisation_id = ?", id, organisationID).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// List retrieves tasks with filters
func (s *TaskService) List(filter TaskFilter) ([]models.Task, error) {
	var tasks []models.Task
	query := s.db.Preload("Project").Preload("CreatedBy").Preload("CreatedByAgent").Preload("AssignedAgent").Preload("Subtasks")

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

	// Filter by organisation if provided
	if filter.OrganisationID != "" {
		query = query.Where("organisation_id = ?", filter.OrganisationID)
	}

	err := query.Order("created_at DESC").Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// Update updates a task
func (s *TaskService) Update(id string, title, description *string, priority *models.TaskPriority, dueDate **time.Time, commitURL *string, assignedAgentID *string) (*models.Task, error) {
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
	if commitURL != nil {
		updates["commit_url"] = *commitURL
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

// UpdateByOrganisation updates a task, verifying it belongs to the organisation
func (s *TaskService) UpdateByOrganisation(id, organisationID string, title, description *string, priority *models.TaskPriority, dueDate **time.Time, commitURL *string, assignedAgentID *string) (*models.Task, error) {
	var task models.Task
	if err := s.db.Where("id = ? AND organisation_id = ?", id, organisationID).First(&task).Error; err != nil {
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
	if commitURL != nil {
		updates["commit_url"] = *commitURL
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

	return s.GetByIDAndOrganisation(id, organisationID)
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

// AssignAgentByOrganisation assigns an agent to a task, verifying both belong to the organisation
func (s *TaskService) AssignAgentByOrganisation(taskID, agentID, organisationID string) (*models.Task, error) {
	var task models.Task
	if err := s.db.Where("id = ? AND organisation_id = ?", taskID, organisationID).First(&task).Error; err != nil {
		return nil, err
	}

	var agent models.Agent
	if err := s.db.Where("id = ? AND organisation_id = ?", agentID, organisationID).First(&agent).Error; err != nil {
		return nil, gorm.ErrRecordNotFound // Agent not found in this organisation
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

// UnassignAgentByOrganisation removes agent assignment, verifying task belongs to organisation
func (s *TaskService) UnassignAgentByOrganisation(taskID, organisationID string) (*models.Task, error) {
	var task models.Task
	if err := s.db.Where("id = ? AND organisation_id = ?", taskID, organisationID).First(&task).Error; err != nil {
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

// DeleteByOrganisation deletes a task, verifying it belongs to the organisation
func (s *TaskService) DeleteByOrganisation(id, organisationID string) error {
	result := s.db.Where("id = ? AND organisation_id = ?", id, organisationID).Delete(&models.Task{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
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

// GetUpcomingDueTasks retrieves tasks with due dates within the specified duration
// that are not completed yet
func (s *TaskService) GetUpcomingDueTasks(within time.Duration) ([]models.Task, error) {
	var tasks []models.Task
	now := time.Now()
	dueBefore := now.Add(within)

	err := s.db.Preload("Project").Preload("CreatedByAgent").Preload("AssignedAgent").Preload("Subtasks").
		Where("due_date IS NOT NULL").
		Where("due_date > ?", now).
		Where("due_date <= ?", dueBefore).
		Where("status != ?", models.TaskStatusCompleted).
		Where("status != ?", models.TaskStatusFailed).
		Order("due_date ASC").
		Find(&tasks).Error

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetUpcomingDueTasksByOrganisation retrieves upcoming due tasks for a specific organisation
func (s *TaskService) GetUpcomingDueTasksByOrganisation(within time.Duration, organisationID string) ([]models.Task, error) {
	var tasks []models.Task
	now := time.Now()
	dueBefore := now.Add(within)

	err := s.db.Preload("Project").Preload("CreatedByAgent").Preload("AssignedAgent").Preload("Subtasks").
		Where("organisation_id = ?", organisationID).
		Where("due_date IS NOT NULL").
		Where("due_date > ?", now).
		Where("due_date <= ?", dueBefore).
		Where("status != ?", models.TaskStatusCompleted).
		Where("status != ?", models.TaskStatusFailed).
		Order("due_date ASC").
		Find(&tasks).Error

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetOverdueTasks retrieves tasks that are past their due date and not completed
func (s *TaskService) GetOverdueTasks() ([]models.Task, error) {
	var tasks []models.Task
	now := time.Now()

	err := s.db.Preload("Project").Preload("CreatedByAgent").Preload("AssignedAgent").Preload("Subtasks").
		Where("due_date IS NOT NULL").
		Where("due_date < ?", now).
		Where("status != ?", models.TaskStatusCompleted).
		Where("status != ?", models.TaskStatusFailed).
		Order("due_date ASC").
		Find(&tasks).Error

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetOverdueTasksByOrganisation retrieves overdue tasks for a specific organisation
func (s *TaskService) GetOverdueTasksByOrganisation(organisationID string) ([]models.Task, error) {
	var tasks []models.Task
	now := time.Now()

	err := s.db.Preload("Project").Preload("CreatedByAgent").Preload("AssignedAgent").Preload("Subtasks").
		Where("organisation_id = ?", organisationID).
		Where("due_date IS NOT NULL").
		Where("due_date < ?", now).
		Where("status != ?", models.TaskStatusCompleted).
		Where("status != ?", models.TaskStatusFailed).
		Order("due_date ASC").
		Find(&tasks).Error

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

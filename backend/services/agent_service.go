package services

import (
	"errors"

	"github.com/formatho/agent-todo/db"
	"github.com/formatho/agent-todo/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AgentService struct {
	db *gorm.DB
}

func NewAgentService() *AgentService {
	return &AgentService{
		db: db.GetDB(),
	}
}

// Create creates a new agent
func (s *AgentService) Create(name, description string) (*models.Agent, error) {
	// Generate API key
	apiKey := "sk_agent_" + uuid.New().String()

	agent := &models.Agent{
		Name:        name,
		APIKey:      apiKey,
		Description: description,
		Role:        models.AgentRoleRegular,
		Enabled:     true,
	}

	if err := s.db.Create(agent).Error; err != nil {
		return nil, err
	}

	return agent, nil
}

// CreateWithRole creates a new agent with a specific role
func (s *AgentService) CreateWithRole(name, description string, role models.AgentRole) (*models.Agent, error) {
	// Generate API key
	apiKey := "sk_agent_" + uuid.New().String()

	agent := &models.Agent{
		Name:        name,
		APIKey:      apiKey,
		Description: description,
		Role:        role,
		Enabled:     true,
	}

	if err := s.db.Create(agent).Error; err != nil {
		return nil, err
	}

	return agent, nil
}

// CreateWithOrganisation creates a new agent with organisation context
func (s *AgentService) CreateWithOrganisation(name, description string, role models.AgentRole, organisationID string) (*models.Agent, error) {
	// Generate API key
	apiKey := "sk_agent_" + uuid.New().String()

	orgUUID := uuid.MustParse(organisationID)
	agent := &models.Agent{
		Name:           name,
		APIKey:         apiKey,
		Description:    description,
		Role:           role,
		Enabled:        true,
		OrganisationID: &orgUUID,
	}

	if err := s.db.Create(agent).Error; err != nil {
		return nil, err
	}

	return agent, nil
}

// GetByID retrieves an agent by ID
func (s *AgentService) GetByID(id string) (*models.Agent, error) {
	var agent models.Agent
	err := s.db.Where("id = ?", id).First(&agent).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

// GetByIDAndOrganisation retrieves an agent by ID, ensuring it belongs to the organisation
func (s *AgentService) GetByIDAndOrganisation(id, organisationID string) (*models.Agent, error) {
	var agent models.Agent
	err := s.db.Where("id = ? AND organisation_id = ?", id, organisationID).First(&agent).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

// GetByAPIKey retrieves an agent by API key
func (s *AgentService) GetByAPIKey(apiKey string) (*models.Agent, error) {
	var agent models.Agent
	err := s.db.Where("api_key = ?", apiKey).First(&agent).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid API key")
		}
		return nil, err
	}
	return &agent, nil
}

// List retrieves all agents
func (s *AgentService) List() ([]models.Agent, error) {
	var agents []models.Agent
	err := s.db.Find(&agents).Error
	if err != nil {
		return nil, err
	}
	return agents, nil
}

// ListByOrganisation retrieves all agents for a specific organisation
func (s *AgentService) ListByOrganisation(organisationID string) ([]models.Agent, error) {
	var agents []models.Agent
	err := s.db.Where("organisation_id = ?", organisationID).Find(&agents).Error
	if err != nil {
		return nil, err
	}
	return agents, nil
}

// Delete deletes an agent
func (s *AgentService) Delete(id string) error {
	return s.db.Where("id = ?", id).Delete(&models.Agent{}).Error
}

// DeleteByOrganisation deletes an agent, verifying it belongs to the organisation
func (s *AgentService) DeleteByOrganisation(id, organisationID string) error {
	result := s.db.Where("id = ? AND organisation_id = ?", id, organisationID).Delete(&models.Agent{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("agent not found in organisation")
	}
	return nil
}

// Update updates an agent's details
func (s *AgentService) Update(id, name, description string, role models.AgentRole, enabled *bool) (*models.Agent, error) {
	var agent models.Agent
	if err := s.db.Where("id = ?", id).First(&agent).Error; err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})

	if name != "" {
		updates["name"] = name
	}
	if description != "" {
		updates["description"] = description
	}
	if role != "" {
		updates["role"] = role
	}
	if enabled != nil {
		updates["enabled"] = *enabled
	}

	if err := s.db.Model(&agent).Updates(updates).Error; err != nil {
		return nil, err
	}

	// Reload to get updated data
	if err := s.db.Where("id = ?", id).First(&agent).Error; err != nil {
		return nil, err
	}

	return &agent, nil
}

// UpdateByOrganisation updates an agent, verifying it belongs to the organisation
func (s *AgentService) UpdateByOrganisation(id, organisationID, name, description string, role models.AgentRole, enabled *bool) (*models.Agent, error) {
	var agent models.Agent
	if err := s.db.Where("id = ? AND organisation_id = ?", id, organisationID).First(&agent).Error; err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})

	if name != "" {
		updates["name"] = name
	}
	if description != "" {
		updates["description"] = description
	}
	if role != "" {
		updates["role"] = role
	}
	if enabled != nil {
		updates["enabled"] = *enabled
	}

	if err := s.db.Model(&agent).Updates(updates).Error; err != nil {
		return nil, err
	}

	// Reload to get updated data
	if err := s.db.Where("id = ?", id).First(&agent).Error; err != nil {
		return nil, err
	}

	return &agent, nil
}

// AgentStatistics represents statistics for an agent
type AgentStatistics struct {
	TotalTasks     int64 `json:"total_tasks"`
	PendingTasks   int64 `json:"pending_tasks"`
	InProgress     int64 `json:"in_progress"`
	CompletedTasks int64 `json:"completed_tasks"`
	FailedTasks    int64 `json:"failed_tasks"`
	BlockedTasks   int64 `json:"blocked_tasks"`
}

// GetStatistics retrieves task statistics for an agent
func (s *AgentService) GetStatistics(agentID string) (*AgentStatistics, error) {
	stats := &AgentStatistics{}

	// Get total tasks
	if err := s.db.Model(&models.Task{}).Where("assigned_agent_id = ?", agentID).Count(&stats.TotalTasks).Error; err != nil {
		return nil, err
	}

	// Get tasks by status
	if err := s.db.Model(&models.Task{}).Where("assigned_agent_id = ? AND status = ?", agentID, "pending").Count(&stats.PendingTasks).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&models.Task{}).Where("assigned_agent_id = ? AND status = ?", agentID, "in_progress").Count(&stats.InProgress).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&models.Task{}).Where("assigned_agent_id = ? AND status = ?", agentID, "completed").Count(&stats.CompletedTasks).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&models.Task{}).Where("assigned_agent_id = ? AND status = ?", agentID, "failed").Count(&stats.FailedTasks).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&models.Task{}).Where("assigned_agent_id = ? AND status = ?", agentID, "blocked").Count(&stats.BlockedTasks).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

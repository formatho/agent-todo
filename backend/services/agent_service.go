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

// GetByID retrieves an agent by ID
func (s *AgentService) GetByID(id string) (*models.Agent, error) {
	var agent models.Agent
	err := s.db.Where("id = ?", id).First(&agent).Error
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

// Delete deletes an agent
func (s *AgentService) Delete(id string) error {
	return s.db.Where("id = ?", id).Delete(&models.Agent{}).Error
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

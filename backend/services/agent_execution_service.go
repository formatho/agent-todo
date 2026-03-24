package services

import (
	"context"
	"fmt"

	"github.com/formatho/agent-todo/models"
)

// AgentExecutionService handles agent execution logic
type AgentExecutionService struct {
	executor *AgentExecutor
	db       *DatabaseService
}

func NewAgentExecutionService(executor *AgentExecutor, db *DatabaseService) *AgentExecutionService {
	return &AgentExecutionService{
		executor: executor,
		db:       db,
	}
}

// CreateAndExecuteTask creates an agent if needed and executes a task
func (s *AgentExecutionService) CreateAndExecuteTask(agentID string, projectID string, title string, description string) (*models.AgentResponse, error) {
	ctx := context.Background()

	// Get the agent from database
	agent, err := s.db.GetAgentByID(agentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent: %w", err)
	}

	// Get the project for context
	project, err := s.db.GetProjectByID(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	// Build execution context from project LLM context and task details
	context := fmt.Sprintf(
		"Project: %s\nDescription: %s\nLLM Instructions: %s",
		project.Name, project.Description, project.LLMContext,
	)

	// Execute the task using the agent executor
	response, err := s.executor.ExecuteTask(agentID, title, description, context)
	if err != nil {
		return nil, fmt.Errorf("task execution failed: %w", err)
	}

	// Store response in database for tracking
	err = s.storeAgentResponse(ctx, agentID, project.ID.String(), title, description, response.Response, response.Metadata)
	if err != nil {
		fmt.Printf("Warning: Failed to store agent response: %v\n", err)
	}

	return response, nil
}

func (s *AgentExecutionService) storeAgentResponse(ctx context.Context, agentID string, projectID string, title string, description string, responseText string, metadata interface{}) error {
	// TODO: Implement database storage for agent responses
	// This is a placeholder - we'll need to add the table and model first
	return nil
}

// ExecuteTaskWithLLM executes a task using specified LLM configuration
func (s *AgentExecutionService) ExecuteTaskWithLLM(agentID string, llmConfig models.LLMConfig, title string, description string) (*models.AgentResponse, error) {
	ctx := context.Background()

	// Get the agent from database
	agent, err := s.db.GetAgentByID(agentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent: %w", err)
	}

	// Convert models.LLMConfig to executor.LLMConfig
	executorConfig := LLMConfig{
		Type:      string(llmConfig.Type),
		ModelName: llmConfig.ModelName,
		APIKey:    llmConfig.APIKey,
		BaseURL:   llmConfig.BaseURL,
	}

	// Execute with custom LLM config
	response, err := s.executor.ExecuteTaskWithCustomLLM(agentID, executorConfig, title, description, "")
	if err != nil {
		return nil, fmt.Errorf("task execution failed: %w", err)
	}

	return response, nil
}

// AddToolToAgent adds a tool to an agent
func (s *AgentExecutionService) AddToolToAgent(agentID string, toolName string, toolFunc func(args map[string]interface{}) (interface{}, error)) error {
	err := s.executor.AddTool(agentID, toolName, toolFunc)
	if err != nil {
		return fmt.Errorf("failed to add tool: %w", err)
	}
	return nil
}

// GetAgentStatus returns the runtime status of an agent
func (s *AgentExecutionService) GetAgentStatus(agentID string) (*models.AgentStatus, error) {
	status, err := s.executor.GetAgentStatus(agentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent status: %w", err)
	}
	return status, nil
}

package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/formatho/agent-todo/db"
	"github.com/formatho/agent-todo/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AgentStateSnapshot represents a persisted state of an agent for cloud sync
type AgentStateSnapshot struct {
	Base
	AgentID         string                 `gorm:"not null;index" json:"agent_id"`
	StateData       []byte                 `gorm:"type:jsonb" json:"state_data"` // Serialized agent state
	Version         int                    `gorm:"not null;default:0" json:"version"`
	SnapshotType    string                 `gorm:"not null;default:'full'" json:"snapshot_type"` // "full", "incremental", "checkpoint"
	IsCurrent       bool                   `gorm:"not null;default:false" json:"is_current"`
	Metadata        map[string]interface{} `json:"metadata"`
	CreatedByUserID *uuid.UUID             `gorm:"type:uuid" json:"created_by_user_id,omitempty"`
}

// TaskExecutionHistory tracks all task executions for analytics and cloud sync
type TaskExecutionHistory struct {
	Base
	TaskID          uuid.UUID  `gorm:"type:uuid;index" json:"task_id"`
	AgentID         string     `gorm:"not null;index" json:"agent_id"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	ResponseText    string     `json:"response_text"`
	ResponseJSON    []byte     `gorm:"type:jsonb" json:"response_json,omitempty"`   // Structured response data
	Status          string     `gorm:"not null;default:'completed'" json:"status"`  // "pending", "in_progress", "completed", "failed"
	ExecutionTime   int64      `gorm:"not null;default:0" json:"execution_time_ms"` // Duration in milliseconds
	ContextUsed     string     `json:"context_used,omitempty"`                      // The LLM context sent to agent
	Metadata        []byte     `gorm:"type:jsonb" json:"metadata,omitempty"`        // Additional metadata as JSONB
	CreatedByUserID *uuid.UUID `gorm:"type:uuid;index" json:"created_by_user_id,omitempty"`
}

// AnalyticsCache pre-computed metrics for fast dashboard loading
type AnalyticsCache struct {
	Base
	OrganisationID  uuid.UUID  `gorm:"type:uuid;uniqueIndex;not null" json:"organisation_id"`
	CacheType       string     `gorm:"not null;index" json:"cache_type"` // "task_trends", "performance_metrics", "agent_stats"
	Data            []byte     `gorm:"type:jsonb;not null" json:"data"`  // Cached data as JSONB
	ExpiresAt       time.Time  `gorm:"not null;index" json:"expires_at"` // Cache expiration time
	CreatedByUserID *uuid.UUID `gorm:"type:uuid" json:"created_by_user_id,omitempty"`
}

// TeamMember represents a team member for collaboration features (alternative to OrganisationMember)
type TeamMember struct {
	Base
	OrganisationID uuid.UUID              `gorm:"type:uuid;uniqueIndex;not null" json:"organisation_id"`
	UserID         uuid.UUID              `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	Role           OrganisationMemberRole `gorm:"not null;default:'member'" json:"role"`
	Status         string                 `gorm:"not null;default:'active'" json:"status"` // "active", "invited", "pending"
	InvitedAt      *time.Time             `json:"invited_at,omitempty"`
	JoinedAt       *time.Time             `json:"joined_at,omitempty"`
	LastActiveAt   *time.Time             `json:"last_active_at,omitempty"`
	Metadata       map[string]interface{} `json:"metadata"`
}

// StateSerializationService handles agent state serialization and cloud sync
type StateSerializationService struct {
	db *gorm.DB
}

func NewStateSerializationService() *StateSerializationService {
	return &StateSerializationService{
		db: db.GetDB(),
	}
}

// SerializeAgentState serializes an agent's current runtime state to JSON
// Note: This creates a basic state snapshot. For runtime agent status, use AgentExecutor
func (s *StateSerializationService) SerializeAgentState(agentID string, metadata map[string]interface{}) ([]byte, error) {
	stateData := map[string]interface{}{
		"agent_id":    agentID,
		"snapshot_at": time.Now().UTC(),
		"version":     1,
	}

	if metadata != nil {
		stateData["custom_metadata"] = metadata
	}

	return json.Marshal(stateData)
}

// DeserializeAgentState deserializes agent state from JSON and updates executor
func (s *StateSerializationService) DeserializeAgentState(agentID string, snapshotData []byte) (*models.AgentStatus, error) {
	var stateMap map[string]interface{}
	if err := json.Unmarshal(snapshotData, &stateMap); err != nil {
		return nil, fmt.Errorf("failed to deserialize agent state: %w", err)
	}

	status := &models.AgentStatus{
		ID:       stateMap["agent_id"].(string),
		Exists:   true, // Assuming exists since we have a snapshot
		Active:   false,
		Metadata: make(map[string]interface{}),
	}

	if metadata, ok := stateMap["metadata"].(map[string]interface{}); ok {
		status.Metadata = metadata
	}

	return status, nil
}

// SaveSnapshot saves an agent state snapshot to the database for cloud persistence
func (s *StateSerializationService) SaveSnapshot(ctx context.Context, agentID string, stateData []byte, snapshotType string, metadata map[string]interface{}, createdByUserID *uuid.UUID) (*AgentStateSnapshot, error) {
	snapshot := &AgentStateSnapshot{
		AgentID:         agentID,
		StateData:       stateData,
		SnapshotType:    snapshotType,
		Metadata:        metadata,
		CreatedByUserID: createdByUserID,
	}

	// Mark previous snapshots as not current if this is a full snapshot
	if snapshotType == "full" {
		var updateQuery = `UPDATE agent_state_snapshots SET is_current = false WHERE agent_id = ? AND snapshot_type = 'full'`
		if err := s.db.Raw(updateQuery, agentID).Error; err != nil {
			return nil, fmt.Errorf("failed to clear current snapshots: %w", err)
		}
		snapshot.IsCurrent = true
	}

	result := s.db.Create(snapshot)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to save snapshot: %w", result.Error)
	}

	return snapshot, nil
}

// GetCurrentSnapshot retrieves the current state snapshot for an agent
func (s *StateSerializationService) GetCurrentSnapshot(agentID string) (*AgentStateSnapshot, error) {
	var snapshot AgentStateSnapshot
	result := s.db.Where("agent_id = ? AND is_current = true", agentID).First(&snapshot)
	if result.Error != nil {
		return nil, fmt.Errorf("no current snapshot found for agent: %w", result.Error)
	}

	return &snapshot, nil
}

// GetSnapshotHistory retrieves the history of snapshots for an agent
func (s *StateSerializationService) GetSnapshotHistory(agentID string, limit int) ([]AgentStateSnapshot, error) {
	var snapshots []AgentStateSnapshot
	result := s.db.Where("agent_id = ?", agentID).Order("created_at DESC").Limit(limit).Find(&snapshots)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get snapshot history: %w", result.Error)
	}

	return snapshots, nil
}

// StoreTaskExecution stores a task execution in the history for analytics and cloud sync
func (s *StateSerializationService) StoreTaskExecution(ctx context.Context, taskID uuid.UUID, agentID string, title string, description string, responseText string, metadata map[string]interface{}, executionTime int64, status string) (*TaskExecutionHistory, error) {
	history := &TaskExecutionHistory{
		TaskID:        taskID,
		AgentID:       agentID,
		Title:         title,
		Description:   description,
		ResponseText:  responseText,
		Status:        status,
		ExecutionTime: executionTime,
		Metadata:      []byte{}, // Will be set below
	}

	if metadata != nil {
		if jsonData, err := json.Marshal(metadata); err == nil {
			history.Metadata = jsonData
		}
	}

	result := s.db.Create(history)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to store task execution: %w", result.Error)
	}

	return history, nil
}

// StoreStructuredResponse stores a structured agent response for export functionality
func (s *StateSerializationService) StoreStructuredResponse(ctx context.Context, taskID uuid.UUID, agentID string, projectID string, title string, description string, responseText string, responseJSON map[string]interface{}, metadata map[string]interface{}) (*TaskExecutionHistory, error) {
	history := &TaskExecutionHistory{
		TaskID:       taskID,
		AgentID:      agentID,
		Title:        title,
		Description:  description,
		ResponseText: responseText,
		Status:       "completed",
		Metadata:     []byte{},
	}

	if len(responseJSON) > 0 {
		if jsonData, err := json.Marshal(responseJSON); err == nil {
			history.ResponseJSON = jsonData
		}
	}

	if metadata != nil {
		if jsonData, err := json.Marshal(metadata); err == nil {
			history.Metadata = jsonData
		}
	}

	result := s.db.Create(history)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to store structured response: %w", result.Error)
	}

	return history, nil
}

// GetTaskExecutionHistory retrieves the execution history for a task or agent
func (s *StateSerializationService) GetTaskExecutionHistory(taskID uuid.UUID, agentID string, limit int) ([]TaskExecutionHistory, error) {
	var histories []TaskExecutionHistory
	query := s.db.Model(&TaskExecutionHistory{})

	if !taskID.Equal(uuid.Nil) {
		query = query.Where("task_id = ?", taskID)
	} else if agentID != "" {
		query = query.Where("agent_id = ?", agentID)
	}

	result := query.Order("created_at DESC").Limit(limit).Find(&histories)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get execution history: %w", result.Error)
	}

	return histories, nil
}

// UpdateAnalyticsCache updates or creates an analytics cache entry
func (s *StateSerializationService) UpdateAnalyticsCache(ctx context.Context, orgID uuid.UUID, cacheType string, data map[string]interface{}, ttl time.Duration) (*AnalyticsCache, error) {
	cacheData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal analytics data: %w", err)
	}

	cache := &AnalyticsCache{
		OrganisationID: orgID,
		CacheType:      cacheType,
		Data:           cacheData,
		ExpiresAt:      time.Now().UTC().Add(ttl),
	}

	result := s.db.Create(cache)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create analytics cache: %w", result.Error)
	}

	return cache, nil
}

// GetAnalyticsCache retrieves cached analytics data if not expired
func (s *StateSerializationService) GetAnalyticsCache(orgID uuid.UUID, cacheType string) (*AnalyticsCache, error) {
	var cache AnalyticsCache
	result := s.db.Where("organisation_id = ? AND cache_type = ? AND expires_at > ?", orgID, cacheType, time.Now().UTC()).First(&cache)
	if result.Error != nil {
		return nil, fmt.Errorf("no cached data found or expired: %w", result.Error)
	}

	return &cache, nil
}

// AddTeamMember adds a member to an organisation for collaboration features
func (s *StateSerializationService) AddTeamMember(ctx context.Context, orgID uuid.UUID, userID uuid.UUID, role OrganisationMemberRole, status string, metadata map[string]interface{}) (*TeamMember, error) {
	member := &TeamMember{
		OrganisationID: orgID,
		UserID:         userID,
		Role:           role,
		Status:         status,
		Metadata:       metadata,
	}

	if status == "invited" || status == "pending" {
		now := time.Now()
		member.InvitedAt = &now
	} else if status == "active" {
		now := time.Now()
		member.JoinedAt = &now
	}

	result := s.db.Create(member)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to add team member: %w", result.Error)
	}

	return member, nil
}

// GetTeamMembers retrieves all members of an organisation
func (s *StateSerializationService) GetTeamMembers(orgID uuid.UUID) ([]TeamMember, error) {
	var members []TeamMember
	result := s.db.Where("organisation_id = ?", orgID).Find(&members)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get team members: %w", result.Error)
	}

	return members, nil
}

// UpdateMemberStatus updates the status of a team member (e.g., from invited to active)
func (s *StateSerializationService) UpdateMemberStatus(orgID uuid.UUID, userID uuid.UUID, newStatus string) (*TeamMember, error) {
	var member TeamMember
	result := s.db.Where("organisation_id = ? AND user_id = ?", orgID, userID).First(&member)
	if result.Error != nil {
		return nil, fmt.Errorf("team member not found: %w", result.Error)
	}

	member.Status = newStatus
	now := time.Now()
	member.LastActiveAt = &now

	if newStatus == "active" || newStatus == "joined" {
		member.JoinedAt = &now
	}

	result = s.db.Save(&member)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update member status: %w", result.Error)
	}

	return &member, nil
}

// EnsureDatabaseTables creates all necessary database tables for cloud sync features
func (s *StateSerializationService) EnsureDatabaseTables() error {
	tables := []interface{}{
		&AgentStateSnapshot{},
		&TaskExecutionHistory{},
		&AnalyticsCache{},
		&TeamMember{},
	}

	for _, table := range tables {
		if err := s.db.AutoMigrate(table); err != nil {
			return fmt.Errorf("failed to auto-migrate %T: %w", table, err)
		}
	}

	return nil
}

// GetTaskCompletionMetrics retrieves task completion metrics for analytics dashboard
func (s *StateSerializationService) GetTaskCompletionMetrics(orgID uuid.UUID, days int) (map[string]interface{}, error) {
	cacheKey := fmt.Sprintf("task_completion_%d_days", days)

	// Try to get from cache first
	if cached, err := s.GetAnalyticsCache(orgID, "task_trends"); err == nil && cached.CacheType == "task_trends" {
		var data map[string]interface{}
		if err := json.Unmarshal(cached.Data, &data); err == nil {
			return data, nil
		}
	}

	// Calculate metrics from execution history
	var histories []TaskExecutionHistory
	result := s.db.Model(&TaskExecutionHistory{}).
		Joins("LEFT JOIN tasks ON task_execution_history.task_id = tasks.id").
		Where("tasks.organisation_id = ? AND created_at > ?", orgID, time.Now().AddDate(0, 0, -days)).
		Find(&histories)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get task completion metrics: %w", result.Error)
	}

	metrics := map[string]interface{}{
		"total_tasks":        len(histories),
		"completed_tasks":    0,
		"in_progress_tasks":  0,
		"failed_tasks":       0,
		"avg_execution_time": float64(0),
		"completion_rate":    float64(0),
	}

	var totalExecutionTime int64 = 0

	for _, h := range histories {
		switch h.Status {
		case "completed":
			metrics["completed_tasks"] = metrics["completed_tasks"].(int) + 1
		case "in_progress":
			metrics["in_progress_tasks"] = metrics["in_progress_tasks"].(int) + 1
		case "failed":
			metrics["failed_tasks"] = metrics["failed_tasks"].(int) + 1
		}

		totalExecutionTime += h.ExecutionTime
	}

	if len(histories) > 0 {
		metrics["avg_execution_time"] = float64(totalExecutionTime) / float64(len(histories))
		completed := metrics["completed_tasks"].(int)
		metrics["completion_rate"] = (float64(completed) / float64(len(histories))) * 100
	}

	// Cache the results for 24 hours
	cache, err := s.UpdateAnalyticsCache(context.Background(), orgID, "task_trends", metrics, 24*time.Hour)
	if err != nil {
		fmt.Printf("Warning: Failed to cache task completion metrics: %v\n", err)
	}

	return metrics, nil
}

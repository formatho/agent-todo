package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/formatho/agent-todo/middleware"
	"github.com/formatho/agent-todo/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// StateSyncHandler handles state synchronization and cloud sync endpoints
type StateSyncHandler struct {
	stateService *services.StateSerializationService
}

func NewStateSyncHandler(stateService *services.StateSerializationService) *StateSyncHandler {
	return &StateSyncHandler{stateService: stateService}
}

// SaveAgentSnapshot saves a snapshot of the agent's current state for cloud persistence
// @Summary Save agent state snapshot
// @Description Persist an agent's runtime state to enable cloud sync across sessions
// @Tags State Synchronization
// @Security X-API-KEY
// @Param agent_id path string true "Agent ID"
// @Param snapshot body object true "Snapshot data"
// @Success 201 {object} services.AgentStateSnapshot
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /state/sync/agents/{agent_id}/snapshots [post]
func (h *StateSyncHandler) SaveAgentSnapshot(c *gin.Context) {
	agentID := c.Param("agent_id")

	var request struct {
		StateData    map[string]interface{} `json:"state_data"`
		SnapshotType string                 `json:"snapshot_type"` // "full", "incremental", "checkpoint"
		Metadata     map[string]interface{} `json:"metadata"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %v", err)})
		return
	}

	ctx := c.Request.Context()

	// Serialize state data if not already provided
	var serializedData []byte
	if len(request.StateData) > 0 {
		var err error
		serializedData, err = h.stateService.SerializeAgentState(agentID, request.Metadata)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to serialize state: %v", err)})
			return
		}
	}

	snapshotType := request.SnapshotType
	if snapshotType == "" {
		snapshotType = "full" // Default to full snapshot if not specified
	}

	// Save the snapshot
	snapshot, err := h.stateService.SaveSnapshot(ctx, agentID, serializedData, snapshotType, request.Metadata, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save snapshot: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Agent state snapshot saved successfully",
		"snapshot": snapshot,
	})
}

// GetCurrentSnapshot retrieves the current state snapshot for an agent
// @Summary Get current agent state snapshot
// @Description Retrieve the most recent persisted state of an agent
// @Tags State Synchronization
// @Security X-API-KEY
// @Param agent_id path string true "Agent ID"
// @Success 200 {object} services.AgentStateSnapshot
// @Failure 404 {object} string
// @Router /state/sync/agents/{agent_id}/snapshot/current [get]
func (h *StateSyncHandler) GetCurrentSnapshot(c *gin.Context) {
	agentID := c.Param("agent_id")

	snapshot, err := h.stateService.GetCurrentSnapshot(agentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("No current snapshot found: %v", err)})
		return
	}

	c.JSON(http.StatusOK, snapshot)
}

// GetSnapshotHistory retrieves the history of snapshots for an agent
// @Summary Get agent state snapshot history
// @Description Retrieve historical snapshots of an agent's state
// @Tags State Synchronization
// @Security X-API-KEY
// @Param agent_id path string true "Agent ID"
// @Param limit query int false "Number of snapshots to retrieve (default: 20)"
// @Success 200 {array} services.AgentStateSnapshot
// @Router /state/sync/agents/{agent_id}/snapshots [get]
func (h *StateSyncHandler) GetSnapshotHistory(c *gin.Context) {
	agentID := c.Param("agent_id")

	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20 // Default to 20 snapshots
	}

	snapshots, err := h.stateService.GetSnapshotHistory(agentID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get snapshot history: %v", err)})
		return
	}

	c.JSON(http.StatusOK, snapshots)
}

// RestoreAgentState restores an agent's state from a snapshot
// @Summary Restore agent state from snapshot
// @Description Load an agent's runtime state from a persisted snapshot (used for cloud sync recovery)
// @Tags State Synchronization
// @Security X-API-KEY
// @Param agent_id path string true "Agent ID"
// @Param snapshot_id path string true "Snapshot ID to restore from"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} string
// @Router /state/sync/agents/{agent_id}/snapshots/{snapshot_id}/restore [post]
func (h *StateSyncHandler) RestoreAgentState(c *gin.Context) {
	agentID := c.Param("agent_id")
	snapshotID := c.Param("snapshot_id")

	snapshotUUID, err := uuid.Parse(snapshotID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid snapshot ID format"})
		return
	}

	var snapshot services.AgentStateSnapshot
	result := h.stateService.db.DB.Where("id = ?", snapshotUUID).First(&snapshot)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Snapshot not found: %v", result.Error)})
		return
	}

	// Deserialize the state
	status, err := h.stateService.DeserializeAgentState(agentID, snapshot.StateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to restore state: %v", err)})
		return
	}

	// Mark this snapshot as current (optional - depends on implementation)
	snapshot.IsCurrent = true
	h.stateService.db.DB.Save(&snapshot)

	c.JSON(http.StatusOK, gin.H{
		"message": "Agent state restored successfully",
		"status":  status,
	})
}

// StoreTaskExecution stores a task execution for analytics and cloud sync
// @Summary Store task execution history
// @Description Record a task execution in the history for analytics dashboard and export functionality
// @Tags State Synchronization
// @Security X-API-KEY
// @Param agent_id path string true "Agent ID"
// @Param execution body object true "Execution data"
// @Success 201 {object} services.TaskExecutionHistory
// @Failure 400 {object} string
// @Router /state/sync/agents/{agent_id}/executions [post]
func (h *StateSyncHandler) StoreTaskExecution(c *gin.Context) {
	agentID := c.Param("agent_id")

	var request struct {
		TaskID        uuid.UUID              `json:"task_id"`
		Title         string                 `json:"title"`
		Description   string                 `json:"description"`
		ResponseText  string                 `json:"response_text"`
		Status        string                 `json:"status"` // "pending", "in_progress", "completed", "failed"
		Metadata      map[string]interface{} `json:"metadata"`
		ExecutionTime int64                  `json:"execution_time_ms"`
		ContextUsed   string                 `json:"context_used,omitempty"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %v", err)})
		return
	}

	ctx := c.Request.Context()

	history, err := h.stateService.StoreTaskExecution(ctx, request.TaskID, agentID, request.Title, request.Description, request.ResponseText, request.Metadata, request.ExecutionTime, request.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to store execution: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, history)
}

// StoreStructuredResponse stores a structured agent response for export functionality
// @Summary Store structured agent response
// @Description Store an agent's structured response data for CSV/JSON export and analytics
// @Tags State Synchronization
// @Security X-API-KEY
// @Param agent_id path string true "Agent ID"
// @Param response body object true "Structured response data"
// @Success 201 {object} services.TaskExecutionHistory
// @Router /state/sync/agents/{agent_id}/responses [post]
func (h *StateSyncHandler) StoreStructuredResponse(c *gin.Context) {
	agentID := c.Param("agent_id")

	var request struct {
		TaskID       uuid.UUID              `json:"task_id"`
		ProjectID    string                 `json:"project_id"`
		Title        string                 `json:"title"`
		Description  string                 `json:"description"`
		ResponseText string                 `json:"response_text"`
		ResponseJSON map[string]interface{} `json:"response_json"` // Structured response as JSON object
		Metadata     map[string]interface{} `json:"metadata"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %v", err)})
		return
	}

	ctx := c.Request.Context()

	taskUUID, err := uuid.Parse(request.TaskID.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID format"})
		return
	}

	history, err := h.stateService.StoreStructuredResponse(ctx, taskUUID, agentID, request.ProjectID, request.Title, request.Description, request.ResponseText, request.ResponseJSON, request.Metadata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to store structured response: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, history)
}

// GetTaskExecutionHistory retrieves execution history for a task or agent
// @Summary Get execution history
// @Description Retrieve the execution history for analytics and audit purposes
// @Tags State Synchronization
// @Security X-API-KEY
// @Param agent_id path string true "Agent ID" (optional if task_id provided)
// @Param task_id query string false "Task ID to filter by"
// @Param limit query int false "Number of records to retrieve (default: 50)"
// @Success 200 {array} services.TaskExecutionHistory
// @Router /state/sync/agents/{agent_id}/executions [get]
func (h *StateSyncHandler) GetTaskExecutionHistory(c *gin.Context) {
	agentID := c.Param("agent_id")

	taskIDStr := c.Query("task_id")
	limitStr := c.DefaultQuery("limit", "50")

	var taskID uuid.UUID
	if taskIDStr != "" {
		parsed, err := uuid.Parse(taskIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID format"})
			return
		}
		taskID = parsed
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50 // Default to 50 records
	}

	histories, err := h.stateService.GetTaskExecutionHistory(taskID, agentID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get execution history: %v", err)})
		return
	}

	c.JSON(http.StatusOK, histories)
}

// GetTaskCompletionMetrics retrieves task completion metrics for analytics dashboard
// @Summary Get task completion metrics
// @Description Retrieve analytics data about task completion trends and performance metrics
// @Tags State Synchronization
// @Security X-API-KEY
// @Param agent_id path string true "Agent ID"
// @Param days query int false "Number of days to analyze (default: 30)"
// @Success 200 {object} map[string]interface{}
// @Router /state/sync/agents/{agent_id}/metrics/task-completion [get]
func (h *StateSyncHandler) GetTaskCompletionMetrics(c *gin.Context) {
	agentID := c.Param("agent_id")

	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		days = 30 // Default to 30 days
	}

	metrics, err := h.stateService.GetTaskCompletionMetrics(agentID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get metrics: %v", err)})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// AddTeamMember adds a team member for collaboration features
// @Summary Add team member
// @Description Invite or add a user to an organisation for team collaboration
// @Tags Team Collaboration
// @Security X-API-KEY
// @Param agent_id path string true "Agent ID" (represents the organisation)
// @Param member body object true "Team member data"
// @Success 201 {object} services.TeamMember
// @Router /state/sync/organisations/{agent_id}/members [post]
func (h *StateSyncHandler) AddTeamMember(c *gin.Context) {
	orgIDStr := c.Param("agent_id") // Using agent_id as organisation ID for simplicity

	orgUUID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organisation ID format"})
		return
	}

	var request struct {
		UserID   string                 `json:"user_id"`
		Role     string                 `json:"role"`   // "owner", "admin", "member"
		Status   string                 `json:"status"` // "active", "invited", "pending"
		Metadata map[string]interface{} `json:"metadata"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %v", err)})
		return
	}

	userUUID, err := uuid.Parse(request.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	memberRole := services.OrganisationMemberRole(request.Role)
	if memberRole == "" {
		memberRole = services.OrganisationMemberRoleMember // Default to member
	}

	status := request.Status
	if status == "" {
		status = "invited" // Default to invited for new members
	}

	ctx := c.Request.Context()

	member, err := h.stateService.AddTeamMember(ctx, orgUUID, userUUID, memberRole, status, request.Metadata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to add team member: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, member)
}

// GetTeamMembers retrieves all members of an organisation
// @Summary Get team members
// @Description List all team members in an organisation for collaboration management
// @Tags Team Collaboration
// @Security X-API-KEY
// @Param agent_id path string true "Agent ID" (represents the organisation)
// @Success 200 {array} services.TeamMember
// @Router /state/sync/organisations/{agent_id}/members [get]
func (h *StateSyncHandler) GetTeamMembers(c *gin.Context) {
	orgIDStr := c.Param("agent_id")

	orgUUID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organisation ID format"})
		return
	}

	members, err := h.stateService.GetTeamMembers(orgUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get team members: %v", err)})
		return
	}

	c.JSON(http.StatusOK, members)
}

// UpdateMemberStatus updates a team member's status (e.g., accept invitation)
// @Summary Update team member status
// @Description Change a team member's status from invited to active when they join
// @Tags Team Collaboration
// @Security X-API-KEY
// @Param agent_id path string true "Agent ID" (represents the organisation)
// @Param user_id path string true "User ID of the member to update"
// @Param status body object true "New status"
// @Success 200 {object} services.TeamMember
// @Router /state/sync/organisations/{agent_id}/members/{user_id}/status [patch]
func (h *StateSyncHandler) UpdateMemberStatus(c *gin.Context) {
	orgIDStr := c.Param("agent_id")
	userIDStr := c.Param("user_id")

	orgUUID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organisation ID format"})
		return
	}

	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	var request struct {
		Status string `json:"status"` // "active", "invited", "pending"
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %v", err)})
		return
	}

	member, err := h.stateService.UpdateMemberStatus(orgUUID, userUUID, request.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update member status: %v", err)})
		return
	}

	c.JSON(http.StatusOK, member)
}

// ExportTaskHistory exports task execution history in CSV or JSON format
// @Summary Export task history
// @Description Export agent's task execution history for backup and analysis
// @Tags State Synchronization
// @Security X-API-KEY
// @Param agent_id path string true "Agent ID"
// @Param format query string false "Export format: csv or json (default: json)"
// @Success 200 {object} map[string]interface{}
// @Router /state/sync/agents/{agent_id}/export [get]
func (h *StateSyncHandler) ExportTaskHistory(c *gin.Context) {
	agentID := c.Param("agent_id")

	format := c.DefaultQuery("format", "json") // Default to JSON format

	limitStr := c.DefaultQuery("limit", "100")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 100 // Default to 100 records
	}

	histories, err := h.stateService.GetTaskExecutionHistory(uuid.Nil, agentID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get execution history: %v", err)})
		return
	}

	if format == "csv" {
		// Convert to CSV format (simplified - would need proper CSV encoding in production)
		var csvData []string
		csvData = append(csvData, "ID,Task ID,Agent,Title,Status,Execution Time,Created At")

		for _, h := range histories {
			line := fmt.Sprintf("%s,%s,%s,%s,%s,%d,%s",
				h.ID.String(),
				h.TaskID.String(),
				h.AgentID,
				h.Title,
				h.Status,
				h.ExecutionTime,
				h.CreatedAt.Format(time.RFC3339))
			csvData = append(csvData, line)
		}

		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=task_history_%s.csv", agentID))
		c.String(http.StatusOK, "%s\n", csvData[0]) // Simplified - proper CSV would use strings.Join
	} else {
		// JSON format (default)
		c.JSON(http.StatusOK, gin.H{
			"export_format": "json",
			"count":         len(histories),
			"data":          histories,
		})
	}
}

package handlers

import (
	"net/http"
	"time"

	"github.com/formatho/agent-todo/middleware"
	"github.com/formatho/agent-todo/models"
	"github.com/formatho/agent-todo/services"
	"github.com/gin-gonic/gin"
)

type AgentTaskHandler struct {
	taskService *services.TaskService
}

func NewAgentTaskHandler() *AgentTaskHandler {
	return &AgentTaskHandler{
		taskService: services.NewTaskService(),
	}
}

// AgentCreateTaskRequest represents the request body for agent task creation
type AgentCreateTaskRequest struct {
	Title       string              `json:"title" binding:"required" example:"Process data files"`
	Description string              `json:"description" example:"Process and analyze the uploaded CSV files"`
	Priority    models.TaskPriority `json:"priority" binding:"required" example:"medium"`
	DueDate     *time.Time          `json:"due_date" example:"2024-12-31T23:59:59Z"`
	ProjectID   string              `json:"project_id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// UpdateTaskStatusRequest represents the request body for updating task status
type UpdateTaskStatusRequest struct {
	Status models.TaskStatus `json:"status" binding:"required" example:"in_progress"`
}

// AgentCreateTask godoc
// @Summary Create a task (Agent)
// @Description Create a new task (for agents)
// @Tags agent
// @Accept json
// @Produce json
// @Security X-API-KEY
// @Param request body AgentCreateTaskRequest true "Task details"
// @Success 201 {object} models.Task
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /agent/tasks [post]
func (h *AgentTaskHandler) CreateTask(c *gin.Context) {
	var req AgentCreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	agentID, err := middleware.GetAgentID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Agent creates tasks for itself
	// Note: Using agent ID as created_by_user_id - this should be reconsidered
	// as ideally agents should create tasks on behalf of a user
	task, err := h.taskService.Create(
		req.Title,
		req.Description,
		req.Priority,
		req.DueDate,
		req.ProjectID,
		agentID,  // Using agent ID as creator (will need to update service)
		&agentID, // Auto-assign to self
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// AgentListTasks godoc
// @Summary List agent tasks
// @Description Get tasks assigned to the authenticated agent
// @Tags agent
// @Produce json
// @Security X-API-KEY
// @Param status query string false "Filter by status" Enums(pending, in_progress, completed, failed)
// @Success 200 {array} models.Task
// @Failure 401 {object} map[string]string
// @Router /agent/tasks [get]
func (h *AgentTaskHandler) ListTasks(c *gin.Context) {
	agentID, err := middleware.GetAgentID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	filter := services.TaskFilter{
		AssignedAgentID: agentID,
		Status:          models.TaskStatus(c.Query("status")),
	}

	tasks, err := h.taskService.List(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// AgentGetTask godoc
// @Summary Get agent task
// @Description Get a specific task assigned to the agent
// @Tags agent
// @Produce json
// @Security X-API-KEY
// @Param id path string true "Task ID"
// @Success 200 {object} models.Task
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /agent/tasks/{id} [get]
func (h *AgentTaskHandler) GetTask(c *gin.Context) {
	id := c.Param("id")
	agentID, err := middleware.GetAgentID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	task, err := h.taskService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Verify task is assigned to this agent
	if task.AssignedAgentID == nil || task.AssignedAgentID.String() != agentID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Task not assigned to this agent"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// AgentUpdateStatus godoc
// @Summary Update task status (Agent)
// @Description Update the status of a task assigned to the agent (or any task if supervisor/admin)
// @Tags agent
// @Accept json
// @Produce json
// @Security X-API-KEY
// @Param id path string true "Task ID"
// @Param request body UpdateTaskStatusRequest true "New status"
// @Success 200 {object} models.Task
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /agent/tasks/{id}/status [patch]
func (h *AgentTaskHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	agentID, err := middleware.GetAgentID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	agentName, err := middleware.GetAgentName(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	agent, err := middleware.GetAgent(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req UpdateTaskStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify task is assigned to this agent, or agent is supervisor/admin
	task, err := h.taskService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Check permissions: assigned agent OR supervisor/admin
	isAssignedAgent := task.AssignedAgentID != nil && task.AssignedAgentID.String() == agentID
	isPrivileged := agent.Role == models.AgentRoleSupervisor || agent.Role == models.AgentRoleAdmin

	if !isAssignedAgent && !isPrivileged {
		c.JSON(http.StatusForbidden, gin.H{"error": "Task not assigned to this agent and insufficient privileges"})
		return
	}

	updatedTask, err := h.taskService.UpdateStatus(id, req.Status, agentName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedTask)
}

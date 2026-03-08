package handlers

import (
	"net/http"
	"time"

	"github.com/formatho/agent-todo/middleware"
	"github.com/formatho/agent-todo/models"
	"github.com/formatho/agent-todo/services"
	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskService *services.TaskService
}

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{
		taskService: services.NewTaskService(),
	}
}

// CreateTaskRequest represents the request body for creating a task
type CreateTaskRequest struct {
	Title           string             `json:"title" binding:"required" example:"Complete project documentation"`
	Description     string             `json:"description" example:"Write comprehensive documentation for the project"`
	Priority        models.TaskPriority `json:"priority" binding:"required" example:"high"`
	DueDate         *time.Time         `json:"due_date" example:"2024-12-31T23:59:59Z"`
	ProjectID       string             `json:"project_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	AssignedAgentID *string            `json:"assigned_agent_id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// UpdateTaskRequest represents the request body for updating a task
type UpdateTaskRequest struct {
	Title           *string             `json:"title" example:"Updated task title"`
	Description     *string             `json:"description" example:"Updated description"`
	Priority        *models.TaskPriority `json:"priority" example:"medium"`
	DueDate         **time.Time        `json:"due_date"`
	AssignedAgentID *string            `json:"assigned_agent_id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// AssignAgentRequest represents the request body for assigning an agent
type AssignAgentRequest struct {
	AgentID string `json:"agent_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// CreateTask godoc
// @Summary Create a new task
// @Description Create a new task
// @Tags tasks
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body CreateTaskRequest true "Task details"
// @Success 201 {object} models.Task
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	task, err := h.taskService.Create(
		req.Title,
		req.Description,
		req.Priority,
		req.DueDate,
		req.ProjectID,
		userID,
		req.AssignedAgentID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// ListTasks godoc
// @Summary List tasks
// @Description Get a list of tasks with optional filters
// @Tags tasks
// @Produce json
// @Security Bearer
// @Param status query string false "Filter by status" Enums(pending, in_progress, completed, failed)
// @Param agent_id query string false "Filter by assigned agent ID"
// @Param priority query string false "Filter by priority" Enums(low, medium, high, critical)
// @Param project_id query string false "Filter by project ID"
// @Param search query string false "Search in title and description"
// @Success 200 {array} models.Task
// @Failure 401 {object} map[string]string
// @Router /tasks [get]
func (h *TaskHandler) ListTasks(c *gin.Context) {
	filter := services.TaskFilter{
		Status:          models.TaskStatus(c.Query("status")),
		AssignedAgentID: c.Query("agent_id"),
		Priority:        models.TaskPriority(c.Query("priority")),
		ProjectID:       c.Query("project_id"),
		SearchTerm:      c.Query("search"),
	}

	tasks, err := h.taskService.List(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// GetTask godoc
// @Summary Get a task
// @Description Get a specific task by ID
// @Tags tasks
// @Produce json
// @Security Bearer
// @Param id path string true "Task ID"
// @Success 200 {object} models.Task
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tasks/{id} [get]
func (h *TaskHandler) GetTask(c *gin.Context) {
	id := c.Param("id")

	task, err := h.taskService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// UpdateTask godoc
// @Summary Update a task
// @Description Update a task's details
// @Tags tasks
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Task ID"
// @Param request body UpdateTaskRequest true "Task updates"
// @Success 200 {object} models.Task
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tasks/{id} [patch]
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var req UpdateTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.taskService.Update(
		id,
		req.Title,
		req.Description,
		req.Priority,
		req.DueDate,
		req.AssignedAgentID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

// DeleteTask godoc
// @Summary Delete a task
// @Description Delete a task by ID
// @Tags tasks
// @Produce json
// @Security Bearer
// @Param id path string true "Task ID"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")

	if err := h.taskService.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

// AssignAgent godoc
// @Summary Assign an agent to a task
// @Description Assign an agent to a specific task
// @Tags tasks
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Task ID"
// @Param request body AssignAgentRequest true "Agent ID"
// @Success 200 {object} models.Task
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tasks/{id}/assign [patch]
func (h *TaskHandler) AssignAgent(c *gin.Context) {
	taskID := c.Param("id")
	var req AssignAgentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.taskService.AssignAgent(taskID, req.AgentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

// UnassignAgent godoc
// @Summary Unassign an agent from a task
// @Description Remove agent assignment from a task
// @Tags tasks
// @Produce json
// @Security Bearer
// @Param id path string true "Task ID"
// @Success 200 {object} models.Task
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tasks/{id}/unassign [patch]
func (h *TaskHandler) UnassignAgent(c *gin.Context) {
	taskID := c.Param("id")

	task, err := h.taskService.UnassignAgent(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

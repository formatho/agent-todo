package handlers

import (
	"net/http"

	"github.com/formatho/agent-todo/middleware"
	"github.com/formatho/agent-todo/models"
	"github.com/formatho/agent-todo/services"
	"github.com/gin-gonic/gin"
)

type ToolsHandler struct {
	taskService *services.TaskService
}

func NewToolsHandler() *ToolsHandler {
	return &ToolsHandler{
		taskService: services.NewTaskService(),
	}
}

// ToolCreateTaskRequest represents a tool request for creating tasks
type ToolCreateTaskRequest struct {
	Title       string              `json:"title" binding:"required" example:"Analyze data"`
	Description string              `json:"description" example:"Process the CSV files"`
	Priority    models.TaskPriority `json:"priority" example:"high"`
	ProjectID   string              `json:"project_id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// ToolUpdateTaskRequest represents a tool request for updating tasks
type ToolUpdateTaskRequest struct {
	TaskID  string            `json:"task_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Status  models.TaskStatus `json:"status" binding:"required" example:"in_progress"`
	Comment string            `json:"comment" example:"Started working on this task"`
}

// ToolListRequest represents a tool request for listing tasks
type ToolListRequest struct {
	Status models.TaskStatus `json:"status" example:"pending"`
	Limit  int               `json:"limit" example:"10"`
}

// ToolResponse represents a standardized tool response
type ToolResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ToolCreateTask godoc
// @Summary Create task (Tool endpoint)
// @Description OpenClaw-compatible endpoint for creating tasks
// @Tags tools
// @Accept json
// @Produce json
// @Security X-API-KEY
// @Param request body ToolCreateTaskRequest true "Task details"
// @Success 200 {object} ToolResponse
// @Failure 400 {object} ToolResponse
// @Failure 401 {object} ToolResponse
// @Router /tools/tasks/create [post]
func (h *ToolsHandler) CreateTask(c *gin.Context) {
	agentID, err := middleware.GetAgentID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ToolResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	agentName, err := middleware.GetAgentName(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ToolResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	var req ToolCreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ToolResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Agent creates task using CreateByAgent
	task, err := h.taskService.CreateByAgent(
		req.Title,
		req.Description,
		req.Priority,
		nil,
		req.ProjectID,
		agentID,
		agentName,
		&agentID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ToolResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ToolResponse{
		Success: true,
		Data:    task,
	})
}

// ToolUpdateTask godoc
// @Summary Update task (Tool endpoint)
// @Description OpenClaw-compatible endpoint for updating task status
// @Tags tools
// @Accept json
// @Produce json
// @Security X-API-KEY
// @Param request body ToolUpdateTaskRequest true "Update details"
// @Success 200 {object} ToolResponse
// @Failure 400 {object} ToolResponse
// @Failure 401 {object} ToolResponse
// @Router /tools/tasks/update [post]
func (h *ToolsHandler) UpdateTask(c *gin.Context) {
	agentID, err := middleware.GetAgentID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ToolResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	agentName, err := middleware.GetAgentName(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ToolResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	var req ToolUpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ToolResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Verify task ownership
	task, err := h.taskService.GetByID(req.TaskID)
	if err != nil {
		c.JSON(http.StatusNotFound, ToolResponse{
			Success: false,
			Error:   "Task not found",
		})
		return
	}

	if task.AssignedAgentID == nil || task.AssignedAgentID.String() != agentID {
		c.JSON(http.StatusForbidden, ToolResponse{
			Success: false,
			Error:   "Task not assigned to this agent",
		})
		return
	}

	// Update status
	updatedTask, err := h.taskService.UpdateStatus(req.TaskID, req.Status, agentName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ToolResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Add comment if provided
	if req.Comment != "" {
		h.taskService.AddComment(req.TaskID, req.Comment, *task.AssignedAgentID, "agent", agentName)
	}

	c.JSON(http.StatusOK, ToolResponse{
		Success: true,
		Data:    updatedTask,
	})
}

// ToolListTasks godoc
// @Summary List tasks (Tool endpoint)
// @Description OpenClaw-compatible endpoint for listing agent tasks
// @Tags tools
// @Accept json
// @Produce json
// @Security X-API-KEY
// @Param request body ToolListRequest false "Filter options"
// @Success 200 {object} ToolResponse
// @Failure 401 {object} ToolResponse
// @Router /tools/tasks/list [post]
func (h *ToolsHandler) ListTasks(c *gin.Context) {
	agentID, err := middleware.GetAgentID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ToolResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	var req ToolListRequest
	c.ShouldBindJSON(&req)

	if req.Limit == 0 {
		req.Limit = 50
	}

	filter := services.TaskFilter{
		AssignedAgentID: agentID,
		Status:          req.Status,
	}

	tasks, err := h.taskService.List(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ToolResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Apply limit
	if len(tasks) > req.Limit {
		tasks = tasks[:req.Limit]
	}

	c.JSON(http.StatusOK, ToolResponse{
		Success: true,
		Data:    tasks,
	})
}

// ToolGetStatus godoc
// @Summary Get task status (Tool endpoint)
// @Description OpenClaw-compatible endpoint for getting task status
// @Tags tools
// @Produce json
// @Security X-API-KEY
// @Param id path string true "Task ID"
// @Success 200 {object} ToolResponse
// @Failure 401 {object} ToolResponse
// @Failure 404 {object} ToolResponse
// @Router /tools/tasks/status/{id} [get]
func (h *ToolsHandler) GetStatus(c *gin.Context) {
	agentID, err := middleware.GetAgentID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ToolResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	taskID := c.Param("id")

	task, err := h.taskService.GetByID(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, ToolResponse{
			Success: false,
			Error:   "Task not found",
		})
		return
	}

	// Verify ownership
	if task.AssignedAgentID == nil || task.AssignedAgentID.String() != agentID {
		c.JSON(http.StatusForbidden, ToolResponse{
			Success: false,
			Error:   "Task not assigned to this agent",
		})
		return
	}

	c.JSON(http.StatusOK, ToolResponse{
		Success: true,
		Data: gin.H{
			"task_id": task.ID,
			"status":  task.Status,
			"title":   task.Title,
		},
	})
}

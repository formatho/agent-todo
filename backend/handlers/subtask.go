package handlers

import (
	"net/http"

	"github.com/formatho/agent-todo/middleware"
	"github.com/formatho/agent-todo/models"
	"github.com/formatho/agent-todo/services"
	"github.com/gin-gonic/gin"
)

type SubtaskHandler struct {
	subtaskService *services.SubtaskService
	taskService    *services.TaskService
}

func NewSubtaskHandler() *SubtaskHandler {
	return &SubtaskHandler{
		subtaskService: services.NewSubtaskService(),
		taskService:    services.NewTaskService(),
	}
}

// CreateSubtaskRequest represents the request body for creating a subtask
type CreateSubtaskRequest struct {
	Title    string `json:"title" binding:"required" example:"Implement feature X"`
	Position *int   `json:"position" example:"0"`
}

// UpdateSubtaskRequest represents the request body for updating a subtask
type UpdateSubtaskRequest struct {
	Title    string               `json:"title" example:"Updated title"`
	Status   models.SubtaskStatus `json:"status" example:"completed"`
	Position *int                 `json:"position" example:"1"`
}

// ReorderSubtasksRequest represents the request body for reordering subtasks
type ReorderSubtasksRequest struct {
	SubtaskIDs []string `json:"subtask_ids" binding:"required" example:"uuid1,uuid2,uuid3"`
}

// CreateSubtask godoc
// @Summary Create a subtask
// @Description Create a new subtask for a task
// @Tags subtasks
// @Accept json
// @Produce json
// @Security Bearer
// @Param task_id path string true "Task ID"
// @Param request body CreateSubtaskRequest true "Subtask details"
// @Success 201 {object} models.Subtask
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tasks/{task_id}/subtasks [post]
func (h *SubtaskHandler) CreateSubtask(c *gin.Context) {
	taskID := c.Param("id")

	var req CreateSubtaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subtask, err := h.subtaskService.Create(taskID, req.Title, req.Position)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, subtask)
}

// GetSubtask godoc
// @Summary Get a subtask
// @Description Get a specific subtask by ID
// @Tags subtasks
// @Produce json
// @Security Bearer
// @Param id path string true "Subtask ID"
// @Success 200 {object} models.Subtask
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /subtasks/{id} [get]
func (h *SubtaskHandler) GetSubtask(c *gin.Context) {
	id := c.Param("id")

	subtask, err := h.subtaskService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subtask)
}

// UpdateSubtask godoc
// @Summary Update a subtask
// @Description Update a subtask's title, status, or position
// @Tags subtasks
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Subtask ID"
// @Param request body UpdateSubtaskRequest true "Update details"
// @Success 200 {object} models.Subtask
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /subtasks/{id} [patch]
func (h *SubtaskHandler) UpdateSubtask(c *gin.Context) {
	id := c.Param("id")

	var req UpdateSubtaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subtask, err := h.subtaskService.Update(id, req.Title, req.Status, req.Position)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subtask)
}

// DeleteSubtask godoc
// @Summary Delete a subtask
// @Description Delete a subtask
// @Tags subtasks
// @Produce json
// @Security Bearer
// @Param id path string true "Subtask ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /subtasks/{id} [delete]
func (h *SubtaskHandler) DeleteSubtask(c *gin.Context) {
	id := c.Param("id")

	if err := h.subtaskService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListSubtasks godoc
// @Summary List subtasks for a task
// @Description Get all subtasks for a specific task
// @Tags subtasks
// @Produce json
// @Security Bearer
// @Param task_id path string true "Task ID"
// @Success 200 {array} models.Subtask
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tasks/{task_id}/subtasks [get]
func (h *SubtaskHandler) ListSubtasks(c *gin.Context) {
	taskID := c.Param("id")

	subtasks, err := h.subtaskService.ListByTask(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subtasks)
}

// ReorderSubtasks godoc
// @Summary Reorder subtasks
// @Description Update the order of subtasks for a task
// @Tags subtasks
// @Accept json
// @Produce json
// @Security Bearer
// @Param task_id path string true "Task ID"
// @Param request body ReorderSubtasksRequest true "New order of subtask IDs"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tasks/{task_id}/subtasks/reorder [post]
func (h *SubtaskHandler) ReorderSubtasks(c *gin.Context) {
	taskID := c.Param("id")

	var req ReorderSubtasksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.subtaskService.Reorder(taskID, req.SubtaskIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subtasks reordered successfully"})
}

// Agent handlers (for agent API endpoints)

// AgentCreateSubtask godoc
// @Summary Create a subtask (Agent)
// @Description Create a new subtask for a task (for agents)
// @Tags agent
// @Accept json
// @Produce json
// @Security X-API-KEY
// @Param task_id path string true "Task ID"
// @Param request body CreateSubtaskRequest true "Subtask details"
// @Success 201 {object} models.Subtask
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /agent/tasks/{task_id}/subtasks [post]
func (h *SubtaskHandler) AgentCreateSubtask(c *gin.Context) {
	taskID := c.Param("id")

	// Verify task is assigned to this agent
	agentID, err := middleware.GetAgentID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	task, err := h.taskService.GetByID(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if task.AssignedAgentID == nil || task.AssignedAgentID.String() != agentID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Task not assigned to this agent"})
		return
	}

	var req CreateSubtaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subtask, err := h.subtaskService.Create(taskID, req.Title, req.Position)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, subtask)
}

// AgentUpdateSubtask godoc
// @Summary Update a subtask (Agent)
// @Description Update a subtask's title, status, or position (for agents)
// @Tags agent
// @Accept json
// @Produce json
// @Security X-API-KEY
// @Param id path string true "Subtask ID"
// @Param request body UpdateSubtaskRequest true "Update details"
// @Success 200 {object} models.Subtask
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /agent/subtasks/{id} [patch]
func (h *SubtaskHandler) AgentUpdateSubtask(c *gin.Context) {
	id := c.Param("id")

	// Get subtask to verify task assignment
	subtask, err := h.subtaskService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Verify task is assigned to this agent
	agentID, err := middleware.GetAgentID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	task, err := h.taskService.GetByID(subtask.TaskID.String())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if task.AssignedAgentID == nil || task.AssignedAgentID.String() != agentID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Task not assigned to this agent"})
		return
	}

	var req UpdateSubtaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedSubtask, err := h.subtaskService.Update(id, req.Title, req.Status, req.Position)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedSubtask)
}

// AgentDeleteSubtask godoc
// @Summary Delete a subtask (Agent)
// @Description Delete a subtask (for agents)
// @Tags agent
// @Produce json
// @Security X-API-KEY
// @Param id path string true "Subtask ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /agent/subtasks/{id} [delete]
func (h *SubtaskHandler) AgentDeleteSubtask(c *gin.Context) {
	id := c.Param("id")

	// Get subtask to verify task assignment
	subtask, err := h.subtaskService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Verify task is assigned to this agent
	agentID, err := middleware.GetAgentID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	task, err := h.taskService.GetByID(subtask.TaskID.String())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if task.AssignedAgentID == nil || task.AssignedAgentID.String() != agentID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Task not assigned to this agent"})
		return
	}

	if err := h.subtaskService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// AgentListSubtasks godoc
// @Summary List subtasks for a task (Agent)
// @Description Get all subtasks for a specific task (for agents)
// @Tags agent
// @Produce json
// @Security X-API-KEY
// @Param task_id path string true "Task ID"
// @Success 200 {array} models.Subtask
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /agent/tasks/{task_id}/subtasks [get]
func (h *SubtaskHandler) AgentListSubtasks(c *gin.Context) {
	taskID := c.Param("id")

	// Verify task is assigned to this agent
	agentID, err := middleware.GetAgentID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	task, err := h.taskService.GetByID(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if task.AssignedAgentID == nil || task.AssignedAgentID.String() != agentID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Task not assigned to this agent"})
		return
	}

	subtasks, err := h.subtaskService.ListByTask(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subtasks)
}

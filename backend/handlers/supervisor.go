package handlers

import (
	"net/http"

	"github.com/formatho/agent-todo/models"
	"github.com/formatho/agent-todo/services"
	"github.com/gin-gonic/gin"
)

type SupervisorHandler struct {
	agentService *services.AgentService
	taskService  *services.TaskService
}

func NewSupervisorHandler() *SupervisorHandler {
	return &SupervisorHandler{
		agentService: services.NewAgentService(),
		taskService:  services.NewTaskService(),
	}
}

// CreateAgentRequest represents the request body for supervisor creating an agent
type SupervisorCreateAgentRequest struct {
	Name        string           `json:"name" binding:"required" example:"My Agent"`
	Description string           `json:"description" example:"An AI assistant agent"`
	Role        models.AgentRole `json:"role" binding:"required" example:"regular"`
}

// SupervisorCreateAgent godoc
// @Summary Create a new agent (supervisor only)
// @Description Create a new AI agent with specific role (supervisor/admin only)
// @Tags supervisor
// @Accept json
// @Produce json
// @Security X-API-KEY
// @Param request body SupervisorCreateAgentRequest true "Agent details"
// @Success 201 {object} models.Agent
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /supervisor/agents [post]
func (h *SupervisorHandler) CreateAgent(c *gin.Context) {
	var req SupervisorCreateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate role
	if req.Role != models.AgentRoleRegular && req.Role != models.AgentRoleSupervisor && req.Role != models.AgentRoleAdmin {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role. Must be regular, supervisor, or admin"})
		return
	}

	agent, err := h.agentService.CreateWithRole(req.Name, req.Description, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, agent)
}

// SupervisorListAgents godoc
// @Summary List all agents (supervisor only)
// @Description Get a list of all agents (supervisor/admin only)
// @Tags supervisor
// @Produce json
// @Security X-API-KEY
// @Success 200 {array} models.Agent
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /supervisor/agents [get]
func (h *SupervisorHandler) ListAgents(c *gin.Context) {
	agents, err := h.agentService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, agents)
}

// SupervisorUpdateAgentRequest represents the request body for updating an agent
type SupervisorUpdateAgentRequest struct {
	Name        string           `json:"name" example:"Updated Agent Name"`
	Description string           `json:"description" example:"Updated description"`
	Role        models.AgentRole `json:"role" example:"supervisor"`
	Enabled     *bool            `json:"enabled" example:"true"`
}

// SupervisorUpdateAgent godoc
// @Summary Update an agent (supervisor only)
// @Description Update an agent's details (supervisor/admin only)
// @Tags supervisor
// @Accept json
// @Produce json
// @Security X-API-KEY
// @Param id path string true "Agent ID"
// @Param request body SupervisorUpdateAgentRequest true "Update details"
// @Success 200 {object} models.Agent
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /supervisor/agents/{id} [patch]
func (h *SupervisorHandler) UpdateAgent(c *gin.Context) {
	id := c.Param("id")
	var req SupervisorUpdateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	agent, err := h.agentService.Update(id, req.Name, req.Description, req.Role, req.Enabled)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, agent)
}

// SupervisorDeleteAgent godoc
// @Summary Delete an agent (supervisor only)
// @Description Delete an agent by ID (supervisor/admin only)
// @Tags supervisor
// @Produce json
// @Security X-API-KEY
// @Param id path string true "Agent ID"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /supervisor/agents/{id} [delete]
func (h *SupervisorHandler) DeleteAgent(c *gin.Context) {
	id := c.Param("id")

	if err := h.agentService.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Agent deleted successfully"})
}

// SupervisorUpdateTaskStatusRequest represents request to update any task status
type SupervisorUpdateTaskStatusRequest struct {
	Status models.TaskStatus `json:"status" binding:"required" example:"in_progress"`
	Comment string           `json:"comment" example:"Taking over this task"`
}

// SupervisorUpdateTaskStatus godoc
// @Summary Update any task status (supervisor only)
// @Description Update status of any task regardless of assignment (supervisor/admin only)
// @Tags supervisor
// @Accept json
// @Produce json
// @Security X-API-KEY
// @Param id path string true "Task ID"
// @Param request body SupervisorUpdateTaskStatusRequest true "Status update"
// @Success 200 {object} models.Task
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /supervisor/tasks/{id}/status [patch]
func (h *SupervisorHandler) UpdateTaskStatus(c *gin.Context) {
	taskID := c.Param("id")
	var req SupervisorUpdateTaskStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.taskService.UpdateStatus(taskID, req.Status, c.GetString("agent_name"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Add comment if provided
	if req.Comment != "" {
		commentService := services.NewCommentService()
		commentService.Create(taskID, c.GetString("agent_id"), "agent", c.GetString("agent_name"), req.Comment)
	}

	c.JSON(http.StatusOK, task)
}

// SupervisorAssignTaskRequest represents request to assign task to agent
type SupervisorAssignTaskRequest struct {
	AgentID string `json:"agent_id" binding:"required" example:"agent-uuid"`
}

// SupervisorAssignTask godoc
// @Summary Assign task to agent (supervisor only)
// @Description Assign any task to any agent (supervisor/admin only)
// @Tags supervisor
// @Accept json
// @Produce json
// @Security X-API-KEY
// @Param id path string true "Task ID"
// @Param request body SupervisorAssignTaskRequest true "Assignment"
// @Success 200 {object} models.Task
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /supervisor/tasks/{id}/assign [patch]
func (h *SupervisorHandler) AssignTask(c *gin.Context) {
	taskID := c.Param("id")
	var req SupervisorAssignTaskRequest

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

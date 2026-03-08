package handlers

import (
	"net/http"

	"github.com/formatho/agent-todo/services"
	"github.com/gin-gonic/gin"
)

type SupervisorHandler struct {
	agentService      *services.AgentService
	supervisorService *services.SupervisorService
}

func NewSupervisorHandler() *SupervisorHandler {
	return &SupervisorHandler{
		agentService:      services.NewAgentService(),
		supervisorService: services.NewSupervisorService(),
	}
}

// CreateAgent godoc
// @Summary Create a new agent
// @Description Create a new agent (supervisor/admin only)
// @Tags agents
// @Produce json
// @Security Bearer
// @Param request body CreateAgentRequest true "Agent details"
// @Success 201 {object} models.Agent
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /supervisor/agents [post]
func (h *SupervisorHandler) CreateAgent(c *gin.Context) {
	// Use agentHandler logic
	agentHandler := NewAgentHandler()
	agentHandler.CreateAgent(c)
}

// ListAgents godoc
// @Summary List all agents
// @Description List all agents (supervisor/admin only)
// @Tags agents
// @Produce json
// @Security Bearer
// @Success 200 {array} models.Agent
// @Failure 401 {object} map[string]string
// @Router /supervisor/agents [get]
func (h *SupervisorHandler) ListAgents(c *gin.Context) {
	// Use agentHandler logic
	agentHandler := NewAgentHandler()
	agentHandler.ListAgents(c)
}

// UpdateAgent godoc
// @Summary Update an agent
// @Description Update an agent (supervisor/admin only)
// @Tags agents
// @Produce json
// @Security Bearer
// @Param id path string true "Agent ID"
// @Param request body UpdateAgentRequest true "Update details"
// @Success 200 {object} models.Agent
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /supervisor/agents/{id} [patch]
func (h *SupervisorHandler) UpdateAgent(c *gin.Context) {
	// Use agentHandler logic
	agentHandler := NewAgentHandler()
	agentHandler.UpdateAgent(c)
}

// DeleteAgent godoc
// @Summary Delete an agent
// @Description Delete an agent (supervisor/admin only)
// @Tags agents
// @Produce json
// @Security Bearer
// @Param id path string true "Agent ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /supervisor/agents/{id} [delete]
func (h *SupervisorHandler) DeleteAgent(c *gin.Context) {
	// Use agentHandler logic
	agentHandler := NewAgentHandler()
	agentHandler.DeleteAgent(c)
}

// ListTasks godoc
// @Summary List all tasks
// @Description List all tasks (supervisor/admin only)
// @Tags tasks
// @Produce json
// @Security Bearer
// @Success 200 {array} models.Task
// @Failure 401 {object} map[string]string
// @Router /supervisor/tasks [get]
func (h *SupervisorHandler) ListTasks(c *gin.Context) {
	// Use taskHandler logic
	taskHandler := NewTaskHandler()
	taskHandler.ListTasks(c)
}

// UpdateTaskStatus godoc
// @Summary Update task status
// @Description Update any task status (supervisor/admin only)
// @Tags tasks
// @Produce json
// @Security Bearer
// @Param id path string true "Task ID"
// @Param request body UpdateStatusRequest true "Status update"
// @Success 200 {object} models.Task
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /supervisor/tasks/{id}/status [patch]
func (h *SupervisorHandler) UpdateTaskStatus(c *gin.Context) {
	// Use taskHandler logic
	taskHandler := NewTaskHandler()
	c.Set("isSupervisor", true)
	taskHandler.UpdateTask(c)
}

// AssignTask godoc
// @Summary Assign task to agent
// @Description Assign any task to an agent (supervisor/admin only)
// @Tags tasks
// @Produce json
// @Security Bearer
// @Param id path string true "Task ID"
// @Param request body AssignRequest true "Agent assignment"
// @Success 200 {object} models.Task
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /supervisor/tasks/{id}/assign [patch]
func (h *SupervisorHandler) AssignTask(c *gin.Context) {
	// Use taskHandler logic
	taskHandler := NewTaskHandler()
	taskHandler.AssignAgent(c)
}

// GetAgentsWithTasks godoc
// @Summary Get agents with their active tasks
// @Description Get all agents along with their in-progress tasks
// @Tags agents
// @Produce json
// @Security Bearer
// @Success 200 {array} services.AgentWithTasks
// @Failure 401 {object} map[string]string
// @Router /supervisor/agents/activity [get]
func (h *SupervisorHandler) GetAgentsWithTasks(c *gin.Context) {
	agents, err := h.supervisorService.GetAgentsWithActiveTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, agents)
}

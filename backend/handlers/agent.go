package handlers

import (
	"net/http"

	"github.com/formatho/agent-todo/services"
	"github.com/gin-gonic/gin"
)

type AgentHandler struct {
	agentService *services.AgentService
}

func NewAgentHandler() *AgentHandler {
	return &AgentHandler{
		agentService: services.NewAgentService(),
	}
}

// CreateAgentRequest represents the request body for creating an agent
type CreateAgentRequest struct {
	Name        string `json:"name" binding:"required" example:"My Agent"`
	Description string `json:"description" example:"An AI assistant agent"`
}

// CreateAgent godoc
// @Summary Create a new agent
// @Description Create a new AI agent with an API key
// @Tags agents
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body CreateAgentRequest true "Agent details"
// @Success 201 {object} models.Agent
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /agents [post]
func (h *AgentHandler) CreateAgent(c *gin.Context) {
	var req CreateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	agent, err := h.agentService.Create(req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, agent)
}

// ListAgents godoc
// @Summary List all agents
// @Description Get a list of all agents
// @Tags agents
// @Produce json
// @Security Bearer
// @Success 200 {array} models.Agent
// @Failure 401 {object} map[string]string
// @Router /agents [get]
func (h *AgentHandler) ListAgents(c *gin.Context) {
	agents, err := h.agentService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, agents)
}

// GetAgent godoc
// @Summary Get an agent
// @Description Get a specific agent by ID
// @Tags agents
// @Produce json
// @Security Bearer
// @Param id path string true "Agent ID"
// @Success 200 {object} models.Agent
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /agents/{id} [get]
func (h *AgentHandler) GetAgent(c *gin.Context) {
	id := c.Param("id")

	agent, err := h.agentService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	c.JSON(http.StatusOK, agent)
}

// DeleteAgent godoc
// @Summary Delete an agent
// @Description Delete an agent by ID
// @Tags agents
// @Produce json
// @Security Bearer
// @Param id path string true "Agent ID"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /agents/{id} [delete]
func (h *AgentHandler) DeleteAgent(c *gin.Context) {
	id := c.Param("id")

	if err := h.agentService.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Agent deleted successfully"})
}

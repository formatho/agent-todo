package handlers

import (
	"net/http"

	"github.com/formatho/agent-todo/middleware"
	"github.com/formatho/agent-todo/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CommentHandler struct {
	taskService *services.TaskService
}

func NewCommentHandler() *CommentHandler {
	return &CommentHandler{
		taskService: services.NewTaskService(),
	}
}

// CreateCommentRequest represents the request body for creating a comment
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required" example:"This task is in progress"`
}

// CreateComment godoc
// @Summary Create a comment on a task
// @Description Add a comment to a task
// @Tags comments
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Task ID"
// @Param request body CreateCommentRequest true "Comment content"
// @Success 201 {object} models.TaskComment
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tasks/{id}/comments [post]
func (h *CommentHandler) CreateComment(c *gin.Context) {
	taskID := c.Param("id")
	var req CreateCommentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if it's a user or agent
	userID, userErr := middleware.GetUserID(c)
	authorName := ""
	authorID := uuid.Nil
	authorType := "user"

	if userErr == nil {
		// It's a user
		userService := services.NewUserService()
		user, err := userService.GetByID(userID)
		if err == nil {
			authorName = user.Email
			authorID = user.ID
		}
	} else {
		// Check if it's an agent
		agentID, agentErr := middleware.GetAgentID(c)
		if agentErr == nil {
			agentService := services.NewAgentService()
			agent, err := agentService.GetByID(agentID)
			if err == nil {
				authorName = agent.Name
				authorID = agent.ID
				authorType = "agent"
			}
		}
	}

	if authorName == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	comment, err := h.taskService.AddComment(
		taskID,
		req.Content,
		authorID,
		authorType,
		authorName,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// GetComments godoc
// @Summary Get task comments
// @Description Get all comments for a task
// @Tags comments
// @Produce json
// @Security Bearer
// @Param id path string true "Task ID"
// @Success 200 {array} models.TaskComment
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tasks/{id}/comments [get]
func (h *CommentHandler) GetComments(c *gin.Context) {
	taskID := c.Param("id")

	// Verify task exists
	_, err := h.taskService.GetByID(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	comments, err := h.taskService.GetComments(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// AgentCreateComment godoc
// @Summary Create a comment (Agent)
// @Description Add a comment to a task (for agents)
// @Tags agent
// @Accept json
// @Produce json
// @Security X-API-KEY
// @Param id path string true "Task ID"
// @Param request body CreateCommentRequest true "Comment content"
// @Success 201 {object} models.TaskComment
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /agent/tasks/{id}/comments [post]
func (h *CommentHandler) AgentCreateComment(c *gin.Context) {
	taskID := c.Param("id")
	var req CreateCommentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	// Verify task is assigned to this agent
	task, err := h.taskService.GetByID(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if task.AssignedAgentID == nil || task.AssignedAgentID.String() != agentID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Task not assigned to this agent"})
		return
	}

	comment, err := h.taskService.AddComment(
		taskID,
		req.Content,
		uuid.MustParse(agentID),
		"agent",
		agentName,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// AgentGetComments godoc
// @Summary Get task comments (Agent)
// @Description Get all comments for a task (for agents)
// @Tags agent
// @Produce json
// @Security X-API-KEY
// @Param id path string true "Task ID"
// @Success 200 {array} models.TaskComment
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /agent/tasks/{id}/comments [get]
func (h *CommentHandler) AgentGetComments(c *gin.Context) {
	taskID := c.Param("id")
	agentID, err := middleware.GetAgentID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Verify task is assigned to this agent
	task, err := h.taskService.GetByID(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if task.AssignedAgentID == nil || task.AssignedAgentID.String() != agentID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Task not assigned to this agent"})
		return
	}

	comments, err := h.taskService.GetComments(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comments)
}

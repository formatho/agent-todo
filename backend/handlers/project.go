package handlers

import (
	"net/http"

	"github.com/formatho/agent-todo/middleware"
	"github.com/formatho/agent-todo/models"
	"github.com/formatho/agent-todo/services"
	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	projectService *services.ProjectService
}

func NewProjectHandler() *ProjectHandler {
	return &ProjectHandler{
		projectService: services.NewProjectService(),
	}
}

// CreateProjectRequest represents the request body for creating a project
type CreateProjectRequest struct {
	Name            string `json:"name" binding:"required" example:"Website Redesign"`
	Description     string `json:"description" example:"Redesign the company website with new branding"`
	RepositoryURL   string `json:"repository_url" example:"https://github.com/org/repo"`
	DeployedURL     string `json:"deployed_url" example:"https://app.example.com"`
	DocumentationURL string `json:"documentation_url" example:"https://docs.example.com"`
	LLMContext      string `json:"llm_context" example:"Build a REST API with Go and PostgreSQL..."`
}

// UpdateProjectRequest represents the request body for updating a project
type UpdateProjectRequest struct {
	Name            *string               `json:"name" example:"Updated project name"`
	Description     *string               `json:"description" example:"Updated description"`
	Status          *models.ProjectStatus `json:"status" example:"active"`
	RepositoryURL   *string               `json:"repository_url" example:"https://github.com/org/repo"`
	DeployedURL     *string               `json:"deployed_url" example:"https://app.example.com"`
	DocumentationURL *string              `json:"documentation_url" example:"https://docs.example.com"`
	LLMContext      *string               `json:"llm_context" example:"Updated LLM instructions..."`
}

// CreateProject godoc
// @Summary Create a new project
// @Description Create a new project
// @Tags projects
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body CreateProjectRequest true "Project details"
// @Success 201 {object} models.Project
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /projects [post]
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Build project input
	input := services.ProjectInput{
		Name:             &req.Name,
		Description:      &req.Description,
		RepositoryURL:    stringPtr(req.RepositoryURL),
		DeployedURL:      stringPtr(req.DeployedURL),
		DocumentationURL: stringPtr(req.DocumentationURL),
		LLMContext:       stringPtr(req.LLMContext),
	}

	var project *models.Project

	// Use organisation context if available
	if orgID, exists := c.Get("organisation_id"); exists {
		if orgIDStr, ok := orgID.(string); ok && orgIDStr != "" {
			project, err = h.projectService.CreateProject(input, userID, &orgIDStr)
		}
	}

	if project == nil {
		project, err = h.projectService.CreateProject(input, userID, nil)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, project)
}

// Helper function to convert empty string to nil
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// ListProjects godoc
// @Summary List projects
// @Description Get a list of projects with optional filters
// @Tags projects
// @Produce json
// @Security Bearer
// @Param status query string false "Filter by status" Enums(active, archived, completed)
// @Param search query string false "Search in name and description"
// @Success 200 {array} models.Project
// @Failure 401 {object} map[string]string
// @Router /projects [get]
func (h *ProjectHandler) ListProjects(c *gin.Context) {
	filter := services.ProjectFilter{
		Status:     models.ProjectStatus(c.Query("status")),
		SearchTerm: c.Query("search"),
	}

	// Add organisation filter if present in context
	if orgID, exists := c.Get("organisation_id"); exists {
		if orgIDStr, ok := orgID.(string); ok && orgIDStr != "" {
			filter.OrganisationID = orgIDStr
		}
	}

	projects, err := h.projectService.List(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projects)
}

// GetProject godoc
// @Summary Get a project
// @Description Get a specific project by ID
// @Tags projects
// @Produce json
// @Security Bearer
// @Param id path string true "Project ID"
// @Success 200 {object} models.Project
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /projects/{id} [get]
func (h *ProjectHandler) GetProject(c *gin.Context) {
	id := c.Param("id")

	project, err := h.projectService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, project)
}

// UpdateProject godoc
// @Summary Update a project
// @Description Update a project's details
// @Tags projects
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Project ID"
// @Param request body UpdateProjectRequest true "Project updates"
// @Success 200 {object} models.Project
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /projects/{id} [patch]
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	id := c.Param("id")

	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Build project input
	input := services.ProjectInput{
		Name:             req.Name,
		Description:      req.Description,
		Status:           req.Status,
		RepositoryURL:    req.RepositoryURL,
		DeployedURL:      req.DeployedURL,
		DocumentationURL: req.DocumentationURL,
		LLMContext:       req.LLMContext,
	}

	var project *models.Project
	var err error

	// Use organisation context if available
	if orgID, exists := c.Get("organisation_id"); exists {
		if orgIDStr, ok := orgID.(string); ok && orgIDStr != "" {
			project, err = h.projectService.UpdateProjectByOrganisation(id, orgIDStr, input)
		}
	}

	if project == nil {
		project, err = h.projectService.UpdateProject(id, input)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, project)
}

// DeleteProject godoc
// @Summary Delete a project
// @Description Delete a project by ID
// @Tags projects
// @Produce json
// @Security Bearer
// @Param id path string true "Project ID"
// @Success 204
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /projects/{id} [delete]
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	id := c.Param("id")

	if err := h.projectService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetProjectTasks godoc
// @Summary Get project tasks
// @Description Get all tasks for a specific project
// @Tags projects
// @Produce json
// @Security Bearer
// @Param id path string true "Project ID"
// @Success 200 {array} models.Task
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /projects/{id}/tasks [get]
func (h *ProjectHandler) GetProjectTasks(c *gin.Context) {
	id := c.Param("id")

	tasks, err := h.projectService.GetTasks(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// ListProjectsForAgent godoc
// @Summary List projects (for agents)
// @Description Get a list of projects with optional filters (agent accessible)
// @Tags agent
// @Produce json
// @Security X-API-KEY
// @Param status query string false "Filter by status" Enums(active, archived, completed)
// @Param search query string false "Search in name and description"
// @Success 200 {array} models.Project
// @Failure 401 {object} map[string]string
// @Router /agent/projects [get]
func (h *ProjectHandler) ListProjectsForAgent(c *gin.Context) {
	filter := services.ProjectFilter{
		Status:     models.ProjectStatus(c.Query("status")),
		SearchTerm: c.Query("search"),
	}

	projects, err := h.projectService.List(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projects)
}

// GetProjectForAgent godoc
// @Summary Get a project (for agents)
// @Description Get a specific project by ID (agent accessible)
// @Tags agent
// @Produce json
// @Security X-API-KEY
// @Param id path string true "Project ID"
// @Success 200 {object} models.Project
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /agent/projects/{id} [get]
func (h *ProjectHandler) GetProjectForAgent(c *gin.Context) {
	id := c.Param("id")

	project, err := h.projectService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, project)
}

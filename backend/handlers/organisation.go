package handlers

import (
	"net/http"

	"github.com/formatho/agent-todo/middleware"
	"github.com/formatho/agent-todo/models"
	"github.com/formatho/agent-todo/services"
	"github.com/gin-gonic/gin"
)

type OrganisationHandler struct {
	orgService *services.OrganisationService
}

func NewOrganisationHandler() *OrganisationHandler {
	return &OrganisationHandler{
		orgService: services.NewOrganisationService(),
	}
}

// CreateOrganisation godoc
// @Summary Create a new organisation
// @Description Create a new organisation with the authenticated user as owner
// @Tags organisations
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body services.CreateOrganisationInput true "Organisation details"
// @Success 201 {object} models.Organisation
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /organisations [post]
func (h *OrganisationHandler) CreateOrganisation(c *gin.Context) {
	var req services.CreateOrganisationInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	req.UserID = userID

	org, err := h.orgService.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, org)
}

// ListOrganisations godoc
// @Summary List organisations
// @Description Get a list of organisations the authenticated user belongs to
// @Tags organisations
// @Produce json
// @Security Bearer
// @Param status query string false "Filter by status" Enums(active, suspended, archived)
// @Param search query string false "Search in name and description"
// @Success 200 {array} models.Organisation
// @Failure 401 {object} map[string]string
// @Router /organisations [get]
func (h *OrganisationHandler) ListOrganisations(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	filter := services.OrganisationFilter{
		Status:     models.OrganisationStatus(c.Query("status")),
		SearchTerm: c.Query("search"),
	}

	orgs, err := h.orgService.List(userID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orgs)
}

// GetOrganisation godoc
// @Summary Get an organisation
// @Description Get a specific organisation by ID
// @Tags organisations
// @Produce json
// @Security Bearer
// @Param id path string true "Organisation ID"
// @Success 200 {object} models.Organisation
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /organisations/{id} [get]
func (h *OrganisationHandler) GetOrganisation(c *gin.Context) {
	orgID := c.Param("id")
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Check if user is a member
	isMember, err := h.orgService.IsMember(orgID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	org, err := h.orgService.GetByID(orgID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "organisation not found"})
		return
	}

	c.JSON(http.StatusOK, org)
}

// UpdateOrganisation godoc
// @Summary Update an organisation
// @Description Update an organisation's details (requires admin or owner role)
// @Tags organisations
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Organisation ID"
// @Param request body services.UpdateOrganisationInput true "Organisation updates"
// @Success 200 {object} models.Organisation
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /organisations/{id} [patch]
func (h *OrganisationHandler) UpdateOrganisation(c *gin.Context) {
	orgID := c.Param("id")
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Check if user is admin or owner
	isAdmin, err := h.orgService.IsAdmin(orgID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin or owner role required"})
		return
	}

	var req services.UpdateOrganisationInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	org, err := h.orgService.Update(orgID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, org)
}

// DeleteOrganisation godoc
// @Summary Delete an organisation
// @Description Soft delete an organisation (requires owner role)
// @Tags organisations
// @Produce json
// @Security Bearer
// @Param id path string true "Organisation ID"
// @Success 204
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /organisations/{id} [delete]
func (h *OrganisationHandler) DeleteOrganisation(c *gin.Context) {
	orgID := c.Param("id")
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Check if user is owner
	isOwner, err := h.orgService.IsOwner(orgID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "owner role required"})
		return
	}

	if err := h.orgService.Delete(orgID); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// AddOrganisationMember godoc
// @Summary Add a member to organisation
// @Description Add a user to the organisation (requires admin or owner role)
// @Tags organisations
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Organisation ID"
// @Param request body services.AddMemberInput true "Member details"
// @Success 201 {object} models.OrganisationMember
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /organisations/{id}/members [post]
func (h *OrganisationHandler) AddOrganisationMember(c *gin.Context) {
	orgID := c.Param("id")
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Check if user is admin or owner
	isAdmin, err := h.orgService.IsAdmin(orgID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin or owner role required"})
		return
	}

	var req services.AddMemberInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	member, err := h.orgService.AddMember(orgID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, member)
}

// UpdateMemberRole godoc
// @Summary Update member role
// @Description Update a member's role in the organisation (requires owner role)
// @Tags organisations
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Organisation ID"
// @Param member_id path string true "Member ID"
// @Param request body services.UpdateMemberRoleInput true "Role update"
// @Success 200 {object} models.OrganisationMember
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /organisations/{id}/members/{member_id} [patch]
func (h *OrganisationHandler) UpdateMemberRole(c *gin.Context) {
	orgID := c.Param("id")
	memberID := c.Param("member_id")
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Check if user is owner
	isOwner, err := h.orgService.IsOwner(orgID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "owner role required"})
		return
	}

	var req services.UpdateMemberRoleInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	member, err := h.orgService.UpdateMemberRole(orgID, memberID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, member)
}

// RemoveOrganisationMember godoc
// @Summary Remove a member
// @Description Remove a member from the organisation (requires admin or owner role)
// @Tags organisations
// @Produce json
// @Security Bearer
// @Param id path string true "Organisation ID"
// @Param member_id path string true "Member ID"
// @Success 204
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /organisations/{id}/members/{member_id} [delete]
func (h *OrganisationHandler) RemoveOrganisationMember(c *gin.Context) {
	orgID := c.Param("id")
	memberID := c.Param("member_id")
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Check if user is admin or owner
	isAdmin, err := h.orgService.IsAdmin(orgID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin or owner role required"})
		return
	}

	if err := h.orgService.RemoveMember(orgID, memberID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// LeaveOrganisation godoc
// @Summary Leave organisation
// @Description Leave an organisation (not available for owners)
// @Tags organisations
// @Produce json
// @Security Bearer
// @Param id path string true "Organisation ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /organisations/{id}/leave [post]
func (h *OrganisationHandler) LeaveOrganisation(c *gin.Context) {
	orgID := c.Param("id")
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := h.orgService.LeaveOrganisation(orgID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

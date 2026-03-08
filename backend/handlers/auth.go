package handlers

import (
	"net/http"
	"os"

	apperrors "github.com/formatho/agent-todo/errors"
	"github.com/formatho/agent-todo/middleware"
	"github.com/formatho/agent-todo/models"
	"github.com/formatho/agent-todo/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userService *services.UserService
	jwtService  *services.JWTService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		userService: services.NewUserService(),
		jwtService:  services.NewJWTService(),
	}
}

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// AuthResponse represents the response for successful authentication
type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} map[string]string
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	// Check if registration is disabled
	if os.Getenv("DISABLE_REGISTRATION") == "true" {
		middleware.HandleError(c, apperrors.Forbidden("Registration is currently disabled"))
		return
	}

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, apperrors.BadRequest("Invalid request body").WithDetails(map[string]interface{}{
			"validation_error": err.Error(),
		}))
		return
	}

	user, err := h.userService.Register(req.Email, req.Password)
	if err != nil {
		middleware.HandleError(c, apperrors.BadRequest(err.Error()))
		return
	}

	// New users don't have an organisation yet
	token, err := h.jwtService.GenerateTokenWithOrg(user.ID.String(), user.Email, "")
	if err != nil {
		middleware.HandleError(c, apperrors.InternalServerError("Failed to generate token").WithInternal(err))
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User:  *user,
	})
}

// Login godoc
// @Summary Login user
// @Description Authenticate with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} AuthResponse
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, apperrors.BadRequest("Invalid request body"))
		return
	}

	user, err := h.userService.Login(req.Email, req.Password)
	if err != nil {
		middleware.HandleError(c, apperrors.Unauthorized(err.Error()))
		return
	}

	// Get organisation ID for token (use current_org_id if set)
	orgID := ""
	if user.CurrentOrgID != nil {
		orgID = user.CurrentOrgID.String()
	}

	token, err := h.jwtService.GenerateTokenWithOrg(user.ID.String(), user.Email, orgID)
	if err != nil {
		middleware.HandleError(c, apperrors.InternalServerError("Failed to generate token").WithInternal(err))
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User:  *user,
	})
}

// GetCurrentUser godoc
// @Summary Get current user
// @Description Get the currently authenticated user's information
// @Tags auth
// @Produce json
// @Security Bearer
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]string
// @Router /auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// SwitchOrganisationRequest represents the request to switch organisation
type SwitchOrganisationRequest struct {
	OrganisationID string `json:"organisation_id" binding:"required,uuid"`
}

// SwitchOrganisationResponse represents the response after switching organisation
type SwitchOrganisationResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// SwitchOrganisation godoc
// @Summary Switch current organisation
// @Description Switch the user's current organisation context
// @Tags auth
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body SwitchOrganisationRequest true "Organisation ID"
// @Success 200 {object} SwitchOrganisationResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /auth/switch-organisation [post]
func (h *AuthHandler) SwitchOrganisation(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req SwitchOrganisationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Verify user is a member of the organisation
	orgService := services.NewOrganisationService()
	isMember, err := orgService.IsMember(req.OrganisationID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify membership"})
		return
	}
	if !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not a member of this organisation"})
		return
	}

	// Update user's current organisation
	user, err := h.userService.SetCurrentOrganisation(userID, req.OrganisationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update current organisation"})
		return
	}

	// Generate new token with updated organisation
	token, err := h.jwtService.GenerateTokenWithOrg(user.ID.String(), user.Email, req.OrganisationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, SwitchOrganisationResponse{
		Token: token,
		User:  *user,
	})
}

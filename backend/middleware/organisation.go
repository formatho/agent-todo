package middleware

import (
	"errors"
	"net/http"

	"github.com/formatho/agent-todo/models"
	"github.com/formatho/agent-todo/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// OrganisationMiddleware enforces organisation isolation for user requests
// This middleware should be used after AuthMiddleware
func OrganisationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by AuthMiddleware)
		userID, err := GetUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// Get organisation ID from JWT claims (if present)
		orgID, exists := c.Get("organisation_id")
		if exists {
			c.Set("current_organisation_id", orgID)
			c.Next()
			return
		}

		// If no organisation_id in token, try to get user's current organisation from DB
		// This handles users who logged in before organisation support was added
		userService := services.NewUserService()
		user, err := userService.GetByID(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load user"})
			c.Abort()
			return
		}

		if user.CurrentOrgID == nil {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "No organisation context. Please select an organisation.",
			})
			c.Abort()
			return
		}

		// Set organisation ID from user's current organisation
		c.Set("current_organisation_id", user.CurrentOrgID.String())
		c.Next()
	}
}

// AgentOrganisationMiddleware enforces organisation isolation for agent requests
// This middleware should be used after AgentAuthMiddleware
func AgentOrganisationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get agent from context (set by AgentAuthMiddleware)
		agent, err := GetAgent(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Agent not authenticated"})
			c.Abort()
			return
		}

		// Get organisation ID from agent
		if agent.OrganisationID == nil {
			// Agent doesn't belong to any organisation
			// In production, all agents should belong to an organisation
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Agent is not associated with any organisation",
			})
			c.Abort()
			return
		}

		// Set organisation ID in context for use in handlers
		c.Set("current_organisation_id", agent.OrganisationID.String())
		c.Next()
	}
}

// GetOrganisationID retrieves organisation ID from context
func GetOrganisationID(c *gin.Context) (string, error) {
	orgID, exists := c.Get("current_organisation_id")
	if !exists {
		return "", errors.New("no organisation context")
	}
	return orgID.(string), nil
}

// GetOrganisationUUID retrieves organisation ID as UUID from context
func GetOrganisationUUID(c *gin.Context) (uuid.UUID, error) {
	orgIDStr, err := GetOrganisationID(c)
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.Parse(orgIDStr)
}

// RequireOrganisationRole checks if user has required role in organisation
// Role hierarchy: owner > admin > member
func RequireOrganisationRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context
		userID, err := GetUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// Get organisation ID from context
		orgID, err := GetOrganisationID(c)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Organisation context required"})
			c.Abort()
			return
		}

		// Get user's role in organisation
		orgService := services.NewOrganisationService()
		userRole, err := orgService.GetMemberRole(orgID, userID)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "User is not a member of this organisation"})
			c.Abort()
			return
		}

		// Check if user has required role (or higher)
		if !hasRequiredRole(*userRole, requiredRole) {
			c.JSON(http.StatusForbidden, gin.H{
				"error":         "Insufficient permissions",
				"required_role": requiredRole,
				"current_role":  string(*userRole),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// hasRequiredRole checks if the user's role meets the required role level
func hasRequiredRole(userRole models.OrganisationMemberRole, requiredRole string) bool {
	// Role hierarchy map (higher number = more permissions)
	roleLevel := map[models.OrganisationMemberRole]int{
		models.OrganisationMemberRoleMember: 1,
		models.OrganisationMemberRoleAdmin:  2,
		models.OrganisationMemberRoleOwner:  3,
	}

	requiredLevel, exists := roleLevel[models.OrganisationMemberRole(requiredRole)]
	if !exists {
		return false
	}

	userLevel := roleLevel[userRole]
	return userLevel >= requiredLevel
}

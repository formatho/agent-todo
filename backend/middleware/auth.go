package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/formatho/agent-todo/models"
	"github.com/formatho/agent-todo/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID         string `json:"user_id"`
	Email          string `json:"email"`
	OrganisationID string `json:"organisation_id,omitempty"`
	jwt.RegisteredClaims
}

// AuthMiddleware validates JWT tokens for human users
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		claims := &Claims{}
		jwtSecret := []byte(os.Getenv("JWT_SECRET"))

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)

		// Set organisation_id if present in token
		if claims.OrganisationID != "" {
			c.Set("organisation_id", claims.OrganisationID)
		}

		c.Next()
	}
}

// AgentAuthMiddleware validates API keys for agents
func AgentAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-KEY")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API key required"})
			c.Abort()
			return
		}

		agentService := services.NewAgentService()
		agent, err := agentService.GetByAPIKey(apiKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}

		// Check if agent is enabled
		if !agent.Enabled {
			c.JSON(http.StatusForbidden, gin.H{"error": "Agent is disabled"})
			c.Abort()
			return
		}

		c.Set("agent_id", agent.ID.String())
		c.Set("agent_name", agent.Name)
		c.Set("agent", agent)
		c.Next()
	}
}

// GetUserID retrieves user ID from context
func GetUserID(c *gin.Context) (string, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", errors.New("user not authenticated")
	}
	return userID.(string), nil
}

// GetAgentID retrieves agent ID from context
func GetAgentID(c *gin.Context) (string, error) {
	agentID, exists := c.Get("agent_id")
	if !exists {
		return "", errors.New("agent not authenticated")
	}
	return agentID.(string), nil
}

// GetAgentName retrieves agent name from context
func GetAgentName(c *gin.Context) (string, error) {
	agentName, exists := c.Get("agent_name")
	if !exists {
		return "", errors.New("agent not authenticated")
	}
	return agentName.(string), nil
}

// GetAgent retrieves agent object from context
func GetAgent(c *gin.Context) (*models.Agent, error) {
	agent, exists := c.Get("agent")
	if !exists {
		return nil, errors.New("agent not authenticated")
	}
	return agent.(*models.Agent), nil
}

// RequireSupervisor middleware checks if agent has supervisor or admin role
func RequireSupervisor() gin.HandlerFunc {
	return func(c *gin.Context) {
		agent, err := GetAgent(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Agent not authenticated"})
			c.Abort()
			return
		}

		if agent.Role != models.AgentRoleSupervisor && agent.Role != models.AgentRoleAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions. Supervisor or admin role required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

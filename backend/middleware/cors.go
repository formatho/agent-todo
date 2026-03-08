package middleware

import (
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware returns CORS configuration
func CORSMiddleware() gin.HandlerFunc {
	// Get allowed origins from environment variable
	allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")

	var allowOrigins []string
	if allowedOriginsEnv != "" {
		// Split comma-separated origins
		allowOrigins = strings.Split(allowedOriginsEnv, ",")
		// Trim whitespace from each origin
		for i, origin := range allowOrigins {
			allowOrigins[i] = strings.TrimSpace(origin)
		}
	} else {
		// Default to localhost for development
		allowOrigins = []string{
			"http://localhost:3000",
			"http://localhost:8080",
			"http://localhost:5173", // Vite default
		}
	}

	config := cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-API-KEY"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}

	return cors.New(config)
}

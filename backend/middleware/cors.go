package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware returns CORS configuration
func CORSMiddleware() gin.HandlerFunc {
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-API-KEY"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}

	return cors.New(config)
}

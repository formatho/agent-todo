package main

import (
	"log"
	"os"

	"github.com/formatho/agent-todo/db"
	"github.com/formatho/agent-todo/handlers"
	"github.com/formatho/agent-todo/middleware"
	_ "github.com/formatho/agent-todo/docs"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Agent Todo API
// @version 1.0
// @description A task management system for human and AI agent collaboration
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description JWT token for human users

// @securityDefinitions.apikey X-API-KEY
// @in header
// @name X-API-KEY
// @description API key for AI agents
func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Connect to database
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://agent_todo:agent_todo_pass@localhost:5432/agent_todo?sslmode=disable"
	}

	if err := db.Connect(databaseURL); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Set Gin mode
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Create router
	router := gin.Default()

	// Apply middleware
	router.Use(middleware.CORSMiddleware())

	// Initialize handlers
	authHandler := handlers.NewAuthHandler()
	agentHandler := handlers.NewAgentHandler()
	projectHandler := handlers.NewProjectHandler()
	taskHandler := handlers.NewTaskHandler()
	agentTaskHandler := handlers.NewAgentTaskHandler()
	commentHandler := handlers.NewCommentHandler()
	toolsHandler := handlers.NewToolsHandler()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Swagger documentation
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Auth routes
	auth := router.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.GET("/me", middleware.AuthMiddleware(), authHandler.GetCurrentUser)
	}

	// Project routes (human)
	projects := router.Group("/projects")
	projects.Use(middleware.AuthMiddleware())
	{
		projects.POST("", projectHandler.CreateProject)
		projects.GET("", projectHandler.ListProjects)
		projects.GET("/:id", projectHandler.GetProject)
		projects.PATCH("/:id", projectHandler.UpdateProject)
		projects.DELETE("/:id", projectHandler.DeleteProject)
		projects.GET("/:id/tasks", projectHandler.GetProjectTasks)
	}

	// Agent management routes (human only)
	agents := router.Group("/agents")
	agents.Use(middleware.AuthMiddleware())
	{
		agents.POST("", agentHandler.CreateAgent)
		agents.GET("", agentHandler.ListAgents)
		agents.GET("/:id", agentHandler.GetAgent)
		agents.DELETE("/:id", agentHandler.DeleteAgent)
	}

	// Task routes (human)
	tasks := router.Group("/tasks")
	tasks.Use(middleware.AuthMiddleware())
	{
		tasks.POST("", taskHandler.CreateTask)
		tasks.GET("", taskHandler.ListTasks)
		tasks.GET("/:id", taskHandler.GetTask)
		tasks.PATCH("/:id", taskHandler.UpdateTask)
		tasks.DELETE("/:id", taskHandler.DeleteTask)
		tasks.PATCH("/:id/assign", taskHandler.AssignAgent)
		tasks.PATCH("/:id/unassign", taskHandler.UnassignAgent)
		tasks.GET("/:id/comments", commentHandler.GetComments)
		tasks.POST("/:id/comments", commentHandler.CreateComment)
	}

	// Agent task routes (agent only)
	agentTasks := router.Group("/agent")
	agentTasks.Use(middleware.AgentAuthMiddleware())
	{
		agentTasks.POST("/tasks", agentTaskHandler.CreateTask)
		agentTasks.GET("/tasks", agentTaskHandler.ListTasks)
		agentTasks.GET("/tasks/:id", agentTaskHandler.GetTask)
		agentTasks.PATCH("/tasks/:id/status", agentTaskHandler.UpdateStatus)
		agentTasks.GET("/tasks/:id/comments", commentHandler.AgentGetComments)
		agentTasks.POST("/tasks/:id/comments", commentHandler.AgentCreateComment)
	}

	// OpenClaw tool endpoints (agent only)
	tools := router.Group("/tools")
	tools.Use(middleware.AgentAuthMiddleware())
	{
		tools.POST("/tasks/create", toolsHandler.CreateTask)
		tools.POST("/tasks/update", toolsHandler.UpdateTask)
		tools.POST("/tasks/list", toolsHandler.ListTasks)
		tools.GET("/tasks/status/:id", toolsHandler.GetStatus)
	}

	// Start server
	port := "8080"
	log.Printf("Server starting on port %s...", port)
	log.Printf("Swagger documentation available at http://localhost:%s/docs/index.html", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

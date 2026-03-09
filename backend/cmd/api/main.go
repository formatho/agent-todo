package main

import (
	"log"
	"os"

	"github.com/formatho/agent-todo/db"
	_ "github.com/formatho/agent-todo/docs"
	"github.com/formatho/agent-todo/handlers"
	"github.com/formatho/agent-todo/middleware"
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
	router.Use(middleware.RateLimitMiddleware(100)) // 100 requests per minute per IP

	// Initialize handlers
	authHandler := handlers.NewAuthHandler()
	agentHandler := handlers.NewAgentHandler()
	projectHandler := handlers.NewProjectHandler()
	taskHandler := handlers.NewTaskHandler()
	agentTaskHandler := handlers.NewAgentTaskHandler()
	commentHandler := handlers.NewCommentHandler()
	toolsHandler := handlers.NewToolsHandler()
	supervisorHandler := handlers.NewSupervisorHandler()
	activityHandler := handlers.NewActivityHandler()
	reminderHandler := handlers.NewReminderHandler()
	organisationHandler := handlers.NewOrganisationHandler()
	subtaskHandler := handlers.NewSubtaskHandler()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Reminder endpoints (for cron/scheduled tasks)
	reminders := router.Group("/reminders")
	{
		reminders.GET("/upcoming", reminderHandler.GetUpcomingDueTasks)
		reminders.GET("/overdue", reminderHandler.GetOverdueTasks)
	}

	// Swagger documentation
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Auth routes (stricter rate limiting)
	auth := router.Group("/auth")
	auth.Use(middleware.RateLimitMiddleware(20)) // 20 requests per minute for auth
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.GET("/me", middleware.AuthMiddleware(), authHandler.GetCurrentUser)
		auth.POST("/switch-organisation", middleware.AuthMiddleware(), authHandler.SwitchOrganisation)
	}

	// Organisation routes
	organisations := router.Group("/organisations")
	organisations.Use(middleware.AuthMiddleware())
	{
		organisations.POST("", organisationHandler.CreateOrganisation)
		organisations.GET("", organisationHandler.ListOrganisations)
		organisations.GET("/:id", organisationHandler.GetOrganisation)
		organisations.PATCH("/:id", organisationHandler.UpdateOrganisation)
		organisations.DELETE("/:id", organisationHandler.DeleteOrganisation)
		organisations.POST("/:id/members", organisationHandler.AddOrganisationMember)
		organisations.PATCH("/:id/members/:member_id", organisationHandler.UpdateMemberRole)
		organisations.DELETE("/:id/members/:member_id", organisationHandler.RemoveOrganisationMember)
		organisations.POST("/:id/leave", organisationHandler.LeaveOrganisation)
	}

	// Project routes (human)
	projects := router.Group("/projects")
	projects.Use(middleware.AuthMiddleware())
	projects.Use(middleware.OrganisationMiddleware())
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
	agents.Use(middleware.OrganisationMiddleware())
	{
		agents.POST("", agentHandler.CreateAgent)
		agents.GET("", agentHandler.ListAgents)
		agents.GET("/activity", agentHandler.GetAgentsWithTasks)
		agents.GET("/:id", agentHandler.GetAgent)
		agents.GET("/:id/statistics", agentHandler.GetAgentStatistics)
		agents.PATCH("/:id", agentHandler.UpdateAgent)
		agents.DELETE("/:id", agentHandler.DeleteAgent)
	}

	// Task routes (human)
	tasks := router.Group("/tasks")
	tasks.Use(middleware.AuthMiddleware())
	tasks.Use(middleware.OrganisationMiddleware())
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
		// Subtask routes
		tasks.GET("/:task_id/subtasks", subtaskHandler.ListSubtasks)
		tasks.POST("/:task_id/subtasks", subtaskHandler.CreateSubtask)
		tasks.POST("/:task_id/subtasks/reorder", subtaskHandler.ReorderSubtasks)
	}

	// Subtask routes (human - individual subtask operations)
	subtasks := router.Group("/subtasks")
	subtasks.Use(middleware.AuthMiddleware())
	subtasks.Use(middleware.OrganisationMiddleware())
	{
		subtasks.GET("/:id", subtaskHandler.GetSubtask)
		subtasks.PATCH("/:id", subtaskHandler.UpdateSubtask)
		subtasks.DELETE("/:id", subtaskHandler.DeleteSubtask)
	}

	// Activity feed routes
	activity := router.Group("/activity")
	activity.Use(middleware.AuthMiddleware())
	activity.Use(middleware.OrganisationMiddleware())
	{
		activity.GET("", activityHandler.GetActivityFeed)
	}

	// Agent task routes (agent only)
	agentTasks := router.Group("/agent")
	agentTasks.Use(middleware.AgentAuthMiddleware())
	agentTasks.Use(middleware.AgentOrganisationMiddleware())
	agentTasks.Use(middleware.RateLimitByAPIKey(60)) // 60 requests per minute per agent
	{
		agentTasks.POST("/tasks", agentTaskHandler.CreateTask)
		agentTasks.GET("/tasks", agentTaskHandler.ListTasks)
		agentTasks.GET("/tasks/:id", agentTaskHandler.GetTask)
		agentTasks.PATCH("/tasks/:id/status", agentTaskHandler.UpdateStatus)
		agentTasks.GET("/tasks/:id/comments", commentHandler.AgentGetComments)
		agentTasks.POST("/tasks/:id/comments", commentHandler.AgentCreateComment)
		agentTasks.GET("/statistics", agentHandler.GetMyStatistics)
		// Agent subtask routes
		agentTasks.GET("/tasks/:task_id/subtasks", subtaskHandler.AgentListSubtasks)
		agentTasks.POST("/tasks/:task_id/subtasks", subtaskHandler.AgentCreateSubtask)
		agentTasks.PATCH("/subtasks/:id", subtaskHandler.AgentUpdateSubtask)
		agentTasks.DELETE("/subtasks/:id", subtaskHandler.AgentDeleteSubtask)

		// Project routes (read-only for agents)
		agentTasks.GET("/projects", projectHandler.ListProjectsForAgent)
		agentTasks.GET("/projects/:id", projectHandler.GetProjectForAgent)
	}

	// OpenClaw tool endpoints (agent only)
	tools := router.Group("/tools")
	tools.Use(middleware.AgentAuthMiddleware())
	tools.Use(middleware.AgentOrganisationMiddleware())
	tools.Use(middleware.RateLimitByAPIKey(60)) // 60 requests per minute per agent
	{
		tools.POST("/tasks/create", toolsHandler.CreateTask)
		tools.POST("/tasks/update", toolsHandler.UpdateTask)
		tools.POST("/tasks/list", toolsHandler.ListTasks)
		tools.GET("/tasks/status/:id", toolsHandler.GetStatus)
	}

	// Supervisor endpoints (supervisor/admin agents only)
	supervisor := router.Group("/supervisor")
	supervisor.Use(middleware.AgentAuthMiddleware())
	supervisor.Use(middleware.AgentOrganisationMiddleware())
	supervisor.Use(middleware.RequireSupervisor())
	{
		// Agent management
		supervisor.POST("/agents", supervisorHandler.CreateAgent)
		supervisor.GET("/agents", supervisorHandler.ListAgents)
		supervisor.GET("/agents/activity", supervisorHandler.GetAgentsWithTasks)
		supervisor.PATCH("/agents/:id", supervisorHandler.UpdateAgent)
		supervisor.DELETE("/agents/:id", supervisorHandler.DeleteAgent)

		// Task management (any task)
		supervisor.GET("/tasks", supervisorHandler.ListTasks)
		supervisor.PATCH("/tasks/:id/status", supervisorHandler.UpdateTaskStatus)
		supervisor.PATCH("/tasks/:id/assign", supervisorHandler.AssignTask)
	}

	// Start server
	port := "8080"
	log.Printf("Server starting on port %s...", port)
	log.Printf("Swagger documentation available at http://localhost:%s/docs/index.html", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

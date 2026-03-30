package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/formatho/agent-todo/db"
	_ "github.com/formatho/agent-todo/docs"
	"github.com/formatho/agent-todo/handlers"
	"github.com/formatho/agent-todo/middleware"
	"github.com/formatho/agent-todo/models"
	"github.com/formatho/agent-todo/services"
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
	router.Use(metricsMiddleware.Middleware())
	router.Use(middleware.UptimeMiddleware())
	router.Use(middleware.ErrorTrackingMiddleware(metricsService))
	// router.Use(middleware.RateLimitMiddleware(100)) // 100 requests per minute per IP - TEMPORARILY DISABLED

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
	analyticsHandler := handlers.NewAnalyticsHandler()
	subscriptionHandler := handlers.NewSubscriptionHandler()
	emailHandler := handlers.NewEmailHandler()
	feedbackHandler := handlers.NewFeedbackHandler()

	// Initialize metrics service for monitoring
	metricsService := services.NewMetricsService()
	metricsMiddleware := middleware.NewMetricsMiddleware(metricsService)

	// Initialize state sync service and handler for cloud synchronization
	stateService := services.NewStateSerializationService()
	if err := stateService.EnsureDatabaseTables(); err != nil {
		log.Printf("Warning: Failed to ensure database tables for state sync: %v", err)
	}
	stateSyncHandler := handlers.NewStateSyncHandler(stateService)

	// Initialize email sequence service and seed default sequences
	emailSequenceService := services.NewEmailSequenceService(db.DB)
	if err := emailSequenceService.EnsureSequenceTables(); err != nil {
		log.Printf("Warning: Failed to ensure email sequence tables: %v", err)
	}
	if err := emailSequenceService.SeedTrialConversionSequence(); err != nil {
		log.Printf("Warning: Failed to seed trial conversion sequence: %v", err)
	}

	// Auto-migrate subscription models
	if err := db.DB.AutoMigrate(&models.Subscription{}, &models.EmailTemplate{}, &models.EmailSequence{}, &models.EmailSequenceStep{}, &models.EmailQueue{}, &models.EmailLog{}, &models.BetaFeedback{}); err != nil {
		log.Printf("Warning: Failed to migrate email/subscription models: %v", err)
	}

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

	// Analytics endpoints (public for tracking, auth required for viewing)
	analytics := router.Group("/analytics")
	analytics.Use(middleware.ProductHuntMiddleware())
	{
		analytics.POST("/track", analyticsHandler.TrackEvent)
		analytics.POST("/product-hunt-event", analyticsHandler.TrackProductHuntEvent)
	}
	// Protected analytics endpoints
	analyticsProtected := router.Group("/analytics")
	analyticsProtected.Use(middleware.AuthMiddleware())
	{
		analyticsProtected.GET("/funnel", analyticsHandler.GetFunnelStats)
		analyticsProtected.GET("/events", analyticsHandler.GetRecentEvents)
		analyticsProtected.GET("/product-hunt-metrics", analyticsHandler.GetProductHuntMetrics)
		// Task analytics endpoints
		analyticsProtected.GET("/tasks/overview", analyticsHandler.GetTaskOverview)
		analyticsProtected.GET("/tasks/agents", analyticsHandler.GetAgentMetrics)
		analyticsProtected.GET("/tasks/timeline", analyticsHandler.GetTimelineMetrics)
	}

	// Feedback endpoints (public for submission, auth required for viewing)
	feedbackPublic := router.Group("/feedback")
	{
		feedbackPublic.POST("", feedbackHandler.SubmitFeedback)
	}
	// Protected feedback endpoints (admin only)
	feedbackProtected := router.Group("/feedback")
	feedbackProtected.Use(middleware.AuthMiddleware())
	{
		feedbackProtected.GET("", feedbackHandler.ListFeedback)
		feedbackProtected.GET("/stats", feedbackHandler.GetFeedbackStats)
		feedbackProtected.GET("/:id", feedbackHandler.GetFeedback)
		feedbackProtected.PATCH("/:id/status", feedbackHandler.UpdateFeedbackStatus)
		feedbackProtected.PATCH("/:id/notes", feedbackHandler.UpdateFeedbackNotes)
	}

	// Monitoring endpoints
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
			"version": "1.0.0",
		})
	})
	
	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(metricsService.GetMetricsHandler()))

	// Swagger documentation
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Auth routes (stricter rate limiting)
	auth := router.Group("/auth")
	// auth.Use(middleware.RateLimitMiddleware(20)) // 20 requests per minute for auth - TEMPORARILY DISABLED
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
		tasks.GET("/:id/subtasks", subtaskHandler.ListSubtasks)
		tasks.POST("/:id/subtasks", subtaskHandler.CreateSubtask)
		tasks.POST("/:id/subtasks/reorder", subtaskHandler.ReorderSubtasks)
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

	// ============================================================================
	// API-PREFIXED ROUTES (for external integrations and documentation compatibility)
	// ============================================================================

	// API Auth routes - mirrors /auth
	apiAuth := router.Group("/api/auth")
	{
		apiAuth.POST("/register", authHandler.Register)
		apiAuth.POST("/login", authHandler.Login)
		apiAuth.GET("/me", middleware.AuthMiddleware(), authHandler.GetCurrentUser)
		apiAuth.POST("/switch-organisation", middleware.AuthMiddleware(), authHandler.SwitchOrganisation)
	}

	// API Organisation routes - mirrors /organisations
	apiOrganisations := router.Group("/api/organisations")
	apiOrganisations.Use(middleware.AuthMiddleware())
	{
		apiOrganisations.POST("", organisationHandler.CreateOrganisation)
		apiOrganisations.GET("", organisationHandler.ListOrganisations)
		apiOrganisations.GET("/:id", organisationHandler.GetOrganisation)
		apiOrganisations.PATCH("/:id", organisationHandler.UpdateOrganisation)
		apiOrganisations.DELETE("/:id", organisationHandler.DeleteOrganisation)
		apiOrganisations.POST("/:id/members", organisationHandler.AddOrganisationMember)
		apiOrganisations.PATCH("/:id/members/:member_id", organisationHandler.UpdateMemberRole)
		apiOrganisations.DELETE("/:id/members/:member_id", organisationHandler.RemoveOrganisationMember)
		apiOrganisations.POST("/:id/leave", organisationHandler.LeaveOrganisation)
	}

	// API Project routes - mirrors /projects
	apiProjects := router.Group("/api/projects")
	apiProjects.Use(middleware.AuthMiddleware())
	apiProjects.Use(middleware.OrganisationMiddleware())
	{
		apiProjects.POST("", projectHandler.CreateProject)
		apiProjects.GET("", projectHandler.ListProjects)
		apiProjects.GET("/:id", projectHandler.GetProject)
		apiProjects.PATCH("/:id", projectHandler.UpdateProject)
		apiProjects.DELETE("/:id", projectHandler.DeleteProject)
		apiProjects.GET("/:id/tasks", projectHandler.GetProjectTasks)
	}

	// API Agent management routes - mirrors /agents
	apiAgents := router.Group("/api/agents")
	apiAgents.Use(middleware.AuthMiddleware())
	apiAgents.Use(middleware.OrganisationMiddleware())
	{
		apiAgents.POST("", agentHandler.CreateAgent)
		apiAgents.GET("", agentHandler.ListAgents)
		apiAgents.GET("/activity", agentHandler.GetAgentsWithTasks)
		apiAgents.GET("/:id", agentHandler.GetAgent)
		apiAgents.GET("/:id/statistics", agentHandler.GetAgentStatistics)
		apiAgents.PATCH("/:id", agentHandler.UpdateAgent)
		apiAgents.DELETE("/:id", agentHandler.DeleteAgent)
	}

	// API Task routes (human) - mirrors /tasks
	apiTasks := router.Group("/api/tasks")
	apiTasks.Use(middleware.AuthMiddleware())
	apiTasks.Use(middleware.OrganisationMiddleware())
	{
		apiTasks.POST("", taskHandler.CreateTask)
		apiTasks.GET("", taskHandler.ListTasks)
		apiTasks.GET("/:id", taskHandler.GetTask)
		apiTasks.PATCH("/:id", taskHandler.UpdateTask)
		apiTasks.DELETE("/:id", taskHandler.DeleteTask)
		apiTasks.PATCH("/:id/assign", taskHandler.AssignAgent)
		apiTasks.PATCH("/:id/unassign", taskHandler.UnassignAgent)
		apiTasks.GET("/:id/comments", commentHandler.GetComments)
		apiTasks.POST("/:id/comments", commentHandler.CreateComment)
		// Subtask routes
		apiTasks.GET("/:id/subtasks", subtaskHandler.ListSubtasks)
		apiTasks.POST("/:id/subtasks", subtaskHandler.CreateSubtask)
		apiTasks.POST("/:id/subtasks/reorder", subtaskHandler.ReorderSubtasks)
	}

	// API Subtask routes - mirrors /subtasks
	apiSubtasks := router.Group("/api/subtasks")
	apiSubtasks.Use(middleware.AuthMiddleware())
	apiSubtasks.Use(middleware.OrganisationMiddleware())
	{
		apiSubtasks.GET("/:id", subtaskHandler.GetSubtask)
		apiSubtasks.PATCH("/:id", subtaskHandler.UpdateSubtask)
		apiSubtasks.DELETE("/:id", subtaskHandler.DeleteSubtask)
	}

	// API Activity feed routes - mirrors /activity
	apiActivity := router.Group("/api/activity")
	apiActivity.Use(middleware.AuthMiddleware())
	apiActivity.Use(middleware.OrganisationMiddleware())
	{
		apiActivity.GET("", activityHandler.GetActivityFeed)
	}

	// API Analytics routes - mirrors /analytics
	apiAnalytics := router.Group("/api/analytics")
	apiAnalytics.Use(middleware.AuthMiddleware())
	{
		apiAnalytics.GET("/funnel", analyticsHandler.GetFunnelStats)
		apiAnalytics.GET("/events", analyticsHandler.GetRecentEvents)
		apiAnalytics.GET("/product-hunt-metrics", analyticsHandler.GetProductHuntMetrics)
		// Task analytics endpoints
		apiAnalytics.GET("/tasks/overview", analyticsHandler.GetTaskOverview)
		apiAnalytics.GET("/tasks/agents", analyticsHandler.GetAgentMetrics)
		apiAnalytics.GET("/tasks/timeline", analyticsHandler.GetTimelineMetrics)
	}

	// API Agent task routes (agent only) - mirrors /agent
	apiAgentTasks := router.Group("/api/agent")
	apiAgentTasks.Use(middleware.AgentAuthMiddleware())
	apiAgentTasks.Use(middleware.AgentOrganisationMiddleware())
	{
		apiAgentTasks.POST("/tasks", agentTaskHandler.CreateTask)
		apiAgentTasks.GET("/tasks", agentTaskHandler.ListTasks)
		apiAgentTasks.GET("/tasks/:id", agentTaskHandler.GetTask)
		apiAgentTasks.PATCH("/tasks/:id/status", agentTaskHandler.UpdateStatus)
		apiAgentTasks.GET("/tasks/:id/comments", commentHandler.AgentGetComments)
		apiAgentTasks.POST("/tasks/:id/comments", commentHandler.AgentCreateComment)
		apiAgentTasks.GET("/statistics", agentHandler.GetMyStatistics)
		// Agent subtask routes
		apiAgentTasks.GET("/tasks/:id/subtasks", subtaskHandler.AgentListSubtasks)
		apiAgentTasks.POST("/tasks/:id/subtasks", subtaskHandler.AgentCreateSubtask)
		apiAgentTasks.PATCH("/subtasks/:id", subtaskHandler.AgentUpdateSubtask)
		apiAgentTasks.DELETE("/subtasks/:id", subtaskHandler.AgentDeleteSubtask)
		// Project routes (read-only for agents)
		apiAgentTasks.GET("/projects", projectHandler.ListProjectsForAgent)
		apiAgentTasks.GET("/projects/:id", projectHandler.GetProjectForAgent)
	}

	// ============================================================================
	// ORIGINAL ROUTES (non-prefixed)
	// ============================================================================

	// Agent task routes (agent only)
	agentTasks := router.Group("/agent")
	agentTasks.Use(middleware.AgentAuthMiddleware())
	agentTasks.Use(middleware.AgentOrganisationMiddleware())
	// agentTasks.Use(middleware.RateLimitByAPIKey(60)) // 60 requests per minute per agent - TEMPORARILY DISABLED
	{
		agentTasks.POST("/tasks", agentTaskHandler.CreateTask)
		agentTasks.GET("/tasks", agentTaskHandler.ListTasks)
		agentTasks.GET("/tasks/:id", agentTaskHandler.GetTask)
		agentTasks.PATCH("/tasks/:id/status", agentTaskHandler.UpdateStatus)
		agentTasks.GET("/tasks/:id/comments", commentHandler.AgentGetComments)
		agentTasks.POST("/tasks/:id/comments", commentHandler.AgentCreateComment)
		agentTasks.GET("/statistics", agentHandler.GetMyStatistics)
		// Agent subtask routes
		agentTasks.GET("/tasks/:id/subtasks", subtaskHandler.AgentListSubtasks)
		agentTasks.POST("/tasks/:id/subtasks", subtaskHandler.AgentCreateSubtask)
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
	// tools.Use(middleware.RateLimitByAPIKey(60)) // 60 requests per minute per agent - TEMPORARILY DISABLED
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

	// State Synchronization endpoints (for cloud sync and agent persistence)
	stateSync := router.Group("/state/sync")
	stateSync.Use(middleware.AgentAuthMiddleware())
	{
		// Agent state snapshots for cloud synchronization
		stateSync.POST("/agents/:agent_id/snapshots", stateSyncHandler.SaveAgentSnapshot)
		stateSync.GET("/agents/:agent_id/snapshot/current", stateSyncHandler.GetCurrentSnapshot)
		stateSync.GET("/agents/:agent_id/snapshots", stateSyncHandler.GetSnapshotHistory)
		stateSync.POST("/agents/:agent_id/snapshots/:snapshot_id/restore", stateSyncHandler.RestoreAgentState)

		// Task execution history for analytics and export
		stateSync.POST("/agents/:agent_id/executions", stateSyncHandler.StoreTaskExecution)
		stateSync.GET("/agents/:agent_id/executions", stateSyncHandler.GetTaskExecutionHistory)
		stateSync.GET("/agents/:agent_id/metrics/task-completion", stateSyncHandler.GetTaskCompletionMetrics)

		// Structured response storage for export functionality
		stateSync.POST("/agents/:agent_id/responses", stateSyncHandler.StoreStructuredResponse)

		// Task history export (CSV/JSON format)
		stateSync.GET("/agents/:agent_id/export", stateSyncHandler.ExportTaskHistory)

		// Team collaboration endpoints (using agent_id as organisation reference for simplicity)
		stateSync.POST("/organisations/:agent_id/members", stateSyncHandler.AddTeamMember)
		stateSync.GET("/organisations/:agent_id/members", stateSyncHandler.GetTeamMembers)
		stateSync.PATCH("/organisations/:agent_id/members/:user_id/status", stateSyncHandler.UpdateMemberStatus)
	}

	// Subscription routes
	subscriptions := router.Group("/subscriptions")
	subscriptions.Use(middleware.AuthMiddleware())
	{
		subscriptions.POST("/trial", subscriptionHandler.StartTrial)
		subscriptions.GET("/:organisation_id", subscriptionHandler.GetSubscription)
		subscriptions.POST("/:organisation_id/cancel", subscriptionHandler.CancelTrial)
	}

	// Email routes (for cron and admin)
	emails := router.Group("/emails")
	{
		// Public endpoint for cron jobs (should be protected by API key in production)
		emails.POST("/process", emailHandler.ProcessEmailQueue)
	}
	emailsProtected := router.Group("/emails")
	emailsProtected.Use(middleware.AuthMiddleware())
	{
		emailsProtected.GET("/queue", emailHandler.GetEmailQueue)
		emailsProtected.GET("/logs", emailHandler.GetEmailLogs)
		emailsProtected.GET("/sequences", emailHandler.GetEmailSequences)
		emailsProtected.POST("/seed", emailHandler.SeedSequences)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s...", port)
	log.Printf("Swagger documentation available at http://localhost:%s/docs/index.html", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

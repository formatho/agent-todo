// Package db provides database connection and seeding functionality
package db

import (
	"fmt"
	"log"

	"github.com/formatho/agent-todo/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect establishes connection to PostgreSQL database
func Connect(databaseURL string) error {
	var err error

	DB, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate the schema
	err = DB.AutoMigrate(
		&models.User{},
		&models.Agent{},
		&models.Organisation{},
		&models.OrganisationMember{},
		&models.Project{},
		&models.Task{},
		&models.TaskEvent{},
		&models.TaskComment{},
		&models.Subtask{},
	)

	if err != nil {
		return fmt.Errorf("failed to run auto migrations: %w", err)
	}

	// Run manual migrations to fix foreign key constraints
	if err := runManualMigrations(); err != nil {
		log.Printf("Warning: failed to run manual migrations: %v", err)
	}

	// Seed initial data
	if err := SeedData(); err != nil {
		log.Printf("Warning: failed to seed data: %v", err)
	}

	log.Println("Database connected and migrated successfully")
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// SeedData seeds the database with initial data
func SeedData() error {
	// Check if admin user already exists
	var count int64
	DB.Model(&models.User{}).Where("email = ?", "admin@example.com").Count(&count)
	if count > 0 {
		log.Println("Seed data already exists, skipping...")
		return nil
	}

	// Hash password for admin user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Create admin user
	adminUser := &models.User{
		Email:        "admin@example.com",
		PasswordHash: string(hashedPassword),
	}

	if err := DB.Create(adminUser).Error; err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	// Create example agent
	exampleAgent := &models.Agent{
		Name:        "Example Agent",
		APIKey:      "sk_agent_example_key_12345",
		Description: "An example agent for testing",
		Role:        models.AgentRoleRegular,
		Enabled:     true,
	}

	if err := DB.Create(exampleAgent).Error; err != nil {
		return fmt.Errorf("failed to create example agent: %w", err)
	}

	// Create example projects
	projects := []models.Project{
		{
			Name:            "Website Redesign",
			Description:     "Redesign the company website with new branding",
			Status:          models.ProjectStatusActive,
			CreatedByUserID: adminUser.ID,
		},
		{
			Name:            "API Development",
			Description:     "Build RESTful API for mobile applications",
			Status:          models.ProjectStatusActive,
			CreatedByUserID: adminUser.ID,
		},
	}

	for _, project := range projects {
		if err := DB.Create(&project).Error; err != nil {
			return fmt.Errorf("failed to create project: %w", err)
		}
	}

	// Get the first project for task assignment
	var websiteProject models.Project
	if err := DB.Where("name = ?", "Website Redesign").First(&websiteProject).Error; err != nil {
		return fmt.Errorf("failed to find website project: %w", err)
	}

	var apiProject models.Project
	if err := DB.Where("name = ?", "API Development").First(&apiProject).Error; err != nil {
		return fmt.Errorf("failed to find api project: %w", err)
	}

	// Create example tasks
	tasks := []models.Task{
		{
			Title:           "Setup project environment",
			Description:     "Configure development environment with all necessary tools",
			Status:          models.TaskStatusPending,
			Priority:        models.TaskPriorityHigh,
			ProjectID:       &websiteProject.ID,
			CreatedByUserID: &adminUser.ID,
			AssignedAgentID: &exampleAgent.ID,
		},
		{
			Title:           "Create initial documentation",
			Description:     "Write README and API documentation",
			Status:          models.TaskStatusPending,
			Priority:        models.TaskPriorityMedium,
			ProjectID:       &websiteProject.ID,
			CreatedByUserID: &adminUser.ID,
		},
		{
			Title:           "Design database schema",
			Description:     "Design database schema for user management",
			Status:          models.TaskStatusInProgress,
			Priority:        models.TaskPriorityHigh,
			ProjectID:       &apiProject.ID,
			CreatedByUserID: &adminUser.ID,
			AssignedAgentID: &exampleAgent.ID,
		},
	}

	for _, task := range tasks {
		if err := DB.Create(&task).Error; err != nil {
			return fmt.Errorf("failed to create task: %w", err)
		}

		// Create creation event
		event := &models.TaskEvent{
			TaskID:    task.ID,
			EventType: models.TaskEventCreated,
			ChangedBy: "admin@example.com",
		}
		DB.Create(event)
	}

	log.Println("Seed data created successfully")
	return nil
}

// runManualMigrations runs SQL migrations that can't be handled by AutoMigrate
func runManualMigrations() error {
	// Drop old foreign key constraint if it exists
	DB.Exec(`
		DO $$
		BEGIN
			IF EXISTS (
				SELECT 1 FROM pg_constraint
				WHERE conname = 'fk_tasks_created_by'
			) THEN
				ALTER TABLE tasks DROP CONSTRAINT fk_tasks_created_by;
			END IF;
		END $$;
	`)

	// Drop agent creator constraint if it exists
	DB.Exec(`
		DO $$
		BEGIN
			IF EXISTS (
				SELECT 1 FROM pg_constraint
				WHERE conname = 'fk_tasks_created_by_agent'
			) THEN
				ALTER TABLE tasks DROP CONSTRAINT fk_tasks_created_by_agent;
			END IF;
		END $$;
	`)

	// Add proper foreign key constraints with ON DELETE SET NULL
	DB.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint
				WHERE conname = 'fk_tasks_created_by_user'
			) THEN
				ALTER TABLE tasks ADD CONSTRAINT fk_tasks_created_by_user
				FOREIGN KEY (created_by_user_id)
				REFERENCES users(id)
				ON DELETE SET NULL
				ON UPDATE CASCADE;
			END IF;
		END $$;
	`)

	DB.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint
				WHERE conname = 'fk_tasks_created_by_agent'
			) THEN
				ALTER TABLE tasks ADD CONSTRAINT fk_tasks_created_by_agent
				FOREIGN KEY (created_by_agent_id)
				REFERENCES agents(id)
				ON DELETE SET NULL
				ON UPDATE CASCADE;
			END IF;
		END $$;
	`)

	// Ensure check constraint exists for mutual exclusivity
	DB.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint
				WHERE conname = 'chk_task_creator'
			) THEN
				ALTER TABLE tasks ADD CONSTRAINT chk_task_creator
				CHECK (
					(created_by_user_id IS NOT NULL AND created_by_agent_id IS NULL) OR
					(created_by_user_id IS NULL AND created_by_agent_id IS NOT NULL)
				);
			END IF;
		END $$;
	`)

	log.Println("Manual migrations completed successfully")
	return nil
}

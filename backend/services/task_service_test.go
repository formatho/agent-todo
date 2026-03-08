package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaskService_SoftDelete(t *testing.T) {
	// Note: This test requires a test database connection
	// In production, you would use a test database or mock

	// Skip if no database connection
	t.Skip("Requires database connection - run as integration test")

	taskService := NewTaskService()

	// Create a task first
	// task, err := taskService.Create("Test Task", "Description", "medium", nil, "", &userID, nil)
	// assert.NoError(t, err)
	// assert.NotNil(t, task)

	// Delete the task (soft delete)
	// err = taskService.Delete(task.ID.String())
	// assert.NoError(t, err)

	// Try to get the deleted task - should not find it
	// _, err = taskService.GetByID(task.ID.String())
	// assert.Error(t, err)
	// assert.Equal(t, gorm.ErrRecordNotFound, err)

	// Verify it's soft deleted in database (has deleted_at value)
	// This would require direct database query
}

func TestTaskService_ListExcludesSoftDeleted(t *testing.T) {
	// Note: This test requires a test database connection
	t.Skip("Requires database connection - run as integration test")

	// Create tasks
	// Delete one
	// List should not include deleted task
}

func TestTaskService_UnscopedCanRetrieveDeleted(t *testing.T) {
	// Note: This test requires a test database connection
	t.Skip("Requires database connection - run as integration test")

	// If we need to recover a soft-deleted task:
	// db.Unscoped().Where("id = ?", taskID).First(&task)
	// task.DeletedAt = gorm.DeletedAt{Time: time.Time{}, Valid: false}
	// db.Save(&task)
}

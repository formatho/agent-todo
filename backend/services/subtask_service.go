package services

import (
	"errors"
	"fmt"

	"github.com/formatho/agent-todo/db"
	"github.com/formatho/agent-todo/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubtaskService struct{}

func NewSubtaskService() *SubtaskService {
	return &SubtaskService{}
}

// Create creates a new subtask
func (s *SubtaskService) Create(taskID, title string, position *int) (*models.Subtask, error) {
	taskUUID, err := uuid.Parse(taskID)
	if err != nil {
		return nil, fmt.Errorf("invalid task ID: %w", err)
	}

	// Verify task exists
	var task models.Task
	if err := db.DB.First(&task, "id = ?", taskUUID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("error finding task: %w", err)
	}

	// Determine position
	pos := 0
	if position != nil {
		pos = *position
	} else {
		// Get the max position for this task's subtasks
		var maxPos int
		result := db.DB.Model(&models.Subtask{}).Where("task_id = ?", taskUUID).Select("COALESCE(MAX(position), -1)").Scan(&maxPos)
		if result.Error != nil {
			return nil, fmt.Errorf("error getting max position: %w", result.Error)
		}
		pos = maxPos + 1
	}

	subtask := models.Subtask{
		Title:    title,
		Status:   models.SubtaskStatusPending,
		TaskID:   taskUUID,
		Position: pos,
	}

	if err := db.DB.Create(&subtask).Error; err != nil {
		return nil, fmt.Errorf("error creating subtask: %w", err)
	}

	return &subtask, nil
}

// GetByID retrieves a subtask by ID
func (s *SubtaskService) GetByID(id string) (*models.Subtask, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid subtask ID: %w", err)
	}

	var subtask models.Subtask
	if err := db.DB.First(&subtask, "id = ?", uuid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("subtask not found")
		}
		return nil, fmt.Errorf("error finding subtask: %w", err)
	}

	return &subtask, nil
}

// Update updates a subtask
func (s *SubtaskService) Update(id, title string, status models.SubtaskStatus, position *int) (*models.Subtask, error) {
	subtask, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if title != "" {
		updates["title"] = title
	}
	if status != "" {
		updates["status"] = status
	}
	if position != nil {
		updates["position"] = *position
	}

	if len(updates) > 0 {
		if err := db.DB.Model(subtask).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("error updating subtask: %w", err)
		}
	}

	return subtask, nil
}

// Delete deletes a subtask (soft delete)
func (s *SubtaskService) Delete(id string) error {
	subtask, err := s.GetByID(id)
	if err != nil {
		return err
	}

	if err := db.DB.Delete(subtask).Error; err != nil {
		return fmt.Errorf("error deleting subtask: %w", err)
	}

	return nil
}

// ListByTask retrieves all subtasks for a task
func (s *SubtaskService) ListByTask(taskID string) ([]models.Subtask, error) {
	taskUUID, err := uuid.Parse(taskID)
	if err != nil {
		return nil, fmt.Errorf("invalid task ID: %w", err)
	}

	var subtasks []models.Subtask
	if err := db.DB.Where("task_id = ?", taskUUID).Order("position ASC").Find(&subtasks).Error; err != nil {
		return nil, fmt.Errorf("error fetching subtasks: %w", err)
	}

	return subtasks, nil
}

// Reorder updates the positions of multiple subtasks
func (s *SubtaskService) Reorder(taskID string, subtaskIDs []string) error {
	taskUUID, err := uuid.Parse(taskID)
	if err != nil {
		return fmt.Errorf("invalid task ID: %w", err)
	}

	// Use transaction to ensure atomic updates
	return db.DB.Transaction(func(tx *gorm.DB) error {
		for i, id := range subtaskIDs {
			uuid, err := uuid.Parse(id)
			if err != nil {
				return fmt.Errorf("invalid subtask ID %s: %w", id, err)
			}

			result := tx.Model(&models.Subtask{}).
				Where("id = ? AND task_id = ?", uuid, taskUUID).
				Update("position", i)

			if result.Error != nil {
				return fmt.Errorf("error updating position for subtask %s: %w", id, result.Error)
			}

			if result.RowsAffected == 0 {
				return fmt.Errorf("subtask %s not found in task %s", id, taskID)
			}
		}
		return nil
	})
}

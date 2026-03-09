package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/formatho/agent-todo/services"
	"github.com/gin-gonic/gin"
)

type ReminderHandler struct {
	taskService *services.TaskService
}

func NewReminderHandler() *ReminderHandler {
	return &ReminderHandler{
		taskService: services.NewTaskService(),
	}
}

// UpcomingTasksResponse represents the response for upcoming due tasks
type UpcomingTasksResponse struct {
	Tasks []TaskReminder `json:"tasks"`
	Count int            `json:"count"`
}

// TaskReminder represents a task with reminder info
type TaskReminder struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Priority     string `json:"priority"`
	Status       string `json:"status"`
	DueDate      string `json:"due_date"`
	TimeUntilDue string `json:"time_until_due"`
	ProjectName  string `json:"project_name,omitempty"`
	AgentName    string `json:"agent_name,omitempty"`
}

// GetUpcomingDueTasks returns tasks with due dates within a specified time window
// @Summary Get upcoming due tasks
// @Description Get tasks that are due within the specified hours (default 24)
// @Tags reminders
// @Accept json
// @Produce json
// @Param hours query int false "Hours to look ahead (default 24)" minimum(1) maximum(168)
// @Success 200 {object} UpcomingTasksResponse
// @Failure 500 {object} ErrorResponse
// @Router /reminders/upcoming [get]
func (h *ReminderHandler) GetUpcomingDueTasks(c *gin.Context) {
	// Get hours from query param, default to 24
	hoursStr := c.DefaultQuery("hours", "24")
	hours, err := strconv.Atoi(hoursStr)
	if err != nil || hours < 1 {
		hours = 24
	}
	if hours > 168 {
		hours = 168 // Max 1 week
	}

	tasks, err := h.taskService.GetUpcomingDueTasks(time.Duration(hours) * time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch upcoming tasks"})
		return
	}

	reminders := make([]TaskReminder, len(tasks))
	for i, task := range tasks {
		reminder := TaskReminder{
			ID:           task.ID.String(),
			Title:        task.Title,
			Priority:     string(task.Priority),
			Status:       string(task.Status),
			DueDate:      task.DueDate.Format(time.RFC3339),
			TimeUntilDue: formatDuration(time.Until(*task.DueDate)),
		}
		if task.Project != nil {
			reminder.ProjectName = task.Project.Name
		}
		if task.AssignedAgent != nil {
			reminder.AgentName = task.AssignedAgent.Name
		}
		reminders[i] = reminder
	}

	c.JSON(http.StatusOK, UpcomingTasksResponse{
		Tasks: reminders,
		Count: len(reminders),
	})
}

// GetOverdueTasks returns tasks that are past their due date
// @Summary Get overdue tasks
// @Description Get tasks that are past their due date and not completed
// @Tags reminders
// @Accept json
// @Produce json
// @Success 200 {object} UpcomingTasksResponse
// @Failure 500 {object} ErrorResponse
// @Router /reminders/overdue [get]
func (h *ReminderHandler) GetOverdueTasks(c *gin.Context) {
	tasks, err := h.taskService.GetOverdueTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch overdue tasks"})
		return
	}

	reminders := make([]TaskReminder, len(tasks))
	for i, task := range tasks {
		reminder := TaskReminder{
			ID:           task.ID.String(),
			Title:        task.Title,
			Priority:     string(task.Priority),
			Status:       string(task.Status),
			DueDate:      task.DueDate.Format(time.RFC3339),
			TimeUntilDue: formatDuration(time.Until(*task.DueDate)) + " overdue",
		}
		if task.Project != nil {
			reminder.ProjectName = task.Project.Name
		}
		if task.AssignedAgent != nil {
			reminder.AgentName = task.AssignedAgent.Name
		}
		reminders[i] = reminder
	}

	c.JSON(http.StatusOK, UpcomingTasksResponse{
		Tasks: reminders,
		Count: len(reminders),
	})
}

// formatDuration formats a duration in a human-readable way
func formatDuration(d time.Duration) string {
	if d < 0 {
		d = -d
	}

	hours := int(d.Hours())
	if hours < 1 {
		minutes := int(d.Minutes())
		return strconv.Itoa(minutes) + " minutes"
	}
	if hours < 24 {
		return strconv.Itoa(hours) + " hours"
	}
	days := hours / 24
	remainingHours := hours % 24
	if remainingHours == 0 {
		return strconv.Itoa(days) + " days"
	}
	return strconv.Itoa(days) + " days " + strconv.Itoa(remainingHours) + " hours"
}

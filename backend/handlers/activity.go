package handlers

import (
	"net/http"
	"strconv"

	"github.com/formatho/agent-todo/services"
	"github.com/gin-gonic/gin"
)

type ActivityHandler struct {
	activityService *services.ActivityService
}

func NewActivityHandler() *ActivityHandler {
	return &ActivityHandler{
		activityService: services.NewActivityService(),
	}
}

// GetActivityFeed godoc
// @Summary Get activity feed
// @Description Get recent activity across all tasks
// @Tags activity
// @Produce json
// @Security Bearer
// @Param limit query int false "Number of events to return" default(50)
// @Success 200 {array} services.ActivityEvent
// @Failure 401 {object} map[string]string
// @Router /activity [get]
func (h *ActivityHandler) GetActivityFeed(c *gin.Context) {
	limit := 50
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	events, err := h.activityService.GetRecentActivity(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

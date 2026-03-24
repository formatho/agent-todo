package handlers

import (
	"net/http"
	"strconv"

	"github.com/formatho/agent-todo/middleware"
	"github.com/formatho/agent-todo/models"
	"github.com/formatho/agent-todo/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AnalyticsHandler struct {
	analyticsService *services.AnalyticsService
}

func NewAnalyticsHandler() *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: services.NewAnalyticsService(),
	}
}

// TrackEventRequest represents the request body for tracking an event
type TrackEventRequest struct {
	EventType models.AnalyticsEventType `json:"event_type" binding:"required"`
	Page      string                    `json:"page" binding:"required"`
	Plan      string                    `json:"plan"`
	SessionID string                    `json:"session_id"`
	Metadata  map[string]interface{}    `json:"metadata"`
}

// TrackEvent godoc
// @Summary Track an analytics event
// @Description Track an analytics event for conversion funnel tracking
// @Tags analytics
// @Accept json
// @Produce json
// @Param request body TrackEventRequest true "Event data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /analytics/track [post]
func (h *AnalyticsHandler) TrackEvent(c *gin.Context) {
	var req TrackEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate event type
	validTypes := map[models.AnalyticsEventType]bool{
		models.AnalyticsEventPageView:        true,
		models.AnalyticsEventCheckoutStart:   true,
		models.AnalyticsEventCheckoutComplete: true,
	}
	if !validTypes[req.EventType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event type"})
		return
	}

	// Get user ID if authenticated (optional)
	var userUUID *uuid.UUID
	if user, exists := c.Get(middleware.UserContextKey); exists {
		if u, ok := user.(*models.User); ok {
			userUUID = &u.ID
		}
	}

	// Get request metadata
	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()
	referrer := c.GetHeader("Referer")

	err := h.analyticsService.TrackEvent(&services.TrackEventRequest{
		EventType: req.EventType,
		Page:      req.Page,
		Plan:      req.Plan,
		SessionID: req.SessionID,
		Metadata:  req.Metadata,
	}, userUUID, userAgent, ipAddress, referrer)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to track event"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "ok"})
}

// GetFunnelStats godoc
// @Summary Get conversion funnel statistics
// @Description Get analytics stats for the conversion funnel
// @Tags analytics
// @Produce json
// @Param days query int false "Number of days to look back (default: 30)"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /analytics/funnel [get]
func (h *AnalyticsHandler) GetFunnelStats(c *gin.Context) {
	days := 30
	if daysStr := c.Query("days"); daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}

	stats, err := h.analyticsService.GetFunnelStats(days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetRecentEvents godoc
// @Summary Get recent analytics events
// @Description Get recent analytics events
// @Tags analytics
// @Produce json
// @Param limit query int false "Number of events to return (default: 100)"
// @Param event_type query string false "Filter by event type"
// @Success 200 {array} models.AnalyticsEvent
// @Failure 500 {object} map[string]string
// @Router /analytics/events [get]
func (h *AnalyticsHandler) GetRecentEvents(c *gin.Context) {
	limit := 100
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 1000 {
			limit = l
		}
	}

	var eventType *models.AnalyticsEventType
	if et := c.Query("event_type"); et != "" {
		t := models.AnalyticsEventType(et)
		eventType = &t
	}

	events, err := h.analyticsService.GetRecentEvents(limit, eventType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get events"})
		return
	}

	c.JSON(http.StatusOK, events)
}

package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

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

// TrackProductHuntEvent godoc
// @Summary Track a Product Hunt specific event
// @Description Track events related to Product Hunt launch (referrals, upvotes, conversions)
// @Tags analytics
// @Accept json
// @Produce json
// @Param request body TrackProductHuntEventRequest true "Product Hunt event data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /analytics/product-hunt-event [post]
func (h *AnalyticsHandler) TrackProductHuntEvent(c *gin.Context) {
	type TrackProductHuntEventRequest struct {
		EventType  string                 `json:"event_type" binding:"required"`
		Source     string                 `json:"source"`
		VoteType   string                 `json:"vote_type"`
		Conversion string                 `json:"conversion_type"`
		Metadata   map[string]interface{} `json:"metadata"`
	}

	var req TrackProductHuntEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get Product Hunt specific metadata from headers
	userAgent := c.GetHeader("User-Agent")
	referrer := c.GetHeader("Referer")
	sourceType := "unknown"
	
	// Determine source type
	if strings.Contains(strings.ToLower(referrer), "producthunt.com") {
		sourceType = "producthunt_page"
	} else if strings.Contains(strings.ToLower(userAgent), "producthunt") {
		sourceType = "producthunt_bot"
	} else if strings.Contains(strings.ToLower(referrer), "twitter.com") && strings.Contains(strings.ToLower(referrer), "producthunt") {
		sourceType = "twitter_referral"
	} else if strings.Contains(strings.ToLower(referrer), "linkedin.com") && strings.Contains(strings.ToLower(referrer), "producthunt") {
		sourceType = "linkedin_referral"
	} else {
		sourceType = "other_referral"
	}

	// Track using metrics service
	metricsService := services.NewMetricsService()
	
	switch req.EventType {
	case "referral":
		metricsService.RecordPHReferral(sourceType)
	case "upvote":
		metricsService.RecordPHUpvote()
	case "conversion":
		conversionType := req.Conversion
		if conversionType == "" {
			conversionType = "unknown"
		}
		metricsService.RecordPHConversion(conversionType)
	}

	// Also store in database for historical analysis
	err := h.analyticsService.TrackProductHuntEvent(&services.TrackProductHuntEventRequest{
		EventType:  req.EventType,
		Source:     req.Source,
		VoteType:   req.VoteType,
		Conversion: req.Conversion,
		Metadata:  req.Metadata,
	}, userAgent, referrer, sourceType)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to track Product Hunt event"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "ok", "source_type": sourceType})
}

// GetProductHuntMetrics godoc
// @Summary Get Product Hunt launch metrics
// @Description Get aggregated Product Hunt launch performance metrics
// @Tags analytics
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /analytics/product-hunt-metrics [get]
func (h *AnalyticsHandler) GetProductHuntMetrics(c *gin.Context) {
	metricsService := services.NewMetricsService()
	
	// Get metrics from metrics service
	phMetrics := metricsService.GetProductHuntMetrics()
	
	// Get historical data from database
	historicalMetrics, err := h.analyticsService.GetProductHuntHistoricalMetrics(24) // Last 24 hours
	if err != nil {
		// Continue even if historical metrics fail
		historicalMetrics = make(map[string]interface{})
	}

	response := gin.H{
		"current_metrics": phMetrics,
		"historical_metrics": historicalMetrics,
		"timestamp": time.Now().Format(time.RFC3339),
		"launch_status": "active",
	}

	c.JSON(http.StatusOK, response)
}

// GetTaskOverview godoc
// @Summary Get task overview statistics
// @Description Get overall task statistics including counts by status, priority, and completion metrics
// @Tags analytics
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /analytics/tasks/overview [get]
func (h *AnalyticsHandler) GetTaskOverview(c *gin.Context) {
	// Get organisation ID from context if available
	var organisationID *uuid.UUID
	if orgID, exists := c.Get("organisation_id"); exists {
		if id, ok := orgID.(uuid.UUID); ok {
			organisationID = &id
		}
	}

	stats, err := h.analyticsService.GetTaskOverview(organisationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get task overview"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetAgentMetrics godoc
// @Summary Get agent-specific task metrics
// @Description Get task metrics broken down by agent including workload, completion rates, and performance
// @Tags analytics
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /analytics/tasks/agents [get]
func (h *AnalyticsHandler) GetAgentMetrics(c *gin.Context) {
	// Get organisation ID from context if available
	var organisationID *uuid.UUID
	if orgID, exists := c.Get("organisation_id"); exists {
		if id, ok := orgID.(uuid.UUID); ok {
			organisationID = &id
		}
	}

	stats, err := h.analyticsService.GetAgentMetrics(organisationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get agent metrics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetTimelineMetrics godoc
// @Summary Get time-based task metrics
// @Description Get task creation and completion trends over time
// @Tags analytics
// @Produce json
// @Param days query int false "Number of days to look back (default: 30)"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /analytics/tasks/timeline [get]
func (h *AnalyticsHandler) GetTimelineMetrics(c *gin.Context) {
	days := 30
	if daysStr := c.Query("days"); daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}

	// Get organisation ID from context if available
	var organisationID *uuid.UUID
	if orgID, exists := c.Get("organisation_id"); exists {
		if id, ok := orgID.(uuid.UUID); ok {
			organisationID = &id
		}
	}

	stats, err := h.analyticsService.GetTimelineMetrics(days, organisationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get timeline metrics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

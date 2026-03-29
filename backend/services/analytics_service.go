package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/formatho/agent-todo/db"
	"github.com/formatho/agent-todo/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnalyticsService struct {
	db *gorm.DB
}

func NewAnalyticsService() *AnalyticsService {
	return &AnalyticsService{
		db: db.GetDB(),
	}
}

// TrackEventRequest represents the request to track an analytics event
type TrackEventRequest struct {
	EventType models.AnalyticsEventType `json:"event_type" binding:"required"`
	Page      string                    `json:"page" binding:"required"`
	Plan      string                    `json:"plan"`
	SessionID string                    `json:"session_id"`
	Metadata  map[string]interface{}    `json:"metadata"`
}

// TrackEvent records a new analytics event
func (s *AnalyticsService) TrackEvent(req *TrackEventRequest, userID *uuid.UUID, userAgent, ipAddress, referrer string) error {
	var metadataStr string
	if req.Metadata != nil {
		bytes, err := json.Marshal(req.Metadata)
		if err != nil {
			metadataStr = "{}"
		} else {
			metadataStr = string(bytes)
		}
	}

	event := &models.AnalyticsEvent{
		EventType: req.EventType,
		Page:      req.Page,
		Plan:      req.Plan,
		UserID:    userID,
		SessionID: req.SessionID,
		Metadata:  metadataStr,
		UserAgent: userAgent,
		IPAddress: ipAddress,
		Referrer:  referrer,
	}

	return s.db.Create(event).Error
}

// GetFunnelStats returns conversion funnel statistics
func (s *AnalyticsService) GetFunnelStats(days int) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Calculate date threshold
	dateThreshold := fmt.Sprintf("NOW() - INTERVAL '%d day'", days)

	// Count pricing page views
	var pricingViews int64
	s.db.Model(&models.AnalyticsEvent{}).
		Where("event_type = ? AND page = ? AND created_at > ?", models.AnalyticsEventPageView, "pricing", dateThreshold).
		Count(&pricingViews)
	stats["pricing_page_views"] = pricingViews

	// Count checkout starts
	var checkoutStarts int64
	s.db.Model(&models.AnalyticsEvent{}).
		Where("event_type = ? AND created_at > ?", models.AnalyticsEventCheckoutStart, dateThreshold).
		Count(&checkoutStarts)
	stats["checkout_starts"] = checkoutStarts

	// Count checkout completions
	var checkoutCompletions int64
	s.db.Model(&models.AnalyticsEvent{}).
		Where("event_type = ? AND created_at > ?", models.AnalyticsEventCheckoutComplete, dateThreshold).
		Count(&checkoutCompletions)
	stats["checkout_completions"] = checkoutCompletions

	// Calculate conversion rates
	if pricingViews > 0 {
		stats["pricing_to_checkout_rate"] = float64(checkoutStarts) / float64(pricingViews) * 100
		stats["pricing_to_completion_rate"] = float64(checkoutCompletions) / float64(pricingViews) * 100
	} else {
		stats["pricing_to_checkout_rate"] = 0.0
		stats["pricing_to_completion_rate"] = 0.0
	}

	if checkoutStarts > 0 {
		stats["checkout_to_completion_rate"] = float64(checkoutCompletions) / float64(checkoutStarts) * 100
	} else {
		stats["checkout_to_completion_rate"] = 0.0
	}

	// Get stats by plan
	type PlanStat struct {
		Plan      string `json:"plan"`
		EventType string `json:"event_type"`
		Count     int64  `json:"count"`
	}
	var planStats []PlanStat
	s.db.Model(&models.AnalyticsEvent{}).
		Select("plan, event_type, COUNT(*) as count").
		Where("plan IS NOT NULL AND plan != '' AND created_at > ?", dateThreshold).
		Group("plan, event_type").
		Scan(&planStats)
	stats["by_plan"] = planStats

	return stats, nil
}

// GetRecentEvents returns the most recent analytics events
func (s *AnalyticsService) GetRecentEvents(limit int, eventType *models.AnalyticsEventType) ([]models.AnalyticsEvent, error) {
	var events []models.AnalyticsEvent
	query := s.db.Model(&models.AnalyticsEvent{}).Order("created_at DESC").Limit(limit)

	if eventType != nil {
		query = query.Where("event_type = ?", *eventType)
	}

	err := query.Find(&events).Error
	return events, err
}

// TrackProductHuntEventRequest represents the request to track a Product Hunt specific event
type TrackProductHuntEventRequest struct {
	EventType  string                 `json:"event_type" binding:"required"`
	Source     string                 `json:"source"`
	VoteType   string                 `json:"vote_type"`
	Conversion string                 `json:"conversion_type"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// TrackProductHuntEvent records a Product Hunt specific event
func (s *AnalyticsService) TrackProductHuntEvent(req *TrackProductHuntEventRequest, userAgent, referrer, sourceType string) error {
	var metadataStr string
	if req.Metadata != nil {
		bytes, err := json.Marshal(req.Metadata)
		if err != nil {
			metadataStr = "{}"
		} else {
			metadataStr = string(bytes)
		}
	}

	// Create a custom analytics event for Product Hunt
	event := &models.AnalyticsEvent{
		EventType: models.AnalyticsEventType("product_hunt_" + req.EventType), // Cast to AnalyticsEventType
		Page:      "product_hunt",
		Plan:      req.Source,
		UserAgent: userAgent,
		IPAddress: "", // Would need to capture this from context
		Referrer:  referrer,
		Metadata:  metadataStr,
	}

	// Add additional metadata to session tracking
	if req.EventType == "referral" {
		// Create a session identifier for tracking
		sessionID := generateSessionID(userAgent, referrer)
		event.SessionID = sessionID
	}

	return s.db.Create(event).Error
}

// GetProductHuntHistoricalMetrics returns historical Product Hunt metrics for the specified time period
func (s *AnalyticsService) GetProductHuntHistoricalMetrics(hours int) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	// Calculate date threshold
	dateThreshold := fmt.Sprintf("NOW() - INTERVAL '%d hour'", hours)

	// Count Product Hunt referrals
	var phReferrals int64
	s.db.Model(&models.AnalyticsEvent{}).
		Where("event_type = ? AND created_at > ?", "product_hunt_referral", dateThreshold).
		Count(&phReferrals)
	stats["ph_referrals_24h"] = phReferrals

	// Count Product Hunt upvotes
	var phUpvotes int64
	s.db.Model(&models.AnalyticsEvent{}).
		Where("event_type = ? AND created_at > ?", "product_hunt_upvote", dateThreshold).
		Count(&phUpvotes)
	stats["ph_upvotes_24h"] = phUpvotes

	// Count Product Hunt conversions
	var phConversions int64
	s.db.Model(&models.AnalyticsEvent{}).
		Where("event_type = ? AND created_at > ?", "product_hunt_conversion", dateThreshold).
		Count(&phConversions)
	stats["ph_conversions_24h"] = phConversions

	// Calculate conversion rate
	if phReferrals > 0 {
		stats["ph_conversion_rate"] = float64(phConversions) / float64(phReferrals) * 100
	} else {
		stats["ph_conversion_rate"] = 0.0
	}

	// Get hourly breakdown for upvotes (for tracking launch peaks)
	type HourlyStat struct {
		Hour  string `json:"hour"`
		Count int64  `json:"count"`
	}
	var hourlyStats []HourlyStat
	s.db.Model(&models.AnalyticsEvent{}).
		Select("TO_CHAR(created_at, 'HH24') as hour, COUNT(*) as count").
		Where("event_type = ? AND created_at > ?", "product_hunt_upvote", dateThreshold).
		Group("hour").
		Order("hour").
		Scan(&hourlyStats)
	stats["ph_hourly_upvotes"] = hourlyStats

	// Get referral sources breakdown
	type SourceStat struct {
		Source string `json:"source"`
		Count  int64  `json:"count"`
	}
	var sourceStats []SourceStat
	s.db.Model(&models.AnalyticsEvent{}).
		Select("metadata->>'source_type' as source, COUNT(*) as count").
		Where("event_type = ? AND created_at > ?", "product_hunt_referral", dateThreshold).
		Group("source").
		Scan(&sourceStats)
	stats["ph_referral_sources"] = sourceStats

	return stats, nil
}

// generateSessionID generates a simple session ID based on user agent and referrer
func generateSessionID(userAgent, referrer string) string {
	// Simple hash-based session ID generation
	input := userAgent + referrer + time.Now().Format("20060102150405")
	// In a real implementation, you'd use a proper hash function
	return fmt.Sprintf("ph_%d", len(input))
}

// GetTaskOverview returns overall task statistics
func (s *AnalyticsService) GetTaskOverview(organisationID *uuid.UUID) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	query := s.db.Model(&models.Task{})
	if organisationID != nil {
		query = query.Where("organisation_id = ?", organisationID)
	}

	// Total tasks by status
	type StatusCount struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	var statusCounts []StatusCount
	query.Select("status, COUNT(*) as count").
		Where("deleted_at IS NULL").
		Group("status").
		Scan(&statusCounts)
	stats["by_status"] = statusCounts

	// Total tasks by priority
	type PriorityCount struct {
		Priority string `json:"priority"`
		Count    int64  `json:"count"`
	}
	var priorityCounts []PriorityCount
	query.Select("priority, COUNT(*) as count").
		Where("deleted_at IS NULL").
		Group("priority").
		Scan(&priorityCounts)
	stats["by_priority"] = priorityCounts

	// Overall totals
	var totalTasks int64
	query.Where("deleted_at IS NULL").Count(&totalTasks)
	stats["total_tasks"] = totalTasks

	var completedTasks int64
	query.Where("status = ? AND deleted_at IS NULL", models.TaskStatusCompleted).Count(&completedTasks)
	stats["completed_tasks"] = completedTasks

	if totalTasks > 0 {
		stats["completion_rate"] = float64(completedTasks) / float64(totalTasks) * 100
	} else {
		stats["completion_rate"] = 0.0
	}

	// Tasks created in last 7 days
	var recentTasks int64
	query.Where("created_at > NOW() - INTERVAL '7 days' AND deleted_at IS NULL").Count(&recentTasks)
	stats["tasks_last_7_days"] = recentTasks

	// Average completion time (for completed tasks)
	var avgCompletionTime float64
	s.db.Model(&models.Task{}).
		Select("AVG(EXTRACT(EPOCH FROM (updated_at - created_at))/3600)").
		Where("status = ? AND deleted_at IS NULL", models.TaskStatusCompleted).
		Scan(&avgCompletionTime)
	stats["avg_completion_time_hours"] = avgCompletionTime

	return stats, nil
}

// GetAgentMetrics returns agent-specific task metrics
func (s *AnalyticsService) GetAgentMetrics(organisationID *uuid.UUID) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Get all agents first
	var agents []models.Agent
	query := s.db.Model(&models.Agent{})
	if organisationID != nil {
		query = query.Where("organisation_id = ?", organisationID)
	}
	query.Find(&agents)

	// Metrics per agent
	type AgentMetric struct {
		AgentID            string  `json:"agent_id"`
		AgentName          string  `json:"agent_name"`
		TotalTasks         int64   `json:"total_tasks"`
		CompletedTasks     int64   `json:"completed_tasks"`
		InProgressTasks    int64   `json:"in_progress_tasks"`
		PendingTasks       int64   `json:"pending_tasks"`
		BlockedTasks       int64   `json:"blocked_tasks"`
		CompletionRate     float64 `json:"completion_rate"`
		AvgCompletionTime  float64 `json:"avg_completion_time_hours"`
	}

	var agentMetrics []AgentMetric

	for _, agent := range agents {
		metric := AgentMetric{
			AgentID:   agent.ID.String(),
			AgentName: agent.Name,
		}

		// Total tasks assigned to this agent
		taskQuery := s.db.Model(&models.Task{}).Where("assigned_agent_id = ? AND deleted_at IS NULL", agent.ID)
		taskQuery.Count(&metric.TotalTasks)

		// Tasks by status
		taskQuery.Where("status = ?", models.TaskStatusCompleted).Count(&metric.CompletedTasks)
		taskQuery.Where("status = ?", models.TaskStatusInProgress).Count(&metric.InProgressTasks)
		taskQuery.Where("status = ?", models.TaskStatusPending).Count(&metric.PendingTasks)
		taskQuery.Where("status = ?", models.TaskStatusBlocked).Count(&metric.BlockedTasks)

		// Completion rate
		if metric.TotalTasks > 0 {
			metric.CompletionRate = float64(metric.CompletedTasks) / float64(metric.TotalTasks) * 100
		}

		// Average completion time
		var avgTime float64
		s.db.Model(&models.Task{}).
			Select("AVG(EXTRACT(EPOCH FROM (updated_at - created_at))/3600)").
			Where("assigned_agent_id = ? AND status = ? AND deleted_at IS NULL", agent.ID, models.TaskStatusCompleted).
			Scan(&avgTime)
		metric.AvgCompletionTime = avgTime

		agentMetrics = append(agentMetrics, metric)
	}

	stats["agents"] = agentMetrics

	// Overall agent statistics
	var totalAgents int64
	s.db.Model(&models.Agent{}).Count(&totalAgents)
	stats["total_agents"] = totalAgents

	var activeAgents int64
	s.db.Model(&models.Agent{}).
		Joins("JOIN tasks ON tasks.assigned_agent_id = agents.id").
		Where("tasks.status = ? AND tasks.deleted_at IS NULL", models.TaskStatusInProgress).
		Distinct("agents.id").
		Count(&activeAgents)
	stats["active_agents"] = activeAgents

	return stats, nil
}

// GetTimelineMetrics returns time-based task metrics
func (s *AnalyticsService) GetTimelineMetrics(days int, organisationID *uuid.UUID) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	dateThreshold := fmt.Sprintf("NOW() - INTERVAL '%d day'", days)

	// Tasks created over time (by day)
	type DailyCount struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}
	var createdDaily []DailyCount
	query := s.db.Model(&models.Task{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where(fmt.Sprintf("created_at > %s AND deleted_at IS NULL", dateThreshold))
	if organisationID != nil {
		query = query.Where("organisation_id = ?", organisationID)
	}
	query.Group("date").
		Order("date").
		Scan(&createdDaily)
	stats["created_daily"] = createdDaily

	// Tasks completed over time (by day)
	var completedDaily []DailyCount
	query = s.db.Model(&models.Task{}).
		Select("DATE(updated_at) as date, COUNT(*) as count").
		Where(fmt.Sprintf("updated_at > %s AND status = ? AND deleted_at IS NULL", dateThreshold), models.TaskStatusCompleted)
	if organisationID != nil {
		query = query.Where("organisation_id = ?", organisationID)
	}
	query.Group("date").
		Order("date").
		Scan(&completedDaily)
	stats["completed_daily"] = completedDaily

	// Average task age (in hours) for non-completed tasks
	var avgTaskAge float64
	query = s.db.Model(&models.Task{}).
		Select("AVG(EXTRACT(EPOCH FROM (NOW() - created_at))/3600)").
		Where("status != ? AND deleted_at IS NULL", models.TaskStatusCompleted)
	if organisationID != nil {
		query = query.Where("organisation_id = ?", organisationID)
	}
	query.Scan(&avgTaskAge)
	stats["avg_task_age_hours"] = avgTaskAge

	// Tasks by day of week
	type WeekdayCount struct {
		Weekday string `json:"weekday"`
		Count   int64  `json:"count"`
	}
	var byWeekday []WeekdayCount
	query = s.db.Model(&models.Task{}).
		Select("TO_CHAR(created_at, 'Day') as weekday, COUNT(*) as count").
		Where(fmt.Sprintf("created_at > %s AND deleted_at IS NULL", dateThreshold))
	if organisationID != nil {
		query = query.Where("organisation_id = ?", organisationID)
	}
	query.Group("weekday").
		Order("COUNT(*) DESC").
		Scan(&byWeekday)
	stats["by_weekday"] = byWeekday

	return stats, nil
}
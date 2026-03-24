package services

import (
	"encoding/json"
	"fmt"

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
		Where("event_type = ? AND page = ? AND created_at > "+dateThreshold, models.AnalyticsEventPageView, "pricing").
		Count(&pricingViews)
	stats["pricing_page_views"] = pricingViews

	// Count checkout starts
	var checkoutStarts int64
	s.db.Model(&models.AnalyticsEvent{}).
		Where("event_type = ? AND created_at > "+dateThreshold, models.AnalyticsEventCheckoutStart).
		Count(&checkoutStarts)
	stats["checkout_starts"] = checkoutStarts

	// Count checkout completions
	var checkoutCompletions int64
	s.db.Model(&models.AnalyticsEvent{}).
		Where("event_type = ? AND created_at > "+dateThreshold, models.AnalyticsEventCheckoutComplete).
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
		Where("plan IS NOT NULL AND plan != '' AND created_at > "+dateThreshold).
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

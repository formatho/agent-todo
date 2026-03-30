package services

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/formatho/agent-todo/db"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/gorm"
)

// MetricsService handles Prometheus metrics collection
type MetricsService struct {
	db *gorm.DB

	// HTTP request metrics
	httpRequestsTotal    *prometheus.CounterVec
	httpRequestDuration   *prometheus.HistogramVec
	httpResponseSize     *prometheus.HistogramVec
	httpErrorsTotal      *prometheus.CounterVec

	// Database metrics
	dbConnectionsActive   *prometheus.GaugeVec
	dbConnectionsIdle     *prometheus.GaugeVec
	dbConnectionsWait     *prometheus.GaugeVec
	dbQueriesTotal        *prometheus.CounterVec
	dbQueryDuration       *prometheus.HistogramVec
	dbErrorsTotal         *prometheus.CounterVec

	// Application metrics
	appErrorsTotal        *prometheus.CounterVec
	appTasksCreatedTotal  *prometheus.CounterVec
	appTasksCompletedTotal *prometheus.CounterVec
	appUsersTotal         *prometheus.GaugeVec
	appAgentsActive       *prometheus.GaugeVec

	// Product Hunt specific metrics
	phReferralsTotal      *prometheus.CounterVec
	phUpvotesTotal        *prometheus.CounterVec
	phConversionsTotal    *prometheus.CounterVec
}

// NewMetricsService creates a new metrics service
func NewMetricsService() *MetricsService {
	metrics := &MetricsService{
		db: db.GetDB(),

		// HTTP request metrics
		httpRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "endpoint", "status_code"},
		),
		httpRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "endpoint"},
		),
		httpResponseSize: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_response_size_bytes",
				Help:    "HTTP response size in bytes",
				Buckets: prometheus.ExponentialBuckets(100, 10, 12),
			},
			[]string{"method", "endpoint"},
		),
		httpErrorsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_errors_total",
				Help: "Total number of HTTP errors",
			},
			[]string{"method", "endpoint", "error_type"},
		),

		// Database metrics
		dbConnectionsActive: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "db_connections_active",
				Help: "Number of active database connections",
			},
			[]string{},
		),
		dbConnectionsIdle: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "db_connections_idle",
				Help: "Number of idle database connections",
			},
			[]string{},
		),
		dbConnectionsWait: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "db_connections_wait_seconds",
				Help: "Database connection wait time in seconds",
			},
			[]string{},
		),
		dbQueriesTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "db_queries_total",
				Help: "Total number of database queries",
			},
			[]string{"query_type", "table"},
		),
		dbQueryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "db_query_duration_seconds",
				Help:    "Database query duration in seconds",
				Buckets: prometheus.ExponentialBuckets(.001, 2, 15),
			},
			[]string{"query_type", "table"},
		),
		dbErrorsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "db_errors_total",
				Help: "Total number of database errors",
			},
			[]string{"error_type", "table"},
		),

		// Application metrics
		appErrorsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "app_errors_total",
				Help: "Total number of application errors",
			},
			[]string{"error_type", "component"},
		),
		appTasksCreatedTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "app_tasks_created_total",
				Help: "Total number of tasks created",
			},
			[]string{"agent_type", "priority"},
		),
		appTasksCompletedTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "app_tasks_completed_total",
				Help: "Total number of tasks completed",
			},
			[]string{"agent_type", "priority"},
		),
		appUsersTotal: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "app_users_total",
				Help: "Total number of users",
			},
			[]string{"plan_type"},
		),
		appAgentsActive: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "app_agents_active",
				Help: "Number of active agents",
			},
			[]string{"agent_role"},
		),

		// Product Hunt specific metrics
		phReferralsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "ph_referrals_total",
				Help: "Total number of Product Hunt referrals",
			},
			[]string{"source_type"},
		),
		phUpvotesTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "ph_upvotes_total",
				Help: "Total number of Product Hunt upvotes",
			},
			[]string{"hour"},
		),
		phConversionsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "ph_conversions_total",
				Help: "Total number of conversions from Product Hunt traffic",
			},
			[]string{"conversion_type"},
		),
	}

	return metrics
}

// ObserveHTTPRequest records HTTP request metrics
func (m *MetricsService) ObserveHTTPRequest(method, endpoint string, statusCode int, duration time.Duration, responseSize int64) {
	m.httpRequestsTotal.WithLabelValues(method, endpoint, string(rune(statusCode))).Inc()
	m.httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
	m.httpResponseSize.WithLabelValues(method, endpoint).Observe(float64(responseSize))
}

// ObserveHTTPError records HTTP error metrics
func (m *MetricsService) ObserveHTTPError(method, endpoint, errorType string) {
	m.httpErrorsTotal.WithLabelValues(method, endpoint, errorType).Inc()
}

// ObserveDBQuery records database query metrics
func (m *MetricsService) ObserveDBQuery(queryType, table string, duration time.Duration, err error) {
	m.dbQueriesTotal.WithLabelValues(queryType, table).Inc()
	m.dbQueryDuration.WithLabelValues(queryType, table).Observe(duration.Seconds())

	if err != nil {
		errorType := "connection"
		if err == sql.ErrNoRows {
			errorType = "no_rows"
		}
		m.dbErrorsTotal.WithLabelValues(errorType, table).Inc()
	}
}

// UpdateDBConnectionStats updates database connection pool metrics
func (m *MetricsService) UpdateDBConnectionStats(stats sql.DBStats) {
	m.dbConnectionsActive.WithLabelValues().Set(float64(stats.InUse))
	m.dbConnectionsIdle.WithLabelValues().Set(float64(stats.Idle))
}

// ObserveAppError records application error metrics
func (m *MetricsService) ObserveAppError(errorType, component string) {
	m.appErrorsTotal.WithLabelValues(errorType, component).Inc()
}

// RecordTaskCreated records task creation metrics
func (m *MetricsService) RecordTaskCreated(agentType, priority string) {
	m.appTasksCreatedTotal.WithLabelValues(agentType, priority).Inc()
}

// RecordTaskCompleted records task completion metrics
func (m *MetricsService) RecordTaskCompleted(agentType, priority string) {
	m.appTasksCompletedTotal.WithLabelValues(agentType, priority).Inc()
}

// UpdateUserCount updates user count metrics
func (m *MetricsService) UpdateUserCount(planType string, count int64) {
	m.appUsersTotal.WithLabelValues(planType).Set(float64(count))
}

// UpdateActiveAgents updates active agent metrics
func (m *MetricsService) UpdateActiveAgents(agentRole string, count int) {
	m.appAgentsActive.WithLabelValues(agentRole).Set(float64(count))
}

// RecordPHReferral records Product Hunt referral metrics
func (m *MetricsService) RecordPHReferral(sourceType string) {
	m.phReferralsTotal.WithLabelValues(sourceType).Inc()
}

// RecordPHUpvote records Product Hunt upvote metrics
func (m *MetricsService) RecordPHUpvote() {
	hour := time.Now().Format("15") // Hour in 24h format
	m.phUpvotesTotal.WithLabelValues(hour).Inc()
}

// RecordPHConversion records Product Hunt conversion metrics
func (m *MetricsService) RecordPHConversion(conversionType string) {
	m.phConversionsTotal.WithLabelValues(conversionType).Inc()
}

// GetMetricsHandler returns Prometheus metrics HTTP handler
func (m *MetricsService) GetMetricsHandler() http.Handler {
	return promhttp.Handler()
}

// StartMetricsCollection starts periodic metrics collection
func (m *MetricsService) StartMetricsCollection(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.collectSystemMetrics()
		case <-ctx.Done():
			return
		}
	}
}

// collectSystemMetrics collects system-level metrics
func (m *MetricsService) collectSystemMetrics() {
	// Update database connection stats
	sqlDB, err := m.db.DB()
	if err == nil {
		stats := sqlDB.Stats()
		m.UpdateDBConnectionStats(stats)
	}

	// Collect other metrics...
}

// GetProductHuntMetrics returns current Product Hunt metrics
func (m *MetricsService) GetProductHuntMetrics() map[string]interface{} {
	return map[string]interface{}{
		"referrals_total":       m.phReferralsTotal,
		"upvotes_total":         m.phUpvotesTotal,
		"conversions_total":     m.phConversionsTotal,
		"recent_upvotes_count":  m.phUpvotesTotal,
	}
}
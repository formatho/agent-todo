package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/gin-gonic/gin"
)

// Simple metrics service for testing without database dependencies
type TestMetricsService struct {
	httpRequestsTotal    *prometheus.CounterVec
	httpRequestDuration   *prometheus.HistogramVec
	phReferralsTotal      *prometheus.CounterVec
	phUpvotesTotal        *prometheus.CounterVec
}

func NewTestMetricsService() *TestMetricsService {
	return &TestMetricsService{
		httpRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "endpoint"},
		),
		httpRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "endpoint"},
		),
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
	}
}

func (m *TestMetricsService) ObserveHTTPRequest(method, endpoint string, duration time.Duration) {
	m.httpRequestsTotal.WithLabelValues(method, endpoint).Inc()
	m.httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
}

func (m *TestMetricsService) RecordPHReferral(sourceType string) {
	m.phReferralsTotal.WithLabelValues(sourceType).Inc()
}

func (m *TestMetricsService) RecordPHUpvote() {
	hour := time.Now().Format("15")
	m.phUpvotesTotal.WithLabelValues(hour).Inc()
}

func (m *TestMetricsService) GetMetricsHandler() http.Handler {
	return promhttp.Handler()
}

func main() {
	// Create test metrics service
	metricsService := NewTestMetricsService()
	
	// Create router
	router := gin.Default()
	
	// Add metrics middleware
	router.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		
		metricsService.ObserveHTTPRequest(
			c.Request.Method,
			c.Request.URL.Path,
			duration,
		)
	})
	
	// Test endpoint
	router.GET("/test", func(c *gin.Context) {
		start := time.Now()
		c.JSON(http.StatusOK, gin.H{
			"message": "Monitoring test endpoint",
			"timestamp": time.Now().Format(time.RFC3339),
		})
		duration := time.Since(start)
		metricsService.ObserveHTTPRequest("GET", "/test", duration)
	})
	
	// Health endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})
	
	// Metrics endpoint
	router.GET("/metrics", gin.WrapH(metricsService.GetMetricsHandler()))
	
	// Product Hunt test endpoint
	router.POST("/analytics/product-hunt-event", func(c *gin.Context) {
		type TestEvent struct {
			EventType string                 `json:"event_type"`
			Source    string                 `json:"source"`
			Metadata  map[string]interface{} `json:"metadata"`
		}
		
		var event TestEvent
		if err := c.ShouldBindJSON(&event); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		// Record the event
		if event.EventType == "referral" {
			metricsService.RecordPHReferral(event.Source)
		} else if event.EventType == "upvote" {
			metricsService.RecordPHUpvote()
		}
		
		c.JSON(http.StatusCreated, gin.H{
			"status": "success",
			"event": event,
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})
	
	port := "8081"
	fmt.Printf("Starting test monitoring server on port %s...\n", port)
	fmt.Printf("Endpoints:\n")
	fmt.Printf("  GET  http://localhost:%s/health - Health check\n", port)
	fmt.Printf("  GET  http://localhost:%s/test - Test endpoint\n", port)
	fmt.Printf("  GET  http://localhost:%s/metrics - Prometheus metrics\n", port)
	fmt.Printf("  POST http://localhost:%s/analytics/product-hunt-event - PH event tracking\n", port)
	fmt.Printf("\nPress Ctrl+C to stop the server\n")
	
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
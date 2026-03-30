package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/formatho/agent-todo/services"
	"github.com/gin-gonic/gin"
)

func main() {
	// Create metrics service for testing
	metricsService := services.NewMetricsService()
	
	// Create router
	router := gin.Default()
	
	// Add metrics middleware
	router.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		
		// Record basic metrics
		metricsService.ObserveHTTPRequest(
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			duration,
			0, // response size would need to be captured
		)
	})
	
	// Test endpoint
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Monitoring test endpoint",
			"timestamp": time.Now().Format(time.RFC3339),
		})
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
			EventType  string                 `json:"event_type"`
			Source     string                 `json:"source"`
			Metadata   map[string]interface{} `json:"metadata"`
		}
		
		var event TestEvent
		if err := c.ShouldBindJSON(&event); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		// Record the event
		metricsService.RecordPHReferral(event.Source)
		
		c.JSON(http.StatusCreated, gin.H{
			"status": "success",
			"event": event,
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})
	
	port := "8081"
	fmt.Printf("Starting test server on port %s...\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
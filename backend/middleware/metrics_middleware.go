package middleware

import (
	"bufio"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/formatho/agent-todo/services"
	"github.com/gin-gonic/gin"
)

// MetricsMiddleware collects HTTP request metrics
type MetricsMiddleware struct {
	metricsService *services.MetricsService
}

// NewMetricsMiddleware creates a new metrics middleware
func NewMetricsMiddleware(metricsService *services.MetricsService) *MetricsMiddleware {
	return &MetricsMiddleware{
		metricsService: metricsService,
	}
}

// Middleware returns the gin middleware function
func (m *MetricsMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// Wrap the ResponseWriter to capture status code and response size
		wrappedWriter := &responseWriter{
			ResponseWriter: c.Writer,
			statusCode:     http.StatusOK,
		}
		c.Writer = wrappedWriter

		// Process the request
		c.Next()

		// Calculate metrics
		duration := time.Since(start)
		statusCode := wrappedWriter.statusCode
		responseSize := wrappedResponseSize

		// Record metrics
		m.metricsService.ObserveHTTPRequest(
			c.Request.Method,
			c.FullPath(),
			statusCode,
			duration,
			responseSize,
		)

		// Record error metrics if there were errors
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				errorType := "unknown"
				if err.Type == gin.ErrorTypePublic {
					errorType = "public"
				} else if err.Type == gin.ErrorTypePrivate {
					errorType = "private"
				}
				m.metricsService.ObserveHTTPError(
					c.Request.Method,
					c.FullPath(),
					errorType,
				)
			}
		}
	}
}

// responseWriter wraps gin.ResponseWriter to capture status code and response size
type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	writtenBytes int64
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write captures the number of bytes written
func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	if err == nil {
		rw.writtenBytes += int64(n)
	}
	return n, err
}

// GetStatusCode returns the captured status code
func (rw *responseWriter) GetStatusCode() int {
	return rw.statusCode
}

// GetWrittenBytes returns the number of bytes written
func (rw *responseWriter) GetWrittenBytes() int64 {
	return rw.writtenBytes
}

var wrappedResponseSize int64 = 0

// MetricsResponseWriter is a global response writer for wrapped response size
type MetricsResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int64
}

func (w *MetricsResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *MetricsResponseWriter) Write(b []byte) (int, error) {
	w.size += int64(len(b))
	return w.ResponseWriter.Write(b)
}

// ProductHuntMiddleware detects and tracks Product Hunt traffic
func ProductHuntMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userAgent := c.GetHeader("User-Agent")
		referrer := c.GetHeader("Referer")
		
		// Check if this is Product Hunt traffic
		if strings.Contains(strings.ToLower(referrer), "producthunt") || 
		   strings.Contains(strings.ToLower(userAgent), "producthunt") {
			// This would typically use the metrics service, but we'll set it as a header
			// for the handler to track
			c.Set("product_hunt_referral", true)
			c.Set("ph_source_type", "direct_referral")
		}

		c.Next()
	}
}

// UptimeMiddleware tracks service uptime and health
func UptimeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// Add health check endpoint
		if c.Request.URL.Path == "/health" {
			c.JSON(http.StatusOK, gin.H{
				"status": "healthy",
				"timestamp": time.Now().Format(time.RFC3339),
				"uptime": time.Since(start).String(),
			})
			return
		}

		c.Next()
	}
}

// ErrorTrackingMiddleware captures and tracks application errors
func ErrorTrackingMiddleware(metricsService *services.MetricsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Track errors that occurred during the request
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				errorType := "unknown"
				component := "http_handler"
				
				if err.Type == gin.ErrorTypePublic {
					errorType = "public"
				} else if err.Type == gin.ErrorTypePrivate {
					errorType = "private"
				} else if err.Type == gin.ErrorTypeBind {
					errorType = "validation"
				} else if err.Type == gin.ErrorTypeRender {
					errorType = "render"
				}

				// Try to determine component from path
				path := c.FullPath()
				if strings.Contains(path, "/api/") {
					component = "api"
				} else if strings.Contains(path, "/agent/") {
					component = "agent"
				} else if strings.Contains(path, "/analytics/") {
					component = "analytics"
				}

				metricsService.ObserveAppError(errorType, component)
			}
		}
	}
}
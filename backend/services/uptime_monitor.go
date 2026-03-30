package services

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// UptimeMonitor handles external service monitoring
type UptimeMonitor struct {
	checkInterval time.Duration
	checkTimeout  time.Duration
	endpoints     []EndpointCheck
	alerts        chan Alert
}

// EndpointCheck represents a service endpoint to monitor
type EndpointCheck struct {
	Name     string
	URL      string
	Method   string
	Expected int
	Headers  map[string]string
}

// Alert represents a monitoring alert
type Alert struct {
	Type        string
	Message     string
	Severity    string
	Timestamp   time.Time
	Service     string
}

// NewUptimeMonitor creates a new uptime monitor
func NewUptimeMonitor(checkInterval, checkTimeout time.Duration) *UptimeMonitor {
	return &UptimeMonitor{
		checkInterval: checkInterval,
		checkTimeout:  checkTimeout,
		endpoints:     make([]EndpointCheck, 0),
		alerts:        make(chan Alert, 100),
	}
}

// AddEndpoint adds an endpoint to monitor
func (m *UptimeMonitor) AddEndpoint(name, url string, method string, expectedCode int, headers map[string]string) {
	m.endpoints = append(m.endpoints, EndpointCheck{
		Name:     name,
		URL:      url,
		Method:   method,
		Expected: expectedCode,
		Headers:  headers,
	})
}

// Start starts the monitoring process
func (m *UptimeMonitor) Start(ctx context.Context) {
	ticker := time.NewTicker(m.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.runChecks()
		case <-ctx.Done():
			return
		}
	}
}

// runChecks performs all endpoint checks
func (m *UptimeMonitor) runChecks() {
	for _, endpoint := range m.endpoints {
		go m.checkEndpoint(endpoint)
	}
}

// checkEndpoint checks a single endpoint
func (m *UptimeMonitor) checkEndpoint(endpoint EndpointCheck) {
	client := &http.Client{
		Timeout: m.checkTimeout,
	}

	req, err := http.NewRequest(endpoint.Method, endpoint.URL, nil)
	if err != nil {
		m.sendAlert(Alert{
			Type:      "request_error",
			Message:   fmt.Sprintf("Failed to create request for %s: %v", endpoint.Name, err),
			Severity:  "error",
			Timestamp: time.Now(),
			Service:   endpoint.Name,
		})
		return
	}

	// Add headers
	for key, value := range endpoint.Headers {
		req.Header.Set(key, value)
	}

	start := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(start)

	if err != nil {
		m.sendAlert(Alert{
			Type:      "connection_error",
			Message:   fmt.Sprintf("Failed to connect to %s: %v", endpoint.Name, err),
			Severity:  "critical",
			Timestamp: time.Now(),
			Service:   endpoint.Name,
		})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != endpoint.Expected {
		m.sendAlert(Alert{
			Type:      "status_code_mismatch",
			Message:   fmt.Sprintf("Expected %d for %s, got %d", endpoint.Expected, endpoint.Name, resp.StatusCode),
			Severity:  "warning",
			Timestamp: time.Now(),
			Service:   endpoint.Name,
		})
		return
	}

	// Success - log good health
	m.sendAlert(Alert{
		Type:      "health_check",
		Message:   fmt.Sprintf("%s is healthy (%d) in %v", endpoint.Name, resp.StatusCode, duration),
		Severity:  "info",
		Timestamp: time.Now(),
		Service:   endpoint.Name,
	})
}

// sendAlert sends an alert (in a real implementation, this would integrate with Slack, email, etc.)
func (m *UptimeMonitor) sendAlert(alert Alert) {
	select {
	case m.alerts <- alert:
		// Alert sent successfully
	default:
		// Alert channel is full, dropping oldest alert
		<-m.alerts
		m.alerts <- alert
	}
}

// GetAlerts returns the alerts channel
func (m *UptimeMonitor) GetAlerts() <-chan Alert {
	return m.alerts
}

// GetSystemHealth returns overall system health status
func (m *UptimeMonitor) GetSystemHealth() map[string]interface{} {
	health := make(map[string]interface{})
	health["timestamp"] = time.Now().Format(time.RFC3339)
	health["monitored_endpoints"] = len(m.endpoints)
	health["last_check"] = time.Now().Add(-m.checkInterval).Format(time.RFC3339)
	
	// In a real implementation, you'd track actual health status
	health["status"] = "healthy"
	health["overall_health"] = "good"
	
	return health
}

// ConfigureDefaultEndpoints configures monitoring for the default agent-todo endpoints
func (m *UptimeMonitor) ConfigureDefaultEndpoints(baseURL string) {
	// Frontend health check
	m.AddEndpoint("frontend_health", baseURL+"/health", "GET", http.StatusOK, nil)
	
	// Backend API health check
	m.AddEndpoint("backend_health", baseURL+"/api/health", "GET", http.StatusOK, nil)
	
	// Analytics endpoint
	m.AddEndpoint("analytics_endpoint", baseURL+"/analytics/track", "POST", http.StatusCreated, nil)
	
	// Metrics endpoint (Prometheus)
	m.AddEndpoint("metrics_endpoint", baseURL+"/metrics", "GET", http.StatusOK, nil)
}

// ConfigureProductHuntEndpoints configures monitoring specific to Product Hunt launch
func (m *UptimeMonitor) ConfigureProductHuntEndpoints(baseURL string) {
	// Product Hunt analytics endpoints
	m.AddEndpoint("ph_event_tracking", baseURL+"/analytics/product-hunt-event", "POST", http.StatusCreated, nil)
	m.AddEndpoint("ph_metrics", baseURL+"/analytics/product-hunt-metrics", "GET", http.StatusOK, nil)
	
	// High-priority endpoints for launch day
	m.AddEndpoint("auth_api", baseURL+"/api/auth/login", "POST", http.StatusOK, nil)
	m.AddEndpoint("agent_api", baseURL+"/api/agent/tasks", "GET", http.StatusOK, nil)
	m.AddEndpoint("task_creation", baseURL+"/api/agent/tasks", "POST", http.StatusCreated, nil)
}
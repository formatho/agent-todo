package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestRateLimiter_AllowsRequestsWithinLimit(t *testing.T) {
	limiter := NewRateLimiter(5)
	clientID := "test-client"

	// Should allow 5 requests
	for i := 0; i < 5; i++ {
		assert.True(t, limiter.Allow(clientID), "Request %d should be allowed", i+1)
	}

	// 6th request should be denied
	assert.False(t, limiter.Allow(clientID), "Request 6 should be denied")
}

func TestRateLimiter_ResetsAfterWindow(t *testing.T) {
	limiter := NewRateLimiter(2)
	limiter.window = 100 * time.Millisecond // Short window for testing
	clientID := "test-client-reset"

	// Use up the limit
	assert.True(t, limiter.Allow(clientID))
	assert.True(t, limiter.Allow(clientID))
	assert.False(t, limiter.Allow(clientID))

	// Wait for window to reset
	time.Sleep(150 * time.Millisecond)

	// Should allow requests again
	assert.True(t, limiter.Allow(clientID), "Request after window reset should be allowed")
}

func TestRateLimiter_TracksMultipleClients(t *testing.T) {
	limiter := NewRateLimiter(2)

	// Client 1
	assert.True(t, limiter.Allow("client1"))
	assert.True(t, limiter.Allow("client1"))
	assert.False(t, limiter.Allow("client1"))

	// Client 2 should still have quota
	assert.True(t, limiter.Allow("client2"))
	assert.True(t, limiter.Allow("client2"))
	assert.False(t, limiter.Allow("client2"))
}

func TestRateLimiter_CleanupOldEntries(t *testing.T) {
	limiter := NewRateLimiter(10)
	limiter.cleanup = 50 * time.Millisecond // Short cleanup interval for testing
	limiter.window = 50 * time.Millisecond  // Short window so entries become "old" quickly

	// Add some entries
	limiter.Allow("old-client")

	// Wait for both window to expire and cleanup to run
	time.Sleep(120 * time.Millisecond)

	// Make another request to trigger cleanup
	limiter.Allow("new-client")

	// Old entry should be cleaned up
	limiter.mu.Lock()
	_, exists := limiter.requests["old-client"]
	limiter.mu.Unlock()

	assert.False(t, exists, "Old entries should be cleaned up")
}

func TestRateLimitMiddleware_BlocksExcessiveRequests(t *testing.T) {
	router := gin.New()
	router.Use(RateLimitMiddleware(3))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// First 3 requests should succeed
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "Request %d should succeed", i+1)
	}

	// 4th request should be rate limited
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "Rate limit exceeded")
}

func TestRateLimitByAPIKey_TracksPerAPIKey(t *testing.T) {
	router := gin.New()
	router.Use(RateLimitByAPIKey(2))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API key 1 - 2 requests allowed, 3rd denied
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("X-API-KEY", "key1")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-API-KEY", "key1")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)

	// API key 2 should still work
	req = httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-API-KEY", "key2")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRateLimitByAPIKey_SkipsIfNoAPIKey(t *testing.T) {
	router := gin.New()
	router.Use(RateLimitByAPIKey(1))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Requests without API key should pass through (not rate limited by this middleware)
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "Request without API key should pass through")
	}
}

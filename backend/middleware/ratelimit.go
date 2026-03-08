package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements a simple token bucket rate limiter
type RateLimiter struct {
	mu        sync.Mutex
	requests  map[string]*clientInfo
	rate      int           // requests per window
	window    time.Duration // time window
	cleanup   time.Duration // cleanup interval
	lastClean time.Time
}

type clientInfo struct {
	count     int
	windowEnd time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requestsPerMinute int) *RateLimiter {
	return &RateLimiter{
		requests:  make(map[string]*clientInfo),
		rate:      requestsPerMinute,
		window:    time.Minute,
		cleanup:   5 * time.Minute,
		lastClean: time.Now(),
	}
}

// Allow checks if the client is allowed to make a request
func (rl *RateLimiter) Allow(clientID string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Cleanup old entries periodically
	if now.Sub(rl.lastClean) > rl.cleanup {
		for id, info := range rl.requests {
			if now.After(info.windowEnd) {
				delete(rl.requests, id)
			}
		}
		rl.lastClean = now
	}

	info, exists := rl.requests[clientID]
	if !exists || now.After(info.windowEnd) {
		// New window
		rl.requests[clientID] = &clientInfo{
			count:     1,
			windowEnd: now.Add(rl.window),
		}
		return true
	}

	if info.count >= rl.rate {
		return false
	}

	info.count++
	return true
}

// RateLimitMiddleware limits requests per client IP
func RateLimitMiddleware(requestsPerMinute int) gin.HandlerFunc {
	limiter := NewRateLimiter(requestsPerMinute)

	return func(c *gin.Context) {
		clientID := c.ClientIP()

		if !limiter.Allow(clientID) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitByAPIKey limits requests per API key (for agents)
func RateLimitByAPIKey(requestsPerMinute int) gin.HandlerFunc {
	limiters := make(map[string]*RateLimiter)
	var mu sync.Mutex

	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-KEY")
		if apiKey == "" {
			// No API key, skip this middleware (will be caught by auth middleware)
			c.Next()
			return
		}

		mu.Lock()
		limiter, exists := limiters[apiKey]
		if !exists {
			limiter = NewRateLimiter(requestsPerMinute)
			limiters[apiKey] = limiter
		}
		mu.Unlock()

		if !limiter.Allow(apiKey) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

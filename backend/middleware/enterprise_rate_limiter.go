package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/formatho/agent-todo/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/go-redis/redis/v8"
)

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	MaxRequests     int           `json:"max_requests"`
	WindowPeriod    time.Duration `json:"window_period"`
	KeyGenerator    func(*fiber.Ctx) string
	LimitReached    func(*fiber.Ctx) error
	SkipCondition   func(*fiber.Ctx) bool
	StorageBackend  string        `json:"storage_backend"` // "memory", "redis"
}

// EnterpriseRateLimiter provides advanced rate limiting for enterprise customers
type EnterpriseRateLimiter struct {
	config    *RateLimitConfig
	redis     *redis.Client
	inMemory  map[string]*ClientInfo
	mu        sync.RWMutex
	ctx       context.Context
}

// ClientInfo stores client rate limit information
type ClientInfo struct {
	Requests   int       `json:"requests"`
	ResetAt    time.Time `json:"reset_at"`
	Blocked    bool      `json:"blocked"`
	BlockedAt  time.Time `json:"blocked_at"`
}

// NewEnterpriseRateLimiter creates a new enterprise rate limiter
func NewEnterpriseRateLimiter(config *RateLimitConfig, redisClient *redis.Client) *EnterpriseRateLimiter {
	rl := &EnterpriseRateLimiter{
		config:   config,
		redis:    redisClient,
		inMemory: make(map[string]*ClientInfo),
		ctx:      context.Background(),
	}

	// Start cleanup goroutine for in-memory storage
	go rl.cleanupExpiredClients()

	return rl
}

// Middleware returns the Fiber middleware handler
func (erl *EnterpriseRateLimiter) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip rate limiting if condition met
		if erl.config.SkipCondition != nil && erl.config.SkipCondition(c) {
			return c.Next()
		}

		// Generate rate limit key
		key := erl.config.KeyGenerator(c)
		if key == "" {
			key = c.IP()
		}

		// Check rate limit
		allowed, remaining, resetTime, err := erl.checkRateLimit(key)
		if err != nil {
			// On error, allow request (fail open)
			return c.Next()
		}

		// Set rate limit headers
		c.Set("X-RateLimit-Limit", strconv.Itoa(erl.config.MaxRequests))
		c.Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Set("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

		if !allowed {
			// Rate limit exceeded
			if erl.config.LimitReached != nil {
				return erl.config.LimitReached(c)
			}
			return c.Status(http.StatusTooManyRequests).JSON(fiber.Map{
				"error":       "rate limit exceeded",
				"retry_after": resetTime.Sub(time.Now()).Seconds(),
			})
		}

		return c.Next()
	}
}

// checkRateLimit checks if a client has exceeded their rate limit
func (erl *EnterpriseRateLimiter) checkRateLimit(key string) (allowed bool, remaining int, resetTime time.Time, err error) {
	switch erl.config.StorageBackend {
	case "redis":
		return erl.checkRedisRateLimit(key)
	default:
		return erl.checkMemoryRateLimit(key)
	}
}

// checkRedisRateLimit checks rate limit using Redis
func (erl *EnterpriseRateLimiter) checkRedisRateLimit(key string) (bool, int, time.Time, error) {
	ctx := context.Background()
	now := time.Now()
	resetTime := now.Add(erl.config.WindowPeriod)

	// Get current request count
	countStr, err := erl.redis.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return false, 0, resetTime, err
	}

	count := 0
	if countStr != "" {
		count, _ = strconv.Atoi(countStr)
	}

	// Check if limit exceeded
	if count >= erl.config.MaxRequests {
		// Get TTL for reset time
		ttl, _ := erl.redis.TTL(ctx, key).Result()
		resetTime = now.Add(ttl)
		return false, 0, resetTime, nil
	}

	// Increment counter
	pipe := erl.redis.Pipeline()
	pipe.Incr(ctx, key)
	if count == 0 {
		pipe.Expire(ctx, key, erl.config.WindowPeriod)
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return false, 0, resetTime, err
	}

	remaining := erl.config.MaxRequests - count - 1
	return true, remaining, resetTime, nil
}

// checkMemoryRateLimit checks rate limit using in-memory storage
func (erl *EnterpriseRateLimiter) checkMemoryRateLimit(key string) (bool, int, time.Time, error) {
	erl.mu.Lock()
	defer erl.mu.Unlock()

	now := time.Now()
	resetTime := now.Add(erl.config.WindowPeriod)

	// Get or create client info
	client, exists := erl.inMemory[key]
	if !exists || now.After(client.ResetAt) {
		// New window
		client = &ClientInfo{
			Requests: 0,
			ResetAt:  resetTime,
		}
		erl.inMemory[key] = client
	}

	// Check if limit exceeded
	if client.Requests >= erl.config.MaxRequests {
		return false, 0, client.ResetAt, nil
	}

	// Increment counter
	client.Requests++
	remaining := erl.config.MaxRequests - client.Requests

	return true, remaining, client.ResetAt, nil
}

// cleanupExpiredClients removes expired client entries from memory
func (erl *EnterpriseRateLimiter) cleanupExpiredClients() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		erl.mu.Lock()
		now := time.Now()
		for key, client := range erl.inMemory {
			if now.After(client.ResetAt) {
				delete(erl.inMemory, key)
			}
		}
		erl.mu.Unlock()
	}
}

// DefaultRateLimitConfig returns default rate limit configuration
func DefaultRateLimitConfig() *RateLimitConfig {
	return &RateLimitConfig{
		MaxRequests:    1000,
		WindowPeriod:   1 * time.Minute,
		StorageBackend: "redis",
		KeyGenerator: func(c *fiber.Ctx) string {
			// Use organisation ID + user ID for authenticated requests
			if orgID := c.Locals("organisation_id"); orgID != nil {
				if userID := c.Locals("user_id"); userID != nil {
					return fmt.Sprintf("rate_limit:%s:%s", orgID, userID)
				}
			}
			// Fall back to IP for unauthenticated requests
			return fmt.Sprintf("rate_limit:ip:%s", c.IP())
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(http.StatusTooManyRequests).JSON(fiber.Map{
				"error":   "rate limit exceeded",
				"message": "too many requests, please try again later",
			})
		},
		SkipCondition: func(c *fiber.Ctx) bool {
			// Skip rate limiting for health checks
			return c.Path() == "/health"
		},
	}
}

// EnterpriseRateLimitTiers defines rate limits for different subscription tiers
var EnterpriseRateLimitTiers = map[string]*RateLimitConfig{
	"free": {
		MaxRequests:    100,
		WindowPeriod:   1 * time.Minute,
		StorageBackend: "redis",
	},
	"starter": {
		MaxRequests:    500,
		WindowPeriod:   1 * time.Minute,
		StorageBackend: "redis",
	},
	"pro": {
		MaxRequests:    2000,
		WindowPeriod:   1 * time.Minute,
		StorageBackend: "redis",
	},
	"enterprise": {
		MaxRequests:    10000,
		WindowPeriod:   1 * time.Minute,
		StorageBackend: "redis",
	},
}

// GetRateLimitForTier returns rate limit config for a subscription tier
func GetRateLimitForTier(tier string) *RateLimitConfig {
	if config, exists := EnterpriseRateLimitTiers[tier]; exists {
		return config
	}
	return DefaultRateLimitConfig()
}

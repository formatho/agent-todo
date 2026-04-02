package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// CacheService provides caching functionality for enterprise scalability
type CacheService struct {
	redis *redis.Client
	db    *gorm.DB
	ttl   time.Duration
}

// NewCacheService creates a new cache service
func NewCacheService(redisClient *redis.Client, db *gorm.DB, defaultTTL time.Duration) *CacheService {
	return &CacheService{
		redis: redisClient,
		db:    db,
		ttl:   defaultTTL,
	}
}

// Get retrieves a value from cache or loads it from database
func (cs *CacheService) Get(ctx context.Context, key string, dest interface{}, loader func() error) error {
	// Try to get from Redis first
	val, err := cs.redis.Get(ctx, key).Result()
	if err == nil {
		// Cache hit - unmarshal and return
		if err := json.Unmarshal([]byte(val), dest); err == nil {
			return nil
		}
	}

	// Cache miss - load from database
	if err := loader(); err != nil {
		return err
	}

	// Store in cache for future use
	if data, err := json.Marshal(dest); err == nil {
		cs.redis.Set(ctx, key, data, cs.ttl)
	}

	return nil
}

// Set stores a value in cache
func (cs *CacheService) Set(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return cs.redis.Set(ctx, key, data, cs.ttl).Err()
}

// Delete removes a value from cache
func (cs *CacheService) Delete(ctx context.Context, key string) error {
	return cs.redis.Del(ctx, key).Err()
}

// DeletePattern removes all keys matching a pattern
func (cs *CacheService) DeletePattern(ctx context.Context, pattern string) error {
	iter := cs.redis.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := cs.redis.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

// InvalidateOrgCache invalidates all cache entries for an organisation
func (cs *CacheService) InvalidateOrgCache(ctx context.Context, orgID string) error {
	patterns := []string{
		fmt.Sprintf("org:%s:*", orgID),
		fmt.Sprintf("tasks:org:%s:*", orgID),
		fmt.Sprintf("agents:org:%s:*", orgID),
		fmt.Sprintf("projects:org:%s:*", orgID),
	}

	for _, pattern := range patterns {
		if err := cs.DeletePattern(ctx, pattern); err != nil {
			return err
		}
	}

	return nil
}

// Cache keys for common queries
const (
	CacheKeyOrgTasks     = "tasks:org:%s:page:%d"
	CacheKeyOrgAgents    = "agents:org:%s"
	CacheKeyOrgProjects  = "projects:org:%s"
	CacheKeyTaskDetails  = "task:%s"
	CacheKeyUserOrgs     = "user:%s:orgs"
	CacheKeyAgentActivity = "agent:%s:activity"
)

// GetCachedOrgTasks retrieves cached tasks for an organisation
func (cs *CacheService) GetCachedOrgTasks(ctx context.Context, orgID string, page int) ([]map[string]interface{}, error) {
	key := fmt.Sprintf(CacheKeyOrgTasks, orgID, page)
	var tasks []map[string]interface{}
	
	err := cs.Get(ctx, key, &tasks, func() error {
		// Load from database
		return cs.db.Table("task_summary").
			Where("organisation_id = ?", orgID).
			Offset((page - 1) * 20).
			Limit(20).
			Find(&tasks).Error
	})

	return tasks, err
}

// GetCachedOrgAgents retrieves cached agents for an organisation
func (cs *CacheService) GetCachedOrgAgents(ctx context.Context, orgID string) ([]map[string]interface{}, error) {
	key := fmt.Sprintf(CacheKeyOrgAgents, orgID)
	var agents []map[string]interface{}
	
	err := cs.Get(ctx, key, &agents, func() error {
		return cs.db.Table("agents").
			Where("organisation_id = ? AND enabled = ?", orgID, true).
			Find(&agents).Error
	})

	return agents, err
}

// CacheStats returns cache statistics
func (cs *CacheService) CacheStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Get Redis info
	info, err := cs.redis.Info(ctx, "stats").Result()
	if err == nil {
		stats["redis_info"] = info
	}

	// Get key count
	keys, err := cs.redis.DBSize(ctx).Result()
	if err == nil {
		stats["total_keys"] = keys
	}

	// Get memory usage
	memUsage, err := cs.redis.Info(ctx, "memory").Result()
	if err == nil {
		stats["memory_usage"] = memUsage
	}

	return stats, nil
}

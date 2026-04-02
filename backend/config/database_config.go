package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DatabaseConfig holds database configuration for enterprise deployments
type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	SSLMode         string
	
	// Connection Pool Settings
	MaxOpenConns    int           `json:"max_open_conns"`
	MaxIdleConns    int           `json:"max_idle_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time"`
	
	// Performance Settings
	PrepareStmt     bool          `json:"prepare_stmt"`
	StatementCache  int           `json:"statement_cache"`
	
	// Monitoring
	EnableMetrics   bool          `json:"enable_metrics"`
	SlowQueryThreshold time.Duration `json:"slow_query_threshold"`
}

// DefaultDatabaseConfig returns default database configuration
func DefaultDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:            getEnv("DB_HOST", "localhost"),
		Port:            getEnvAsInt("DB_PORT", 5432),
		User:            getEnv("DB_USER", "postgres"),
		Password:        getEnv("DB_PASSWORD", ""),
		Database:        getEnv("DB_NAME", "agent_todo"),
		SSLMode:         getEnv("DB_SSLMODE", "disable"),
		MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 100),
		MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 20),
		ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", 30*time.Minute),
		ConnMaxIdleTime: getEnvAsDuration("DB_CONN_MAX_IDLE_TIME", 10*time.Minute),
		PrepareStmt:     getEnvAsBool("DB_PREPARE_STMT", true),
		StatementCache:  getEnvAsInt("DB_STATEMENT_CACHE", 100),
		EnableMetrics:   getEnvAsBool("DB_ENABLE_METRICS", true),
		SlowQueryThreshold: getEnvAsDuration("DB_SLOW_QUERY_THRESHOLD", 1*time.Second),
	}
}

// EnterpriseDatabaseConfig returns enterprise-optimized database configuration
func EnterpriseDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:            getEnv("DB_HOST", "localhost"),
		Port:            getEnvAsInt("DB_PORT", 5432),
		User:            getEnv("DB_USER", "postgres"),
		Password:        getEnv("DB_PASSWORD", ""),
		Database:        getEnv("DB_NAME", "agent_todo"),
		SSLMode:         getEnv("DB_SSLMODE", "require"),
		MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 200),
		MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 50),
		ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", 1*time.Hour),
		ConnMaxIdleTime: getEnvAsDuration("DB_CONN_MAX_IDLE_TIME", 15*time.Minute),
		PrepareStmt:     true,
		StatementCache:  200,
		EnableMetrics:   true,
		SlowQueryThreshold: 500*time.Millisecond,
	}
}

// Connect establishes a database connection with the given configuration
func (dc *DatabaseConfig) Connect() (*gorm.DB, error) {
	dsn := dc.ConnectionString()
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: dc.PrepareStmt,
		// Enable logging for slow queries
		// Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(dc.MaxOpenConns)
	sqlDB.SetMaxIdleConns(dc.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(dc.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(dc.ConnMaxIdleTime)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// ConnectionString returns the PostgreSQL connection string
func (dc *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dc.Host,
		dc.Port,
		dc.User,
		dc.Password,
		dc.Database,
		dc.SSLMode,
	)
}

// ConnectionPoolStats returns current connection pool statistics
func (dc *DatabaseConfig) ConnectionPoolStats(db *gorm.DB) (map[string]interface{}, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	stats := sqlDB.Stats()

	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":              stats.InUse,
		"idle":                stats.Idle,
		"wait_count":          stats.WaitCount,
		"wait_duration_ms":    stats.WaitDuration.Milliseconds(),
		"max_idle_closed":     stats.MaxIdleClosed,
		"max_lifetime_closed": stats.MaxLifetimeClosed,
		"config": map[string]interface{}{
			"max_open_conns":    dc.MaxOpenConns,
			"max_idle_conns":    dc.MaxIdleConns,
			"conn_max_lifetime": dc.ConnMaxLifetime.String(),
			"conn_max_idle_time": dc.ConnMaxIdleTime.String(),
		},
	}, nil
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// DatabaseConfigPresets returns pre-configured database settings for different scenarios
var DatabaseConfigPresets = map[string]*DatabaseConfig{
	"development": {
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 2 * time.Minute,
		PrepareStmt:     false,
		EnableMetrics:   false,
	},
	"staging": {
		MaxOpenConns:    50,
		MaxIdleConns:    10,
		ConnMaxLifetime: 15 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
		PrepareStmt:     true,
		StatementCache:  50,
		EnableMetrics:   true,
		SlowQueryThreshold: 1 * time.Second,
	},
	"production": {
		MaxOpenConns:    100,
		MaxIdleConns:    20,
		ConnMaxLifetime: 30 * time.Minute,
		ConnMaxIdleTime: 10 * time.Minute,
		PrepareStmt:     true,
		StatementCache:  100,
		EnableMetrics:   true,
		SlowQueryThreshold: 500 * time.Millisecond,
	},
	"enterprise": {
		MaxOpenConns:    200,
		MaxIdleConns:    50,
		ConnMaxLifetime: 1 * time.Hour,
		ConnMaxIdleTime: 15 * time.Minute,
		PrepareStmt:     true,
		StatementCache:  200,
		EnableMetrics:   true,
		SlowQueryThreshold: 200 * time.Millisecond,
	},
}

// GetDatabaseConfig returns database config for a given environment
func GetDatabaseConfig(environment string) *DatabaseConfig {
	if config, exists := DatabaseConfigPresets[environment]; exists {
		return config
	}
	return DefaultDatabaseConfig()
}

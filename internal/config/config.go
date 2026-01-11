package config

import (
	"fmt"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	JWT       JWTConfig
	RateLimit RateLimitConfig
	Log       LogConfig
	Env       string
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host         string        `env:"SERVER_HOST" envDefault:"0.0.0.0"`
	Port         int           `env:"SERVER_PORT" envDefault:"8080"`
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT" envDefault:"30s"`
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT" envDefault:"30s"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	DatabaseURL     string        `env:"DATABASE_URL"` // Full database URL (takes precedence)
	Host            string        `env:"DB_HOST" envDefault:"localhost"`
	Port            int           `env:"DB_PORT" envDefault:"5432"`
	User            string        `env:"DB_USER" envDefault:"kg_user"`
	Password        string        `env:"DB_PASSWORD" envDefault:""`
	DBName          string        `env:"DB_NAME" envDefault:"knowledge_garden"`
	SSLMode         string        `env:"DB_SSL_MODE" envDefault:"disable"`
	MaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS" envDefault:"25"`
	MaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS" envDefault:"5"`
	ConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME" envDefault:"5m"`
}

// DSN returns the PostgreSQL data source name
// If DATABASE_URL is set, it takes precedence over individual DB_* variables
func (c *DatabaseConfig) DSN() string {
	if c.DatabaseURL != "" {
		return c.DatabaseURL
	}
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret            string        `env:"JWT_SECRET" envDefault:"change-me-in-production"`
	AccessExpiration  time.Duration `env:"JWT_ACCESS_EXPIRATION" envDefault:"168h"`  // 7 days for CLI apps
	RefreshExpiration time.Duration `env:"JWT_REFRESH_EXPIRATION" envDefault:"720h"` // 30 days
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Enabled  bool          `env:"RATE_LIMIT_ENABLED" envDefault:"true"`
	Requests int           `env:"RATE_LIMIT_REQUESTS" envDefault:"100"`
	Window   time.Duration `env:"RATE_LIMIT_WINDOW" envDefault:"1m"`
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level  string `env:"LOG_LEVEL" envDefault:"info"`
	Format string `env:"LOG_FORMAT" envDefault:"json"`
}

// Address returns the server address
func (c *ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

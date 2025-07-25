package redis

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// Config holds Redis configuration
type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// Client wraps redis.Client with additional functionality
type Client struct {
	*redis.Client
	logger *logrus.Logger
}

// NewClient creates a new Redis client with configuration from environment
func NewClient(logger *logrus.Logger) (*Client, error) {
	config := getConfigFromEnv()
	
	tracer := otel.Tracer("cart-service")
	ctx, span := tracer.Start(context.Background(), "redis.connect")
	defer span.End()

	span.SetAttributes(
		attribute.String("redis.host", config.Host),
		attribute.Int("redis.port", config.Port),
		attribute.Int("redis.db", config.DB),
	)

	// Create Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password:     config.Password,
		DB:           config.DB,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
		IdleTimeout:  5 * time.Minute,
		MaxRetries:   3,
	})

	// Test connection
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	if err := rdb.Ping(ctx).Err(); err != nil {
		span.RecordError(err)
		logger.WithError(err).WithFields(logrus.Fields{
			"component": "redis",
			"action":    "connect",
			"host":      config.Host,
			"port":      config.Port,
		}).Error("Failed to connect to Redis")
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	span.SetAttributes(attribute.String("redis.status", "connected"))

	logger.WithFields(logrus.Fields{
		"component": "redis",
		"action":    "connect",
		"host":      config.Host,
		"port":      config.Port,
		"db":        config.DB,
	}).Info("Connected to Redis")

	return &Client{
		Client: rdb,
		logger: logger,
	}, nil
}

// HealthCheck verifies Redis connection is healthy
func (c *Client) HealthCheck(ctx context.Context) error {
	tracer := otel.Tracer("cart-service")
	redisCtx, span := tracer.Start(ctx, "redis.health_check")
	defer span.End()

	span.SetAttributes(
		attribute.String("redis.operation", "PING"),
	)

	start := time.Now()
	pong, err := c.Ping(redisCtx).Result()
	duration := time.Since(start)

	span.SetAttributes(
		attribute.Int64("redis.duration_ms", duration.Milliseconds()),
	)

	if err != nil {
		span.RecordError(err)
		c.logger.WithError(err).WithFields(logrus.Fields{
			"component": "redis",
			"action":    "health_check",
			"duration":  duration,
		}).Error("Redis health check failed")
		return fmt.Errorf("Redis health check failed: %v", err)
	}

	if pong != "PONG" {
		err := fmt.Errorf("unexpected ping response: %s", pong)
		span.RecordError(err)
		c.logger.WithFields(logrus.Fields{
			"component": "redis",
			"action":    "health_check",
			"response":  pong,
			"duration":  duration,
		}).Error("Redis ping returned unexpected response")
		return err
	}

	span.SetAttributes(
		attribute.String("redis.status", "healthy"),
		attribute.String("redis.response", pong),
	)

	return nil
}

// Close closes the Redis connection
func (c *Client) Close() error {
	c.logger.WithFields(logrus.Fields{
		"component": "redis",
		"action":    "close",
	}).Info("Closing Redis connection")

	return c.Client.Close()
}

// GetStats returns Redis connection statistics
func (c *Client) GetStats() *redis.PoolStats {
	return c.PoolStats()
}

// getConfigFromEnv reads Redis configuration from environment variables
func getConfigFromEnv() Config {
	config := Config{
		Host:     getEnvOrDefault("REDIS_HOST", "localhost"),
		Port:     getEnvIntOrDefault("REDIS_PORT", 6379),
		Password: getEnvOrDefault("REDIS_PASSWORD", ""),
		DB:       getEnvIntOrDefault("REDIS_DB", 0),
	}
	return config
}

// getEnvOrDefault gets environment variable or returns default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvIntOrDefault gets environment variable as int or returns default value
func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
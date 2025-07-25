package redis

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"cart-service/internal/logger"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// Client wraps redis.Client with additional functionality
type Client struct {
	*redis.Client
}

// Connect establishes a connection to Redis
func Connect() (*Client, error) {
	tracer := otel.Tracer("cart-redis")
	ctx, span := tracer.Start(context.Background(), "redis.Connect")
	defer span.End()

	// Get Redis configuration from environment
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "redis.cart.svc.cluster.local"
	}

	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = "6379"
	}

	password := os.Getenv("REDIS_PASSWORD")
	// Default to no password for our lab setup

	dbStr := os.Getenv("REDIS_DB")
	db := 0
	if dbStr != "" {
		if parsed, err := strconv.Atoi(dbStr); err == nil {
			db = parsed
		}
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	span.SetAttributes(
		attribute.String("redis.host", host),
		attribute.String("redis.port", port),
		attribute.Int("redis.db", db),
	)

	logger.WithFields(logrus.Fields{
		"component": "redis",
		"action":    "connect",
		"host":      host,
		"port":      port,
		"db":        db,
	}).Info("Connecting to Redis")

	// Create Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 2,
	})

	// Test connection
	ctx, pingSpan := tracer.Start(ctx, "redis.Ping")
	pong, err := rdb.Ping(ctx).Result()
	pingSpan.End()

	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	logger.WithFields(logrus.Fields{
		"component": "redis",
		"action":    "ping",
		"response":  pong,
	}).Info("Redis connection successful")

	span.SetAttributes(attribute.String("redis.ping_response", pong))

	return &Client{Client: rdb}, nil
}

// Close closes the Redis connection
func (c *Client) Close() error {
	logger.WithFields(logrus.Fields{
		"component": "redis",
		"action":    "close",
	}).Info("Closing Redis connection")

	return c.Client.Close()
}

// HealthCheck performs a health check on the Redis connection
func (c *Client) HealthCheck(ctx context.Context) error {
	tracer := otel.Tracer("cart-redis")
	ctx, span := tracer.Start(ctx, "redis.HealthCheck")
	defer span.End()

	_, err := c.Ping(ctx).Result()
	if err != nil {
		span.RecordError(err)
		logger.WithError(err).WithFields(logrus.Fields{
			"component": "redis",
			"action":    "health_check",
		}).Error("Redis health check failed")
		return err
	}

	span.SetAttributes(attribute.Bool("redis.healthy", true))
	return nil
}
package main

import (
	"os"

	"cart-service/internal/logger"
	"cart-service/internal/redis"
	"cart-service/internal/server"
	"cart-service/internal/tracing"

	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize OpenTelemetry tracing
	cleanup, err := tracing.Setup("cart-service")
	if err != nil {
		logger.WithError(err).WithFields(logrus.Fields{
			"component": "tracing",
			"action":    "setup",
		}).Fatal("Failed to initialize tracing")
	}
	defer cleanup()

	logger.WithFields(logrus.Fields{
		"component": "tracing",
		"action":    "initialize",
	}).Info("OpenTelemetry tracing initialized")

	// Get Redis connection
	redisClient, err := redis.Connect()
	if err != nil {
		logger.WithError(err).WithFields(logrus.Fields{
			"component": "redis",
			"action":    "connect",
		}).Fatal("Failed to connect to Redis")
	}
	defer redisClient.Close()

	logger.WithFields(logrus.Fields{
		"component": "redis",
		"action":    "connect",
	}).Info("Connected to Redis")

	// Create server with Redis client
	srv := server.NewServer(redisClient)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	logger.WithFields(logrus.Fields{
		"component": "server",
		"action":    "start",
		"port":      port,
	}).Info("Starting cart service")

	if err := srv.Start(port); err != nil {
		logger.WithError(err).WithFields(logrus.Fields{
			"component": "server",
			"action":    "start",
		}).Fatal("Failed to start server")
	}
}
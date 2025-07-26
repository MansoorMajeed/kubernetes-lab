package main

import (
	"os"

	"cart-service/internal/redis"
	"cart-service/internal/server"

	"github.com/sirupsen/logrus"
)

func main() {
	// Create logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})

	logger.WithFields(logrus.Fields{
		"component": "main",
		"action":    "start",
		"service":   "cart-service",
	}).Info("Starting cart service")

	// Create Redis client
	redisClient, err := redis.NewClient(logger)
	if err != nil {
		logger.WithError(err).WithFields(logrus.Fields{
			"component": "main",
			"action":    "redis_connect",
		}).Fatal("Failed to connect to Redis")
	}
	defer redisClient.Close()

	// Create and start server
	srv := server.NewServer(redisClient.Client, logger)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.WithFields(logrus.Fields{
		"component": "main",
		"action":    "server_start",
		"port":      port,
	}).Info("Starting HTTP server")

	if err := srv.Start(port); err != nil {
		logger.WithError(err).WithFields(logrus.Fields{
			"component": "main",
			"action":    "server_start",
		}).Fatal("Failed to start server")
	}
}
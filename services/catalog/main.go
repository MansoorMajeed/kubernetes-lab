package main

import (
	"os"

	"catalog-service/internal/db"
	"catalog-service/internal/logger"
	"catalog-service/internal/server"
	"catalog-service/internal/tracing"

	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize OpenTelemetry tracing
	cleanup, err := tracing.Setup("catalog-service")
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

	// Get database connection
	database, err := db.Connect()
	if err != nil {
		logger.WithError(err).WithFields(logrus.Fields{
			"component": "database",
			"action":    "connect",
		}).Fatal("Failed to connect to database")
	}
	defer database.Close()

	// Initialize database schema
	if err := database.InitSchema(); err != nil {
		logger.WithError(err).WithFields(logrus.Fields{
			"component": "database",
			"action":    "schema_init",
		}).Fatal("Failed to initialize database schema")
	}

	// Create server with the underlying sql.DB
	srv := server.NewServer(database.DB)

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
	}).Info("Starting catalog service")

	if err := srv.Start(port); err != nil {
		logger.WithError(err).WithFields(logrus.Fields{
			"component": "server",
			"action":    "start",
		}).Fatal("Failed to start server")
	}
}

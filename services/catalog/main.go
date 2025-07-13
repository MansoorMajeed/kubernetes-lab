package main

import (
	"log"
	"os"

	"catalog-service/internal/db"
	"catalog-service/internal/server"
)

func main() {
	// Initialize database connection
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Initialize database schema
	if err := database.InitSchema(); err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}

	// Create server with database
	srv := server.New(database)

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Starting catalog service on port %s", port)

	if err := srv.Start(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

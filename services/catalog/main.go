package main

import (
	"log"
	"os"

	"catalog-service/internal/server"
)

func main() {
	// Create server
	srv := server.New()

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

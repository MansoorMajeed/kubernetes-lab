package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"catalog-service/internal/db"
	"catalog-service/internal/handlers"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	server *http.Server
	db     *db.Database
}

func New(database *db.Database) *Server {
	// Set gin to release mode for production
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Add basic middleware
	router.Use(gin.Recovery())
	router.Use(LoggingMiddleware())

	// Create server instance
	s := &Server{
		router: router,
		db:     database,
	}

	// Setup routes
	s.setupRoutes()

	return s
}

func (s *Server) setupRoutes() {
	// Health check endpoint with database
	s.router.GET("/health", handlers.HealthCheckWithDB(s.db))

	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		// Products routes will go here
		_ = v1 // avoid unused variable warning for now
	}
}

func (s *Server) Start(addr string) error {
	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.server != nil {
		return s.server.Shutdown(ctx)
	}
	return nil
}

// LoggingMiddleware provides basic request logging
func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Simple structured log format for now
		log.Printf(`{"method":"%s","path":"%s","status":%d,"latency":"%v","ip":"%s"}`,
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.ClientIP,
		)
		return ""
	})
}

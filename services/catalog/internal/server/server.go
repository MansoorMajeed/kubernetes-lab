package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"catalog-service/internal/handlers"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	server *http.Server
}

func New() *Server {
	// Set gin to release mode for production
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Add basic middleware
	router.Use(gin.Recovery())
	router.Use(LoggingMiddleware())

	// Create server instance
	s := &Server{
		router: router,
	}

	// Setup routes
	s.setupRoutes()

	return s
}

func (s *Server) setupRoutes() {
	// Health check endpoint
	s.router.GET("/health", handlers.HealthCheck)

	// API v1 routes (for future)
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

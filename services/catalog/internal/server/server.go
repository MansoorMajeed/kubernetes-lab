package server

import (
	"database/sql"
	"strconv"
	"time"

	"catalog-service/internal/handlers"
	"catalog-service/internal/logger"
	"catalog-service/internal/metrics"
	"catalog-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// Server represents the HTTP server
type Server struct {
	router  *gin.Engine
	db      *sql.DB
	metrics *metrics.HTTPMetrics
}

// NewServer creates a new server instance
func NewServer(database *sql.DB) *Server {
	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)

	// Create router without default middleware (no default logging)
	router := gin.New()

	// Add recovery middleware (but not logging - we'll add our own)
	router.Use(gin.Recovery())

	// Initialize metrics
	httpMetrics := metrics.NewHTTPMetrics()

	server := &Server{
		router:  router,
		db:      database,
		metrics: httpMetrics,
	}

	// Add our custom middleware
	server.router.Use(server.loggingMiddleware())
	server.router.Use(server.metricsMiddleware())

	// Setup routes
	server.setupRoutes()

	return server
}

// loggingMiddleware logs HTTP requests with structured JSON
func (s *Server) loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Skip logging for health check endpoints to reduce noise
		if c.Request.URL.Path == "/health" {
			return
		}

		// Log request details
		duration := time.Since(start)

		logger.WithFields(logrus.Fields{
			"component":   "http",
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"status_code": c.Writer.Status(),
			"duration_ms": duration.Milliseconds(),
			"client_ip":   c.ClientIP(),
			"user_agent":  c.Request.UserAgent(),
		}).Info("HTTP request processed")
	}
}

// metricsMiddleware collects HTTP metrics for Prometheus
func (s *Server) metricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip metrics collection for the metrics endpoint itself
		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		// Record in-flight request
		s.metrics.IncInFlight()
		defer s.metrics.DecInFlight()

		// Record start time
		start := time.Now()

		// Process request
		c.Next()

		// Calculate duration and record metrics
		duration := time.Since(start).Seconds()
		statusCode := strconv.Itoa(c.Writer.Status())

		s.metrics.RecordRequest(
			c.Request.Method,
			c.FullPath(), // Use route pattern instead of actual path (e.g., "/api/v1/products/:id")
			statusCode,
			duration,
		)
	}
}

// setupRoutes configures all the routes for the server
func (s *Server) setupRoutes() {
	// Health check endpoint
	healthHandler := handlers.NewHealthHandler(s.db)
	s.router.GET("/health", healthHandler.HealthCheck)

	// Metrics endpoint for Prometheus
	s.router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Create product service and handler
	productService := models.NewProductService(s.db)
	productHandler := handlers.NewProductHandler(productService)

	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		// Product routes
		products := v1.Group("/products")
		{
			products.GET("", productHandler.GetProducts)          // GET /api/v1/products
			products.POST("", productHandler.CreateProduct)       // POST /api/v1/products
			products.GET("/:id", productHandler.GetProduct)       // GET /api/v1/products/:id
			products.PUT("/:id", productHandler.UpdateProduct)    // PUT /api/v1/products/:id
			products.DELETE("/:id", productHandler.DeleteProduct) // DELETE /api/v1/products/:id
		}
	}

	// Log all registered routes
	logger.WithFields(logrus.Fields{
		"component": "server",
		"action":    "route_register",
	}).Info("Registered routes")

	for _, route := range s.router.Routes() {
		logger.WithFields(logrus.Fields{
			"component": "server",
			"action":    "route_register",
			"method":    route.Method,
			"path":      route.Path,
		}).Debug("Route registered")
	}
}

// Start starts the HTTP server
func (s *Server) Start(port string) error {
	logger.WithFields(logrus.Fields{
		"component": "server",
		"action":    "start",
		"port":      port,
	}).Info("Starting server")

	return s.router.Run(":" + port)
}

// Stop gracefully stops the server
func (s *Server) Stop() error {
	if s.db != nil {
		logger.WithFields(logrus.Fields{
			"component": "server",
			"action":    "shutdown",
		}).Info("Closing database connection")

		return s.db.Close()
	}
	return nil
}

// GetDB returns the database connection (for testing purposes)
func (s *Server) GetDB() *sql.DB {
	return s.db
}

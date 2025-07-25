package server

import (
	"database/sql"
	"net"
	"strconv"
	"time"

	"catalog-service/internal/grpc"
	"catalog-service/internal/handlers"
	"catalog-service/internal/logger"
	"catalog-service/internal/metrics"
	"catalog-service/internal/models"
	"catalog-service/internal/services"

	catalogpb "github.com/mansoormajeed/kubernetes-lab/proto/catalog"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"
	grpcServer "google.golang.org/grpc"
)

// Server represents the HTTP and gRPC server
type Server struct {
	router     *gin.Engine
	grpcServer *grpcServer.Server
	db         *sql.DB
	metrics    *metrics.HTTPMetrics
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

	// Create gRPC server (OpenTelemetry instrumentation added separately)
	grpcSrv := grpcServer.NewServer()

	// Register gRPC service
	catalogGRPCServer := grpc.NewCatalogGRPCServer(database)
	catalogpb.RegisterCatalogServiceServer(grpcSrv, catalogGRPCServer)

	server := &Server{
		router:     router,
		grpcServer: grpcSrv,
		db:         database,
		metrics:    httpMetrics,
	}

	// Add middleware in order:
	// 1. OpenTelemetry tracing (creates spans)
	router.Use(otelgin.Middleware("catalog-service"))

	// 2. Our custom logging middleware (can use trace context)
	server.router.Use(server.loggingMiddleware())

	// 3. Our metrics middleware
	server.router.Use(server.metricsMiddleware())

	// Setup routes
	server.setupRoutes()

	return server
}

// loggingMiddleware logs HTTP requests with structured JSON and trace correlation
func (s *Server) loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Skip logging for health check endpoints to reduce noise
		if c.Request.URL.Path == "/health" {
			return
		}

		// Get trace information for correlation
		spanCtx := trace.SpanContextFromContext(c.Request.Context())
		fields := logrus.Fields{
			"component":   "http",
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"status_code": c.Writer.Status(),
			"duration_ms": time.Since(start).Milliseconds(),
			"client_ip":   c.ClientIP(),
			"user_agent":  c.Request.UserAgent(),
		}

		// Add trace correlation if available
		if spanCtx.IsValid() {
			fields["trace_id"] = spanCtx.TraceID().String()
			fields["span_id"] = spanCtx.SpanID().String()
		}

		// Log request details
		logger.WithFields(fields).Info("HTTP request processed")
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

	// Create product service and analysis service
	productService := models.NewProductService(s.db)
	analysisService := services.NewAnalysisService(productService)
	productHandler := handlers.NewProductHandler(productService, analysisService)

	// Create frontend metrics handler
	frontendMetricsHandler := handlers.NewFrontendMetricsHandler()

	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		// Frontend metrics endpoint
		v1.POST("/frontend-metrics", frontendMetricsHandler.HandleFrontendMetrics)

		// Product routes
		products := v1.Group("/products")
		{
			products.GET("", productHandler.GetProducts)            // GET /api/v1/products
			products.POST("", productHandler.CreateProduct)         // POST /api/v1/products
			products.GET("/analyze", productHandler.AnalyzeProduct) // GET /api/v1/products/analyze
			products.GET("/:id", productHandler.GetProduct)         // GET /api/v1/products/:id
			products.PUT("/:id", productHandler.UpdateProduct)      // PUT /api/v1/products/:id
			products.DELETE("/:id", productHandler.DeleteProduct)   // DELETE /api/v1/products/:id
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

// Start starts both HTTP and gRPC servers
func (s *Server) Start(port string) error {
	httpPort := port
	grpcPort := "9090" // Fixed gRPC port
	
	// Start gRPC server in a goroutine
	go func() {
		logger.WithFields(logrus.Fields{
			"component": "server",
			"action":    "start_grpc",
			"port":      grpcPort,
		}).Info("Starting gRPC server")

		lis, err := net.Listen("tcp", ":"+grpcPort)
		if err != nil {
			logger.WithError(err).WithFields(logrus.Fields{
				"component": "server",
				"action":    "start_grpc",
				"port":      grpcPort,
			}).Fatal("Failed to listen on gRPC port")
		}

		if err := s.grpcServer.Serve(lis); err != nil {
			logger.WithError(err).WithFields(logrus.Fields{
				"component": "server",
				"action":    "start_grpc",
			}).Fatal("Failed to start gRPC server")
		}
	}()

	// Start HTTP server
	logger.WithFields(logrus.Fields{
		"component": "server",
		"action":    "start_http",
		"port":      httpPort,
	}).Info("Starting HTTP server")

	return s.router.Run(":" + httpPort)
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

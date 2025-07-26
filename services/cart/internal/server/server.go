package server

import (
	"context"
	"net/http"
	"time"

	"cart-service/internal/handlers"
	"cart-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// Server represents the HTTP server
type Server struct {
	router      *gin.Engine
	cartService *models.CartService
	logger      *logrus.Logger
}

// NewServer creates a new HTTP server
func NewServer(redisClient *redis.Client, logger *logrus.Logger) *Server {
	// Create cart service
	cartService := models.NewCartService(redisClient, logger)

	// Create server
	server := &Server{
		cartService: cartService,
		logger:      logger,
	}

	// Setup router
	server.setupRouter()

	return server
}

// setupRouter configures the Gin router with middleware and routes
func (s *Server) setupRouter() {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create router
	s.router = gin.New()

	// Add OpenTelemetry middleware
	s.router.Use(otelgin.Middleware("cart-service"))

	// Add custom middleware
	s.router.Use(s.loggingMiddleware())
	s.router.Use(s.corsMiddleware())
	s.router.Use(gin.Recovery())

	// Health check endpoint
	s.router.GET("/health", s.healthCheck)

	// Metrics endpoint (placeholder for now)
	s.router.GET("/metrics", s.metrics)

	// API routes
	s.setupAPIRoutes()
}

// setupAPIRoutes configures the API routes
func (s *Server) setupAPIRoutes() {
	// Create cart handler
	cartHandler := handlers.NewCartHandler(s.cartService, s.logger)

	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		cart := v1.Group("/cart")
		{
			cart.POST("/items", cartHandler.AddItem)                      // Add item to cart
			cart.GET("", cartHandler.GetCart)                             // Get cart contents
			cart.PUT("/items/:productId", cartHandler.UpdateItemQuantity) // Update item quantity
			cart.DELETE("/items/:productId", cartHandler.RemoveItem)      // Remove item from cart
			cart.DELETE("", cartHandler.ClearCart)                        // Clear entire cart
		}
	}
}

// loggingMiddleware adds request logging
func (s *Server) loggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		s.logger.WithFields(logrus.Fields{
			"component":   "http_server",
			"method":      param.Method,
			"path":        param.Path,
			"status":      param.StatusCode,
			"latency":     param.Latency,
			"client_ip":   param.ClientIP,
			"user_agent":  param.Request.UserAgent(),
		}).Info("HTTP request processed")

		return ""
	})
}

// corsMiddleware adds CORS headers
func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Session-ID")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// healthCheck handles health check requests
func (s *Server) healthCheck(c *gin.Context) {
	tracer := otel.Tracer("cart-service")
	ctx, span := tracer.Start(c.Request.Context(), "handler.health_check")
	defer span.End()

	// Check Redis connection
	redisClient := s.cartService
	if redisClient == nil {
		span.SetAttributes(attribute.String("health.status", "unhealthy"))
		span.SetAttributes(attribute.String("health.error", "redis_unavailable"))
		
		s.logger.WithFields(logrus.Fields{
			"component": "health_check",
			"status":    "unhealthy",
			"reason":    "redis_unavailable",
		}).Error("Health check failed: Redis unavailable")

		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"error":  "Redis unavailable",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	// Try to get an empty cart to test Redis connectivity
	_, err := s.cartService.GetCart(ctx, "health-check")
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("health.status", "unhealthy"))
		
		s.logger.WithError(err).WithFields(logrus.Fields{
			"component": "health_check",
			"status":    "unhealthy",
			"reason":    "redis_connection_failed",
		}).Error("Health check failed: Redis connection error")

		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"error":  "Redis connection failed",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	span.SetAttributes(attribute.String("health.status", "healthy"))

	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"service": "cart-service",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"version": "1.0.0",
	})
}

// metrics handles metrics requests (placeholder)
func (s *Server) metrics(c *gin.Context) {
	// TODO: Implement Prometheus metrics
	c.String(http.StatusOK, "# HELP cart_service_info Information about cart service\n# TYPE cart_service_info gauge\ncart_service_info{version=\"1.0.0\"} 1\n")
}

// Start starts the HTTP server
func (s *Server) Start(port string) error {
	s.logger.WithFields(logrus.Fields{
		"component": "http_server",
		"action":    "start",
		"port":      port,
	}).Info("Starting HTTP server")

	return s.router.Run(":" + port)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.WithFields(logrus.Fields{
		"component": "http_server",
		"action":    "shutdown",
	}).Info("Shutting down HTTP server")

	// TODO: Implement graceful shutdown
	return nil
}
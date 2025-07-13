package server

import (
	"database/sql"
	"log"

	"catalog-service/internal/handlers"
	"catalog-service/internal/models"

	"github.com/gin-gonic/gin"
)

// Server represents the HTTP server
type Server struct {
	router *gin.Engine
	db     *sql.DB
}

// NewServer creates a new server instance
func NewServer(database *sql.DB) *Server {
	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	server := &Server{
		router: router,
		db:     database,
	}

	// Setup routes
	server.setupRoutes()

	return server
}

// setupRoutes configures all the routes for the server
func (s *Server) setupRoutes() {
	// Health check endpoint
	healthHandler := handlers.NewHealthHandler(s.db)
	s.router.GET("/health", healthHandler.HealthCheck)

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
	log.Println("Registered routes:")
	for _, route := range s.router.Routes() {
		log.Printf("  %s %s", route.Method, route.Path)
	}
}

// Start starts the HTTP server
func (s *Server) Start(port string) error {
	log.Printf("Starting server on port %s", port)
	return s.router.Run(":" + port)
}

// Stop gracefully stops the server
func (s *Server) Stop() error {
	if s.db != nil {
		log.Println("Closing database connection...")
		return s.db.Close()
	}
	return nil
}

// GetDB returns the database connection (for testing purposes)
func (s *Server) GetDB() *sql.DB {
	return s.db
}

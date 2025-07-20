package models

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"catalog-service/internal/logger"

	"math/rand"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Product represents a product in the catalog
type Product struct {
	ID          int     `json:"id" db:"id"`
	Name        string  `json:"name" db:"name"`
	Description string  `json:"description" db:"description"`
	Price       float64 `json:"price" db:"price"`
	StockQty    int     `json:"stock_quantity" db:"stock_quantity"`
}

// ProductCreateRequest represents the request to create a new product
type ProductCreateRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	StockQty    int     `json:"stock_quantity" binding:"gte=0"`
}

// ProductUpdateRequest represents the request to update a product
type ProductUpdateRequest struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	StockQty    *int     `json:"stock_quantity,omitempty"`
}

// ProductResponse represents the response when returning a product
type ProductResponse struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	StockQty    int     `json:"stock_quantity"`
}

// AnalysisResult represents the response from product analysis
type AnalysisResult struct {
	ProductID       *int                   `json:"product_id,omitempty"`
	ComputeStats    ComputeStats           `json:"compute_stats"`
	DatabaseStats   DatabaseStats          `json:"database_stats"`
	ExternalData    ExternalData           `json:"external_data"`
	TotalDurationMs int64                  `json:"total_duration_ms"`
	Timestamp       string                 `json:"timestamp"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// ComputeStats represents computation analysis results
type ComputeStats struct {
	CalculationsPerformed int     `json:"calculations_performed"`
	ProcessingTimeMs      int64   `json:"processing_time_ms"`
	MemoryUsed            int64   `json:"memory_used_bytes"`
	ComplexityScore       float64 `json:"complexity_score"`
}

// DatabaseStats represents database performance metrics
type DatabaseStats struct {
	QueriesExecuted    int   `json:"queries_executed"`
	TotalQueryTimeMs   int64 `json:"total_query_time_ms"`
	AverageLatencyMs   int64 `json:"average_latency_ms"`
	SlowQueriesCount   int   `json:"slow_queries_count"`
	ConnectionPoolUsed int   `json:"connection_pool_used"`
}

// ExternalData represents data from external service calls
type ExternalData struct {
	ServiceCalled  string                 `json:"service_called"`
	ResponseTimeMs int64                  `json:"response_time_ms"`
	Success        bool                   `json:"success"`
	DataRetrieved  map[string]interface{} `json:"data_retrieved,omitempty"`
	ErrorMessage   string                 `json:"error_message,omitempty"`
}

// ProductService handles database operations for products
type ProductService struct {
	db *sql.DB
}

// NewProductService creates a new product service
func NewProductService(db *sql.DB) *ProductService {
	return &ProductService{db: db}
}

// CreateProduct creates a new product in the database
func (s *ProductService) CreateProduct(ctx *gin.Context, req ProductCreateRequest) (*Product, error) {
	// Start a database span
	tracer := otel.Tracer("catalog-service")
	dbCtx, span := tracer.Start(ctx.Request.Context(), "db.create_product")
	defer span.End()

	// Add span attributes
	span.SetAttributes(
		attribute.String("db.operation", "INSERT"),
		attribute.String("db.table", "products"),
		attribute.String("product.name", req.Name),
		attribute.Float64("product.price", req.Price),
	)

	query := `
		INSERT INTO products (name, description, price, stock_quantity)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, description, price, stock_quantity`

	var product Product
	err := s.db.QueryRowContext(dbCtx, query, req.Name, req.Description, req.Price, req.StockQty).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.StockQty,
	)
	if err != nil {
		span.RecordError(err)
		logger.WithError(err).WithFields(logrus.Fields{
			"component": "product",
			"action":    "create",
			"name":      req.Name,
			"price":     req.Price,
		}).Error("Error creating product")
		return nil, fmt.Errorf("failed to create product: %v", err)
	}

	span.SetAttributes(
		attribute.Int("product.id", product.ID),
	)

	logger.WithFields(logrus.Fields{
		"component":  "product",
		"action":     "create",
		"product_id": product.ID,
		"name":       product.Name,
		"price":      product.Price,
	}).Info("Created product")

	return &product, nil
}

// GetProduct retrieves a product by ID
func (s *ProductService) GetProduct(ctx *gin.Context, id int) (*Product, error) {
	// Start a database span
	tracer := otel.Tracer("catalog-service")
	dbCtx, span := tracer.Start(ctx.Request.Context(), "db.get_product")
	defer span.End()

	// Add span attributes
	span.SetAttributes(
		attribute.String("db.operation", "SELECT"),
		attribute.String("db.table", "products"),
		attribute.Int("product.id", id),
	)

	query := `SELECT id, name, description, price, stock_quantity FROM products WHERE id = $1`

	var product Product
	err := s.db.QueryRowContext(dbCtx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.StockQty,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			span.SetAttributes(attribute.String("db.result", "not_found"))
			return nil, fmt.Errorf("product not found")
		}
		span.RecordError(err)
		logger.WithError(err).WithFields(logrus.Fields{
			"component":  "product",
			"action":     "get",
			"product_id": id,
		}).Error("Error getting product")
		return nil, fmt.Errorf("failed to get product: %v", err)
	}

	span.SetAttributes(
		attribute.String("db.result", "found"),
		attribute.String("product.name", product.Name),
	)

	return &product, nil
}

// GetAllProducts retrieves all products with basic pagination
func (s *ProductService) GetAllProducts(ctx *gin.Context, offset, limit int) ([]Product, error) {
	// Start a database span
	tracer := otel.Tracer("catalog-service")
	dbCtx, span := tracer.Start(ctx.Request.Context(), "db.get_all_products")
	defer span.End()

	// Add span attributes
	span.SetAttributes(
		attribute.String("db.operation", "SELECT"),
		attribute.String("db.table", "products"),
		attribute.Int("query.offset", offset),
		attribute.Int("query.limit", limit),
	)

	query := `SELECT id, name, description, price, stock_quantity FROM products ORDER BY id LIMIT $1 OFFSET $2`

	rows, err := s.db.QueryContext(dbCtx, query, limit, offset)
	if err != nil {
		span.RecordError(err)
		logger.WithError(err).WithFields(logrus.Fields{
			"component": "product",
			"action":    "list",
			"offset":    offset,
			"limit":     limit,
		}).Error("Error getting products")
		return nil, fmt.Errorf("failed to get products: %v", err)
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.StockQty,
		)
		if err != nil {
			span.RecordError(err)
			logger.WithError(err).WithFields(logrus.Fields{
				"component": "product",
				"action":    "list",
				"operation": "scan",
			}).Error("Error scanning product row")
			return nil, fmt.Errorf("failed to scan product: %v", err)
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		span.RecordError(err)
		logger.WithError(err).WithFields(logrus.Fields{
			"component": "product",
			"action":    "list",
			"operation": "iterate",
		}).Error("Error iterating products")
		return nil, fmt.Errorf("failed to iterate products: %v", err)
	}

	span.SetAttributes(
		attribute.Int("products.count", len(products)),
	)

	return products, nil
}

// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(ctx *gin.Context, id int, req ProductUpdateRequest) (*Product, error) {
	// First, get the current product
	current, err := s.GetProduct(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update only the fields that were provided
	if req.Name != nil {
		current.Name = *req.Name
	}
	if req.Description != nil {
		current.Description = *req.Description
	}
	if req.Price != nil {
		current.Price = *req.Price
	}
	if req.StockQty != nil {
		current.StockQty = *req.StockQty
	}

	// Start a database span for the update
	tracer := otel.Tracer("catalog-service")
	dbCtx, span := tracer.Start(ctx.Request.Context(), "db.update_product")
	defer span.End()

	// Add span attributes
	span.SetAttributes(
		attribute.String("db.operation", "UPDATE"),
		attribute.String("db.table", "products"),
		attribute.Int("product.id", id),
	)

	// Update the database
	query := `
		UPDATE products 
		SET name = $1, description = $2, price = $3, stock_quantity = $4
		WHERE id = $5
		RETURNING id, name, description, price, stock_quantity`

	var product Product
	err = s.db.QueryRowContext(dbCtx, query, current.Name, current.Description, current.Price, current.StockQty, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.StockQty,
	)
	if err != nil {
		span.RecordError(err)
		logger.WithError(err).WithFields(logrus.Fields{
			"component":  "product",
			"action":     "update",
			"product_id": id,
		}).Error("Error updating product")
		return nil, fmt.Errorf("failed to update product: %v", err)
	}

	span.SetAttributes(
		attribute.String("product.name", product.Name),
		attribute.Float64("product.price", product.Price),
	)

	logger.WithFields(logrus.Fields{
		"component":  "product",
		"action":     "update",
		"product_id": product.ID,
		"name":       product.Name,
		"price":      product.Price,
	}).Info("Updated product")

	return &product, nil
}

// DeleteProduct deletes a product by ID
func (s *ProductService) DeleteProduct(ctx *gin.Context, id int) error {
	// Start a database span
	tracer := otel.Tracer("catalog-service")
	dbCtx, span := tracer.Start(ctx.Request.Context(), "db.delete_product")
	defer span.End()

	// Add span attributes
	span.SetAttributes(
		attribute.String("db.operation", "DELETE"),
		attribute.String("db.table", "products"),
		attribute.Int("product.id", id),
	)

	query := `DELETE FROM products WHERE id = $1`

	result, err := s.db.ExecContext(dbCtx, query, id)
	if err != nil {
		span.RecordError(err)
		logger.WithError(err).WithFields(logrus.Fields{
			"component":  "product",
			"action":     "delete",
			"product_id": id,
		}).Error("Error deleting product")
		return fmt.Errorf("failed to delete product: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		span.RecordError(err)
		logger.WithError(err).WithFields(logrus.Fields{
			"component":  "product",
			"action":     "delete",
			"product_id": id,
			"operation":  "rows_affected",
		}).Error("Error getting rows affected for product")
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		span.SetAttributes(attribute.String("db.result", "not_found"))
		return fmt.Errorf("product not found")
	}

	span.SetAttributes(
		attribute.Int64("db.rows_affected", rowsAffected),
		attribute.String("db.result", "deleted"),
	)

	logger.WithFields(logrus.Fields{
		"component":  "product",
		"action":     "delete",
		"product_id": id,
	}).Info("Deleted product")

	return nil
}

// ToResponse converts a Product to a ProductResponse
func (p *Product) ToResponse() ProductResponse {
	return ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		StockQty:    p.StockQty,
	}
}

// AnalyzeProduct performs complex analysis with multiple spans for tracing demonstration
func (s *ProductService) AnalyzeProduct(ctx *gin.Context, productID *int) (*AnalysisResult, error) {
	// Start the main analysis span
	tracer := otel.Tracer("catalog-service")
	mainCtx, mainSpan := tracer.Start(ctx.Request.Context(), "product.analyze")
	defer mainSpan.End()

	mainSpan.SetAttributes(
		attribute.String("analysis.type", "comprehensive"),
		attribute.Bool("analysis.has_product_id", productID != nil),
	)
	if productID != nil {
		mainSpan.SetAttributes(attribute.Int("product.id", *productID))
	}

	startTime := time.Now()

	// Perform computation analysis
	computeStats, err := s.performComputeAnalysis(mainCtx, tracer)
	if err != nil {
		mainSpan.RecordError(err)
		return nil, fmt.Errorf("compute analysis failed: %w", err)
	}

	// Perform database analysis with multiple operations
	dbStats, err := s.performDatabaseAnalysis(mainCtx, tracer, productID)
	if err != nil {
		mainSpan.RecordError(err)
		return nil, fmt.Errorf("database analysis failed: %w", err)
	}

	// Perform external service call
	externalData, err := s.performExternalAnalysis(mainCtx, tracer)
	if err != nil {
		// Don't fail the whole request if external service fails
		logger.WithError(err).Warn("External analysis failed, continuing with empty data")
		externalData = &ExternalData{
			ServiceCalled:  "httpbin.org",
			Success:        false,
			ErrorMessage:   err.Error(),
			ResponseTimeMs: 0,
		}
	}

	duration := time.Since(startTime)
	mainSpan.SetAttributes(attribute.Int64("analysis.total_duration_ms", duration.Milliseconds()))

	result := &AnalysisResult{
		ProductID:       productID,
		ComputeStats:    *computeStats,
		DatabaseStats:   *dbStats,
		ExternalData:    *externalData,
		TotalDurationMs: duration.Milliseconds(),
		Timestamp:       time.Now().Format(time.RFC3339),
		Metadata: map[string]interface{}{
			"version":   "1.0.0",
			"algorithm": "comprehensive-v1",
			"cluster":   "local-k8s",
			"trace_id":  mainSpan.SpanContext().TraceID().String(),
		},
	}

	logger.WithFields(logrus.Fields{
		"component":      "product",
		"action":         "analyze",
		"product_id":     productID,
		"duration_ms":    duration.Milliseconds(),
		"compute_ops":    computeStats.CalculationsPerformed,
		"db_queries":     dbStats.QueriesExecuted,
		"external_calls": 1,
	}).Info("Product analysis completed")

	return result, nil
}

// performComputeAnalysis simulates CPU-intensive computation with tracing
func (s *ProductService) performComputeAnalysis(ctx context.Context, tracer trace.Tracer) (*ComputeStats, error) {
	computeCtx, span := tracer.Start(ctx, "compute.analysis")
	defer span.End()

	span.SetAttributes(attribute.String("compute.type", "mathematical"))

	startTime := time.Now()

	// Simulate complex computation with multiple operations
	calculations := 0
	var results []float64

	// Phase 1: Matrix operations simulation
	_, matrixSpan := tracer.Start(computeCtx, "compute.matrix_operations")
	matrixSpan.SetAttributes(attribute.String("compute.phase", "matrix"))

	for i := 0; i < 1000; i++ {
		// Simulate matrix multiplication
		result := float64(i*i) * 1.414213562373095 // sqrt(2)
		results = append(results, result)
		calculations++
	}
	time.Sleep(50 * time.Millisecond) // Simulate processing time
	matrixSpan.End()

	// Phase 2: Statistical analysis simulation
	_, statsSpan := tracer.Start(computeCtx, "compute.statistical_analysis")
	statsSpan.SetAttributes(attribute.String("compute.phase", "statistics"))

	var sum, mean, variance float64
	for _, val := range results {
		sum += val
		calculations++
	}
	mean = sum / float64(len(results))

	for _, val := range results {
		variance += (val - mean) * (val - mean)
		calculations++
	}
	variance /= float64(len(results))

	time.Sleep(30 * time.Millisecond) // Simulate processing time
	statsSpan.SetAttributes(
		attribute.Float64("compute.mean", mean),
		attribute.Float64("compute.variance", variance),
	)
	statsSpan.End()

	// Phase 3: Complexity scoring
	_, complexSpan := tracer.Start(computeCtx, "compute.complexity_scoring")
	complexSpan.SetAttributes(attribute.String("compute.phase", "complexity"))

	complexityScore := variance / (mean + 1) * 100 // Arbitrary complexity metric
	time.Sleep(20 * time.Millisecond)
	complexSpan.SetAttributes(attribute.Float64("compute.complexity_score", complexityScore))
	complexSpan.End()

	duration := time.Since(startTime)
	memoryUsed := int64(len(results) * 8) // Approximate bytes for float64 slice

	span.SetAttributes(
		attribute.Int("compute.calculations", calculations),
		attribute.Int64("compute.duration_ms", duration.Milliseconds()),
		attribute.Int64("compute.memory_bytes", memoryUsed),
		attribute.Float64("compute.final_score", complexityScore),
	)

	return &ComputeStats{
		CalculationsPerformed: calculations,
		ProcessingTimeMs:      duration.Milliseconds(),
		MemoryUsed:            memoryUsed,
		ComplexityScore:       complexityScore,
	}, nil
}

// performDatabaseAnalysis executes simple database operations to demonstrate DB spans
func (s *ProductService) performDatabaseAnalysis(ctx context.Context, tracer trace.Tracer, productID *int) (*DatabaseStats, error) {
	dbCtx, span := tracer.Start(ctx, "database.analysis")
	defer span.End()

	span.SetAttributes(attribute.String("db.analysis_type", "simple_demo"))

	startTime := time.Now()
	queriesExecuted := 0

	// Query 1: Count all products (demonstrates basic SELECT span)
	countCtx, countSpan := tracer.Start(dbCtx, "db.count_products")
	countSpan.SetAttributes(
		attribute.String("db.operation", "COUNT"),
		attribute.String("db.table", "products"),
	)

	// Simple fixed delay to show span duration
	time.Sleep(80 * time.Millisecond)

	var count int
	err := s.db.QueryRowContext(countCtx, "SELECT COUNT(*) FROM products").Scan(&count)
	if err != nil {
		countSpan.RecordError(err)
		return nil, fmt.Errorf("count query failed: %w", err)
	}

	queriesExecuted++
	countSpan.SetAttributes(
		attribute.Int("db.result_count", count),
		attribute.Int64("db.duration_ms", 80),
	)
	countSpan.End()

	// Query 2: Optional product lookup (demonstrates conditional spans)
	if productID != nil {
		lookupCtx, lookupSpan := tracer.Start(dbCtx, "db.product_lookup")
		lookupSpan.SetAttributes(
			attribute.String("db.operation", "SELECT"),
			attribute.String("db.table", "products"),
			attribute.Int("db.product_id", *productID),
		)

		// Slightly longer delay for specific product lookup
		time.Sleep(120 * time.Millisecond)

		var product Product
		err = s.db.QueryRowContext(lookupCtx, "SELECT id, name, price, stock_quantity FROM products WHERE id = $1", *productID).Scan(
			&product.ID, &product.Name, &product.Price, &product.StockQty,
		)

		queriesExecuted++

		if err == sql.ErrNoRows {
			lookupSpan.SetAttributes(attribute.String("db.result", "not_found"))
		} else if err != nil {
			lookupSpan.RecordError(err)
			return nil, fmt.Errorf("product lookup failed: %w", err)
		} else {
			lookupSpan.SetAttributes(
				attribute.String("db.result", "found"),
				attribute.String("db.product_name", product.Name),
			)
		}

		lookupSpan.SetAttributes(attribute.Int64("db.duration_ms", 120))
		lookupSpan.End()
	}

	totalTime := time.Since(startTime)
	averageLatency := totalTime.Milliseconds() / int64(queriesExecuted)

	span.SetAttributes(
		attribute.Int("db.queries_executed", queriesExecuted),
		attribute.Int64("db.total_time_ms", totalTime.Milliseconds()),
		attribute.Int64("db.average_latency_ms", averageLatency),
	)

	return &DatabaseStats{
		QueriesExecuted:    queriesExecuted,
		TotalQueryTimeMs:   totalTime.Milliseconds(),
		AverageLatencyMs:   averageLatency,
		SlowQueriesCount:   0, // Simplified - no slow query detection
		ConnectionPoolUsed: 2, // Fixed value for simplicity
	}, nil
}

// performExternalAnalysis makes an external HTTP call with tracing
func (s *ProductService) performExternalAnalysis(ctx context.Context, tracer trace.Tracer) (*ExternalData, error) {
	externalCtx, span := tracer.Start(ctx, "external.api_call")
	defer span.End()

	span.SetAttributes(
		attribute.String("http.method", "GET"),
		attribute.String("http.url", "https://httpbin.org/delay/1"),
		attribute.String("external.service", "httpbin"),
	)

	startTime := time.Now()

	// Make HTTP request to httpbin for demonstration
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequestWithContext(externalCtx, "GET", "https://httpbin.org/delay/1", nil)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("external call failed: %w", err)
	}
	defer resp.Body.Close()

	responseTime := time.Since(startTime)

	span.SetAttributes(
		attribute.Int("http.status_code", resp.StatusCode),
		attribute.Int64("http.response_time_ms", responseTime.Milliseconds()),
		attribute.String("http.response_content_type", resp.Header.Get("Content-Type")),
	)

	// Simulate parsing response data
	data := map[string]interface{}{
		"status_code":    resp.StatusCode,
		"content_type":   resp.Header.Get("Content-Type"),
		"response_size":  resp.ContentLength,
		"server":         resp.Header.Get("Server"),
		"simulated_data": "This is fake external data for demo purposes",
		"random_value":   rand.Float64() * 100,
	}

	success := resp.StatusCode == 200

	return &ExternalData{
		ServiceCalled:  "httpbin.org",
		ResponseTimeMs: responseTime.Milliseconds(),
		Success:        success,
		DataRetrieved:  data,
	}, nil
}

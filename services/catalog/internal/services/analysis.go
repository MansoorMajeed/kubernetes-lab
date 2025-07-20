package services

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"catalog-service/internal/logger"
	"catalog-service/internal/models"

	"math/rand"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

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

// AnalysisService handles complex analysis operations with distributed tracing
type AnalysisService struct {
	productService *models.ProductService
}

// NewAnalysisService creates a new analysis service
func NewAnalysisService(productService *models.ProductService) *AnalysisService {
	return &AnalysisService{
		productService: productService,
	}
}

/* This whole analysis is just so that you can see cooler traces. They add no other value.
 */
// AnalyzeProduct performs complex analysis with multiple spans for tracing demonstration
func (s *AnalysisService) AnalyzeProduct(ctx *gin.Context, productID *int) (*AnalysisResult, error) {
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
		"component":      "analysis",
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
func (s *AnalysisService) performComputeAnalysis(ctx context.Context, tracer trace.Tracer) (*ComputeStats, error) {
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
func (s *AnalysisService) performDatabaseAnalysis(ctx context.Context, tracer trace.Tracer, productID *int) (*DatabaseStats, error) {
	dbCtx, span := tracer.Start(ctx, "database.analysis")
	defer span.End()

	span.SetAttributes(attribute.String("db.analysis_type", "simple_demo"))

	startTime := time.Now()
	queriesExecuted := 0

	// Query 1: Count all products (demonstrates basic SELECT span)
	count, err := s.productService.GetProductCount(dbCtx)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("count query failed: %w", err)
	}
	queriesExecuted++

	// Query 2: Optional product lookup (demonstrates conditional spans)
	if productID != nil {
		_, err := s.productService.GetProductByID(dbCtx, *productID)
		if err != nil && err.Error() != "product not found" {
			span.RecordError(err)
			return nil, fmt.Errorf("product lookup failed: %w", err)
		}
		queriesExecuted++
	}

	totalTime := time.Since(startTime)
	averageLatency := totalTime.Milliseconds() / int64(queriesExecuted)

	span.SetAttributes(
		attribute.Int("db.queries_executed", queriesExecuted),
		attribute.Int64("db.total_time_ms", totalTime.Milliseconds()),
		attribute.Int64("db.average_latency_ms", averageLatency),
		attribute.Int("db.product_count", count),
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
func (s *AnalysisService) performExternalAnalysis(ctx context.Context, tracer trace.Tracer) (*ExternalData, error) {
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

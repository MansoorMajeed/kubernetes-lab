package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"catalog-service/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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

// GetProductCount returns the total number of products (helper for analysis)
func (s *ProductService) GetProductCount(ctx context.Context) (int, error) {
	tracer := otel.Tracer("catalog-service")
	dbCtx, span := tracer.Start(ctx, "db.count_products")
	defer span.End()

	span.SetAttributes(
		attribute.String("db.operation", "COUNT"),
		attribute.String("db.table", "products"),
	)

	// Simple fixed delay to demonstrate span duration
	time.Sleep(80 * time.Millisecond)

	var count int
	err := s.db.QueryRowContext(dbCtx, "SELECT COUNT(*) FROM products").Scan(&count)
	if err != nil {
		span.RecordError(err)
		return 0, err
	}

	span.SetAttributes(
		attribute.Int("db.result_count", count),
		attribute.Int64("db.duration_ms", 80),
	)

	return count, nil
}

// GetProductByID returns a product by ID (helper for analysis)
func (s *ProductService) GetProductByID(ctx context.Context, id int) (*Product, error) {
	tracer := otel.Tracer("catalog-service")
	dbCtx, span := tracer.Start(ctx, "db.product_lookup")
	defer span.End()

	span.SetAttributes(
		attribute.String("db.operation", "SELECT"),
		attribute.String("db.table", "products"),
		attribute.Int("db.product_id", id),
	)

	// Slightly longer delay for specific lookup
	time.Sleep(120 * time.Millisecond)

	query := `SELECT id, name, description, price, stock_quantity FROM products WHERE id = $1`
	var product Product

	err := s.db.QueryRowContext(dbCtx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.StockQty,
	)

	if err == sql.ErrNoRows {
		span.SetAttributes(attribute.String("db.result", "not_found"))
		return nil, fmt.Errorf("product not found")
	} else if err != nil {
		span.RecordError(err)
		return nil, err
	}

	span.SetAttributes(
		attribute.String("db.result", "found"),
		attribute.String("db.product_name", product.Name),
		attribute.Int64("db.duration_ms", 120),
	)

	return &product, nil
}

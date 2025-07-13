package models

import (
	"database/sql"
	"fmt"
	"log"
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
func (s *ProductService) CreateProduct(req ProductCreateRequest) (*Product, error) {
	query := `
		INSERT INTO products (name, description, price, stock_quantity)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, description, price, stock_quantity`

	var product Product
	err := s.db.QueryRow(query, req.Name, req.Description, req.Price, req.StockQty).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.StockQty,
	)
	if err != nil {
		log.Printf("Error creating product: %v", err)
		return nil, fmt.Errorf("failed to create product: %v", err)
	}

	log.Printf("Created product: %+v", product)
	return &product, nil
}

// GetProduct retrieves a product by ID
func (s *ProductService) GetProduct(id int) (*Product, error) {
	query := `SELECT id, name, description, price, stock_quantity FROM products WHERE id = $1`

	var product Product
	err := s.db.QueryRow(query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.StockQty,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		log.Printf("Error getting product %d: %v", id, err)
		return nil, fmt.Errorf("failed to get product: %v", err)
	}

	return &product, nil
}

// GetAllProducts retrieves all products with basic pagination
func (s *ProductService) GetAllProducts(offset, limit int) ([]Product, error) {
	query := `SELECT id, name, description, price, stock_quantity FROM products ORDER BY id LIMIT $1 OFFSET $2`

	rows, err := s.db.Query(query, limit, offset)
	if err != nil {
		log.Printf("Error getting products: %v", err)
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
			log.Printf("Error scanning product row: %v", err)
			return nil, fmt.Errorf("failed to scan product: %v", err)
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating products: %v", err)
		return nil, fmt.Errorf("failed to iterate products: %v", err)
	}

	return products, nil
}

// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(id int, req ProductUpdateRequest) (*Product, error) {
	// First, get the current product
	current, err := s.GetProduct(id)
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

	// Update the database
	query := `
		UPDATE products 
		SET name = $1, description = $2, price = $3, stock_quantity = $4
		WHERE id = $5
		RETURNING id, name, description, price, stock_quantity`

	var product Product
	err = s.db.QueryRow(query, current.Name, current.Description, current.Price, current.StockQty, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.StockQty,
	)
	if err != nil {
		log.Printf("Error updating product %d: %v", id, err)
		return nil, fmt.Errorf("failed to update product: %v", err)
	}

	log.Printf("Updated product: %+v", product)
	return &product, nil
}

// DeleteProduct deletes a product by ID
func (s *ProductService) DeleteProduct(id int) error {
	query := `DELETE FROM products WHERE id = $1`

	result, err := s.db.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting product %d: %v", id, err)
		return fmt.Errorf("failed to delete product: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected for product %d: %v", id, err)
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}

	log.Printf("Deleted product %d", id)
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

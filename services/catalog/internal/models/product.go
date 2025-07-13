package models

import (
	"database/sql/driver"
	"fmt"
)

// Product represents a product in the catalog
type Product struct {
	ID            int     `json:"id" db:"id"`
	Name          string  `json:"name" db:"name"`
	Description   string  `json:"description" db:"description"`
	Price         float64 `json:"price" db:"price"`
	StockQuantity int     `json:"stock_quantity" db:"stock_quantity"`
}

// CreateProductRequest represents the request payload for creating a product
type CreateProductRequest struct {
	Name          string  `json:"name" binding:"required"`
	Description   string  `json:"description"`
	Price         float64 `json:"price" binding:"required,gt=0"`
	StockQuantity int     `json:"stock_quantity" binding:"gte=0"`
}

// UpdateProductRequest represents the request payload for updating a product
type UpdateProductRequest struct {
	Name          *string  `json:"name,omitempty"`
	Description   *string  `json:"description,omitempty"`
	Price         *float64 `json:"price,omitempty" binding:"omitempty,gt=0"`
	StockQuantity *int     `json:"stock_quantity,omitempty" binding:"omitempty,gte=0"`
}

// ProductListResponse represents the response for listing products
type ProductListResponse struct {
	Products []Product `json:"products"`
	Count    int       `json:"count"`
}

// Validate checks if the product data is valid
func (p *Product) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("product name is required")
	}
	if p.Price <= 0 {
		return fmt.Errorf("product price must be greater than 0")
	}
	if p.StockQuantity < 0 {
		return fmt.Errorf("stock quantity cannot be negative")
	}
	return nil
}

// Value implements the driver.Valuer interface for database operations
func (p Product) Value() (driver.Value, error) {
	return fmt.Sprintf("(%d,%s,%s,%.2f,%d)",
		p.ID, p.Name, p.Description, p.Price, p.StockQuantity), nil
}

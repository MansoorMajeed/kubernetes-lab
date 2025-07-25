package models

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// CartItem represents a single item in a shopping cart
type CartItem struct {
	ProductID int     `json:"product_id" redis:"product_id"`
	Quantity  int     `json:"quantity" redis:"quantity"`
	Price     float64 `json:"price" redis:"price"`
	Name      string  `json:"name" redis:"name"`
	AddedAt   string  `json:"added_at" redis:"added_at"`
}

// Cart represents a user's shopping cart
type Cart struct {
	SessionID string     `json:"session_id" redis:"session_id"`
	Items     []CartItem `json:"items" redis:"items"`
	Total     float64    `json:"total" redis:"total"`
	CreatedAt string     `json:"created_at" redis:"created_at"`
	UpdatedAt string     `json:"updated_at" redis:"updated_at"`
}

// CartAddItemRequest represents the request to add an item to cart
type CartAddItemRequest struct {
	ProductID int `json:"product_id" binding:"required,gt=0"`
	Quantity  int `json:"quantity" binding:"required,gt=0"`
}

// CartUpdateItemRequest represents the request to update an item quantity
type CartUpdateItemRequest struct {
	Quantity int `json:"quantity" binding:"required,gte=0"`
}

// CartResponse represents the response when returning a cart
type CartResponse struct {
	SessionID string     `json:"session_id"`
	Items     []CartItem `json:"items"`
	Total     float64    `json:"total"`
	ItemCount int        `json:"item_count"`
	CreatedAt string     `json:"created_at"`
	UpdatedAt string     `json:"updated_at"`
}

// CartService handles Redis operations for carts
type CartService struct {
	redis  *redis.Client
	logger *logrus.Logger
}

// NewCartService creates a new cart service
func NewCartService(redisClient *redis.Client, logger *logrus.Logger) *CartService {
	return &CartService{
		redis:  redisClient,
		logger: logger,
	}
}

// getCartKey generates Redis key for cart storage
func (s *CartService) getCartKey(sessionID string) string {
	return fmt.Sprintf("cart:%s", sessionID)
}

// GetCart retrieves a cart by session ID
func (s *CartService) GetCart(ctx context.Context, sessionID string) (*Cart, error) {
	tracer := otel.Tracer("cart-service")
	redisCtx, span := tracer.Start(ctx, "redis.get_cart")
	defer span.End()

	span.SetAttributes(
		attribute.String("redis.operation", "GET"),
		attribute.String("cart.session_id", sessionID),
	)

	key := s.getCartKey(sessionID)
	data, err := s.redis.Get(redisCtx, key).Result()
	if err == redis.Nil {
		// Cart doesn't exist, return empty cart
		span.SetAttributes(attribute.String("redis.result", "not_found"))
		return s.createEmptyCart(sessionID), nil
	}
	if err != nil {
		span.RecordError(err)
		s.logger.WithError(err).WithFields(logrus.Fields{
			"component":  "cart",
			"action":     "get",
			"session_id": sessionID,
		}).Error("Error getting cart from Redis")
		return nil, fmt.Errorf("failed to get cart: %v", err)
	}

	var cart Cart
	if err := json.Unmarshal([]byte(data), &cart); err != nil {
		span.RecordError(err)
		s.logger.WithError(err).WithFields(logrus.Fields{
			"component":  "cart",
			"action":     "get",
			"session_id": sessionID,
			"operation":  "unmarshal",
		}).Error("Error unmarshaling cart data")
		return nil, fmt.Errorf("failed to unmarshal cart: %v", err)
	}

	span.SetAttributes(
		attribute.String("redis.result", "found"),
		attribute.Int("cart.item_count", len(cart.Items)),
		attribute.Float64("cart.total", cart.Total),
	)

	return &cart, nil
}

// SaveCart saves a cart to Redis with TTL
func (s *CartService) SaveCart(ctx context.Context, cart *Cart) error {
	tracer := otel.Tracer("cart-service")
	redisCtx, span := tracer.Start(ctx, "redis.save_cart")
	defer span.End()

	span.SetAttributes(
		attribute.String("redis.operation", "SET"),
		attribute.String("cart.session_id", cart.SessionID),
		attribute.Int("cart.item_count", len(cart.Items)),
		attribute.Float64("cart.total", cart.Total),
	)

	// Update timestamps
	now := time.Now().UTC().Format(time.RFC3339)
	cart.UpdatedAt = now
	if cart.CreatedAt == "" {
		cart.CreatedAt = now
	}

	// Recalculate total
	cart.calculateTotal()

	// Marshal cart to JSON
	data, err := json.Marshal(cart)
	if err != nil {
		span.RecordError(err)
		s.logger.WithError(err).WithFields(logrus.Fields{
			"component":  "cart",
			"action":     "save",
			"session_id": cart.SessionID,
			"operation":  "marshal",
		}).Error("Error marshaling cart data")
		return fmt.Errorf("failed to marshal cart: %v", err)
	}

	// Save to Redis with 24 hour TTL
	key := s.getCartKey(cart.SessionID)
	ttl := 24 * time.Hour
	if err := s.redis.Set(redisCtx, key, data, ttl).Err(); err != nil {
		span.RecordError(err)
		s.logger.WithError(err).WithFields(logrus.Fields{
			"component":  "cart",
			"action":     "save",
			"session_id": cart.SessionID,
			"ttl_hours":  24,
		}).Error("Error saving cart to Redis")
		return fmt.Errorf("failed to save cart: %v", err)
	}

	span.SetAttributes(
		attribute.String("redis.result", "saved"),
		attribute.Int64("redis.ttl_seconds", int64(ttl.Seconds())),
	)

	s.logger.WithFields(logrus.Fields{
		"component":  "cart",
		"action":     "save",
		"session_id": cart.SessionID,
		"item_count": len(cart.Items),
		"total":      cart.Total,
	}).Info("Saved cart to Redis")

	return nil
}

// DeleteCart removes a cart from Redis
func (s *CartService) DeleteCart(ctx context.Context, sessionID string) error {
	tracer := otel.Tracer("cart-service")
	redisCtx, span := tracer.Start(ctx, "redis.delete_cart")
	defer span.End()

	span.SetAttributes(
		attribute.String("redis.operation", "DEL"),
		attribute.String("cart.session_id", sessionID),
	)

	key := s.getCartKey(sessionID)
	result, err := s.redis.Del(redisCtx, key).Result()
	if err != nil {
		span.RecordError(err)
		s.logger.WithError(err).WithFields(logrus.Fields{
			"component":  "cart",
			"action":     "delete",
			"session_id": sessionID,
		}).Error("Error deleting cart from Redis")
		return fmt.Errorf("failed to delete cart: %v", err)
	}

	span.SetAttributes(
		attribute.Int64("redis.keys_deleted", result),
		attribute.String("redis.result", "deleted"),
	)

	s.logger.WithFields(logrus.Fields{
		"component":     "cart",
		"action":        "delete",
		"session_id":    sessionID,
		"keys_deleted":  result,
	}).Info("Deleted cart from Redis")

	return nil
}

// AddItem adds or updates an item in the cart
func (s *CartService) AddItem(ctx context.Context, sessionID string, productID int, quantity int, price float64, name string) (*Cart, error) {
	// Get existing cart
	cart, err := s.GetCart(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Check if item already exists
	found := false
	for i, item := range cart.Items {
		if item.ProductID == productID {
			// Update existing item quantity
			cart.Items[i].Quantity += quantity
			found = true
			break
		}
	}

	// Add new item if not found
	if !found {
		newItem := CartItem{
			ProductID: productID,
			Quantity:  quantity,
			Price:     price,
			Name:      name,
			AddedAt:   time.Now().UTC().Format(time.RFC3339),
		}
		cart.Items = append(cart.Items, newItem)
	}

	// Save cart
	if err := s.SaveCart(ctx, cart); err != nil {
		return nil, err
	}

	return cart, nil
}

// UpdateItemQuantity updates the quantity of an item in the cart
func (s *CartService) UpdateItemQuantity(ctx context.Context, sessionID string, productID int, quantity int) (*Cart, error) {
	// Get existing cart
	cart, err := s.GetCart(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Find and update item
	found := false
	for i, item := range cart.Items {
		if item.ProductID == productID {
			if quantity == 0 {
				// Remove item if quantity is 0
				cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			} else {
				cart.Items[i].Quantity = quantity
			}
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("item not found in cart")
	}

	// Save cart
	if err := s.SaveCart(ctx, cart); err != nil {
		return nil, err
	}

	return cart, nil
}

// RemoveItem removes an item from the cart
func (s *CartService) RemoveItem(ctx context.Context, sessionID string, productID int) (*Cart, error) {
	// Get existing cart
	cart, err := s.GetCart(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Find and remove item
	found := false
	for i, item := range cart.Items {
		if item.ProductID == productID {
			cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("item not found in cart")
	}

	// Save cart
	if err := s.SaveCart(ctx, cart); err != nil {
		return nil, err
	}

	return cart, nil
}

// createEmptyCart creates a new empty cart
func (s *CartService) createEmptyCart(sessionID string) *Cart {
	now := time.Now().UTC().Format(time.RFC3339)
	return &Cart{
		SessionID: sessionID,
		Items:     []CartItem{},
		Total:     0,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// calculateTotal recalculates the cart total
func (c *Cart) calculateTotal() {
	total := 0.0
	for _, item := range c.Items {
		total += item.Price * float64(item.Quantity)
	}
	c.Total = total
}

// ToResponse converts a Cart to a CartResponse
func (c *Cart) ToResponse() CartResponse {
	return CartResponse{
		SessionID: c.SessionID,
		Items:     c.Items,
		Total:     c.Total,
		ItemCount: len(c.Items),
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// GetItemCount returns the total number of items in the cart
func (c *Cart) GetItemCount() int {
	count := 0
	for _, item := range c.Items {
		count += item.Quantity
	}
	return count
}
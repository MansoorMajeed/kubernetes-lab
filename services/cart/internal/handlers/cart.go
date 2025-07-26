package handlers

import (
	"net/http"
	"strconv"

	"cart-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// CartHandler handles HTTP requests for cart operations
type CartHandler struct {
	cartService *models.CartService
	logger      *logrus.Logger
}

// NewCartHandler creates a new cart handler
func NewCartHandler(cartService *models.CartService, logger *logrus.Logger) *CartHandler {
	return &CartHandler{
		cartService: cartService,
		logger:      logger,
	}
}

// getSessionID extracts session ID from request headers or generates one
func (h *CartHandler) getSessionID(c *gin.Context) string {
	// Try to get session ID from header first
	sessionID := c.GetHeader("X-Session-ID")
	if sessionID != "" {
		return sessionID
	}
	
	// Try to get from query parameter
	sessionID = c.Query("session_id")
	if sessionID != "" {
		return sessionID
	}
	
	// For now, use a default session (in production, this would be from JWT or cookies)
	return "default-session"
}

// AddItem adds an item to the cart
// POST /api/v1/cart/items
func (h *CartHandler) AddItem(c *gin.Context) {
	tracer := otel.Tracer("cart-service")
	ctx, span := tracer.Start(c.Request.Context(), "handler.add_item")
	defer span.End()

	sessionID := h.getSessionID(c)
	span.SetAttributes(attribute.String("cart.session_id", sessionID))

	var req models.CartAddItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		span.RecordError(err)
		h.logger.WithError(err).WithFields(logrus.Fields{
			"component":  "cart_handler",
			"action":     "add_item",
			"session_id": sessionID,
		}).Error("Invalid request body")
		
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	span.SetAttributes(
		attribute.Int("item.product_id", req.ProductID),
		attribute.Int("item.quantity", req.Quantity),
	)

	h.logger.WithFields(logrus.Fields{
		"component":  "cart_handler",
		"action":     "add_item",
		"session_id": sessionID,
		"product_id": req.ProductID,
		"quantity":   req.Quantity,
	}).Info("Adding item to cart")

	// For now, use placeholder price and name (Phase 2 will add catalog validation)
	price := 29.99 // TODO: Get from catalog service
	name := "Product Name" // TODO: Get from catalog service

	cart, err := h.cartService.AddItem(ctx, sessionID, req.ProductID, req.Quantity, price, name)
	if err != nil {
		span.RecordError(err)
		h.logger.WithError(err).WithFields(logrus.Fields{
			"component":  "cart_handler",
			"action":     "add_item",
			"session_id": sessionID,
			"product_id": req.ProductID,
		}).Error("Failed to add item to cart")
		
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add item to cart",
		})
		return
	}

	span.SetAttributes(
		attribute.Int("cart.item_count", len(cart.Items)),
		attribute.Float64("cart.total", cart.Total),
	)

	h.logger.WithFields(logrus.Fields{
		"component":  "cart_handler",
		"action":     "add_item",
		"session_id": sessionID,
		"product_id": req.ProductID,
		"item_count": len(cart.Items),
		"total":      cart.Total,
	}).Info("Item added to cart successfully")

	c.JSON(http.StatusOK, gin.H{
		"message": "Item added to cart",
		"cart":    cart.ToResponse(),
	})
}

// GetCart retrieves the current cart
// GET /api/v1/cart
func (h *CartHandler) GetCart(c *gin.Context) {
	tracer := otel.Tracer("cart-service")
	ctx, span := tracer.Start(c.Request.Context(), "handler.get_cart")
	defer span.End()

	sessionID := h.getSessionID(c)
	span.SetAttributes(attribute.String("cart.session_id", sessionID))

	h.logger.WithFields(logrus.Fields{
		"component":  "cart_handler",
		"action":     "get_cart",
		"session_id": sessionID,
	}).Info("Getting cart")

	cart, err := h.cartService.GetCart(ctx, sessionID)
	if err != nil {
		span.RecordError(err)
		h.logger.WithError(err).WithFields(logrus.Fields{
			"component":  "cart_handler",
			"action":     "get_cart",
			"session_id": sessionID,
		}).Error("Failed to get cart")
		
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get cart",
		})
		return
	}

	span.SetAttributes(
		attribute.Int("cart.item_count", len(cart.Items)),
		attribute.Float64("cart.total", cart.Total),
	)

	c.JSON(http.StatusOK, cart.ToResponse())
}

// UpdateItemQuantity updates the quantity of an item in the cart
// PUT /api/v1/cart/items/:productId
func (h *CartHandler) UpdateItemQuantity(c *gin.Context) {
	tracer := otel.Tracer("cart-service")
	ctx, span := tracer.Start(c.Request.Context(), "handler.update_item_quantity")
	defer span.End()

	sessionID := h.getSessionID(c)
	span.SetAttributes(attribute.String("cart.session_id", sessionID))

	// Get product ID from URL parameter
	productIDStr := c.Param("productId")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		span.RecordError(err)
		h.logger.WithError(err).WithFields(logrus.Fields{
			"component":   "cart_handler",
			"action":      "update_item_quantity",
			"session_id":  sessionID,
			"product_id":  productIDStr,
		}).Error("Invalid product ID")
		
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	var req models.CartUpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		span.RecordError(err)
		h.logger.WithError(err).WithFields(logrus.Fields{
			"component":  "cart_handler",
			"action":     "update_item_quantity",
			"session_id": sessionID,
			"product_id": productID,
		}).Error("Invalid request body")
		
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	span.SetAttributes(
		attribute.Int("item.product_id", productID),
		attribute.Int("item.quantity", req.Quantity),
	)

	h.logger.WithFields(logrus.Fields{
		"component":  "cart_handler",
		"action":     "update_item_quantity",
		"session_id": sessionID,
		"product_id": productID,
		"quantity":   req.Quantity,
	}).Info("Updating item quantity")

	cart, err := h.cartService.UpdateItemQuantity(ctx, sessionID, productID, req.Quantity)
	if err != nil {
		span.RecordError(err)
		h.logger.WithError(err).WithFields(logrus.Fields{
			"component":  "cart_handler",
			"action":     "update_item_quantity",
			"session_id": sessionID,
			"product_id": productID,
		}).Error("Failed to update item quantity")
		
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update item quantity",
		})
		return
	}

	span.SetAttributes(
		attribute.Int("cart.item_count", len(cart.Items)),
		attribute.Float64("cart.total", cart.Total),
	)

	h.logger.WithFields(logrus.Fields{
		"component":  "cart_handler",
		"action":     "update_item_quantity",
		"session_id": sessionID,
		"product_id": productID,
		"item_count": len(cart.Items),
		"total":      cart.Total,
	}).Info("Item quantity updated successfully")

	c.JSON(http.StatusOK, gin.H{
		"message": "Item quantity updated",
		"cart":    cart.ToResponse(),
	})
}

// RemoveItem removes an item from the cart
// DELETE /api/v1/cart/items/:productId
func (h *CartHandler) RemoveItem(c *gin.Context) {
	tracer := otel.Tracer("cart-service")
	ctx, span := tracer.Start(c.Request.Context(), "handler.remove_item")
	defer span.End()

	sessionID := h.getSessionID(c)
	span.SetAttributes(attribute.String("cart.session_id", sessionID))

	// Get product ID from URL parameter
	productIDStr := c.Param("productId")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		span.RecordError(err)
		h.logger.WithError(err).WithFields(logrus.Fields{
			"component":  "cart_handler",
			"action":     "remove_item",
			"session_id": sessionID,
			"product_id": productIDStr,
		}).Error("Invalid product ID")
		
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	span.SetAttributes(attribute.Int("item.product_id", productID))

	h.logger.WithFields(logrus.Fields{
		"component":  "cart_handler",
		"action":     "remove_item",
		"session_id": sessionID,
		"product_id": productID,
	}).Info("Removing item from cart")

	cart, err := h.cartService.RemoveItem(ctx, sessionID, productID)
	if err != nil {
		span.RecordError(err)
		h.logger.WithError(err).WithFields(logrus.Fields{
			"component":  "cart_handler",
			"action":     "remove_item",
			"session_id": sessionID,
			"product_id": productID,
		}).Error("Failed to remove item from cart")
		
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to remove item from cart",
		})
		return
	}

	span.SetAttributes(
		attribute.Int("cart.item_count", len(cart.Items)),
		attribute.Float64("cart.total", cart.Total),
	)

	h.logger.WithFields(logrus.Fields{
		"component":  "cart_handler",
		"action":     "remove_item",
		"session_id": sessionID,
		"product_id": productID,
		"item_count": len(cart.Items),
		"total":      cart.Total,
	}).Info("Item removed from cart successfully")

	c.JSON(http.StatusOK, gin.H{
		"message": "Item removed from cart",
		"cart":    cart.ToResponse(),
	})
}

// ClearCart clears all items from the cart
// DELETE /api/v1/cart
func (h *CartHandler) ClearCart(c *gin.Context) {
	tracer := otel.Tracer("cart-service")
	ctx, span := tracer.Start(c.Request.Context(), "handler.clear_cart")
	defer span.End()

	sessionID := h.getSessionID(c)
	span.SetAttributes(attribute.String("cart.session_id", sessionID))

	h.logger.WithFields(logrus.Fields{
		"component":  "cart_handler",
		"action":     "clear_cart",
		"session_id": sessionID,
	}).Info("Clearing cart")

	err := h.cartService.DeleteCart(ctx, sessionID)
	if err != nil {
		span.RecordError(err)
		h.logger.WithError(err).WithFields(logrus.Fields{
			"component":  "cart_handler",
			"action":     "clear_cart",
			"session_id": sessionID,
		}).Error("Failed to clear cart")
		
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to clear cart",
		})
		return
	}

	span.SetAttributes(attribute.String("cart.status", "cleared"))

	h.logger.WithFields(logrus.Fields{
		"component":  "cart_handler",
		"action":     "clear_cart",
		"session_id": sessionID,
	}).Info("Cart cleared successfully")

	c.JSON(http.StatusOK, gin.H{
		"message": "Cart cleared successfully",
	})
}
package handlers

import (
	"net/http"
	"strconv"

	"catalog-service/internal/logger"
	"catalog-service/internal/models"
	"catalog-service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ProductHandler handles product-related HTTP requests
type ProductHandler struct {
	productService  *models.ProductService
	analysisService *services.AnalysisService
}

// NewProductHandler creates a new product handler
func NewProductHandler(productService *models.ProductService, analysisService *services.AnalysisService) *ProductHandler {
	return &ProductHandler{
		productService:  productService,
		analysisService: analysisService,
	}
}

// GetProducts handles GET /api/v1/products
func (h *ProductHandler) GetProducts(c *gin.Context) {
	// Parse query parameters for pagination
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 50
	}

	offset := (page - 1) * limit

	// Get products from database
	products, err := h.productService.GetAllProducts(c, offset, limit)
	if err != nil {
		logger.WithError(err).WithFields(logrus.Fields{
			"component": "handler",
			"action":    "get_products",
			"page":      page,
			"limit":     limit,
		}).Error("Failed to retrieve products")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve products",
		})
		return
	}

	// Convert to response format
	var responses []models.ProductResponse
	for _, product := range products {
		responses = append(responses, product.ToResponse())
	}

	logger.WithFields(logrus.Fields{
		"component": "handler",
		"action":    "get_products",
		"page":      page,
		"limit":     limit,
		"count":     len(responses),
	}).Info("Retrieved products")

	c.JSON(http.StatusOK, gin.H{
		"data":  responses,
		"page":  page,
		"limit": limit,
		"count": len(responses),
	})
}

// GetProduct handles GET /api/v1/products/:id
func (h *ProductHandler) GetProduct(c *gin.Context) {
	// Parse product ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.WithError(err).WithFields(logrus.Fields{
			"component": "handler",
			"action":    "get_product",
			"id_param":  idStr,
		}).Error("Invalid product ID")

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	// Get product from database
	product, err := h.productService.GetProduct(c, id)
	if err != nil {
		if err.Error() == "product not found" {
			logger.WithFields(logrus.Fields{
				"component":  "handler",
				"action":     "get_product",
				"product_id": id,
			}).Warn("Product not found")

			c.JSON(http.StatusNotFound, gin.H{
				"error": "Product not found",
			})
			return
		}

		logger.WithError(err).WithFields(logrus.Fields{
			"component":  "handler",
			"action":     "get_product",
			"product_id": id,
		}).Error("Failed to retrieve product")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve product",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": product.ToResponse(),
	})
}

// CreateProduct handles POST /api/v1/products
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req models.ProductCreateRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithError(err).WithFields(logrus.Fields{
			"component": "handler",
			"action":    "create_product",
		}).Error("Invalid request data")

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Create product in database
	product, err := h.productService.CreateProduct(c, req)
	if err != nil {
		logger.WithError(err).WithFields(logrus.Fields{
			"component": "handler",
			"action":    "create_product",
			"name":      req.Name,
			"price":     req.Price,
		}).Error("Failed to create product")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create product",
		})
		return
	}

	logger.WithFields(logrus.Fields{
		"component":  "handler",
		"action":     "create_product",
		"product_id": product.ID,
		"name":       product.Name,
		"price":      product.Price,
	}).Info("Product created successfully")

	c.JSON(http.StatusCreated, gin.H{
		"data": product.ToResponse(),
	})
}

// UpdateProduct handles PUT /api/v1/products/:id
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	// Parse product ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.WithError(err).WithFields(logrus.Fields{
			"component": "handler",
			"action":    "update_product",
			"id_param":  idStr,
		}).Error("Invalid product ID")

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	var req models.ProductUpdateRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithError(err).WithFields(logrus.Fields{
			"component":  "handler",
			"action":     "update_product",
			"product_id": id,
		}).Error("Invalid request data")

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Update product in database
	product, err := h.productService.UpdateProduct(c, id, req)
	if err != nil {
		if err.Error() == "product not found" {
			logger.WithFields(logrus.Fields{
				"component":  "handler",
				"action":     "update_product",
				"product_id": id,
			}).Warn("Product not found")

			c.JSON(http.StatusNotFound, gin.H{
				"error": "Product not found",
			})
			return
		}

		logger.WithError(err).WithFields(logrus.Fields{
			"component":  "handler",
			"action":     "update_product",
			"product_id": id,
		}).Error("Failed to update product")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update product",
		})
		return
	}

	logger.WithFields(logrus.Fields{
		"component":  "handler",
		"action":     "update_product",
		"product_id": product.ID,
		"name":       product.Name,
		"price":      product.Price,
	}).Info("Product updated successfully")

	c.JSON(http.StatusOK, product.ToResponse())
}

// DeleteProduct handles DELETE /api/v1/products/:id
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	// Parse product ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.WithError(err).WithFields(logrus.Fields{
			"component": "handler",
			"action":    "delete_product",
			"id_param":  idStr,
		}).Error("Invalid product ID")

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	// Delete product from database
	err = h.productService.DeleteProduct(c, id)
	if err != nil {
		if err.Error() == "product not found" {
			logger.WithFields(logrus.Fields{
				"component":  "handler",
				"action":     "delete_product",
				"product_id": id,
			}).Warn("Product not found")

			c.JSON(http.StatusNotFound, gin.H{
				"error": "Product not found",
			})
			return
		}

		logger.WithError(err).WithFields(logrus.Fields{
			"component":  "handler",
			"action":     "delete_product",
			"product_id": id,
		}).Error("Failed to delete product")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete product",
		})
		return
	}

	logger.WithFields(logrus.Fields{
		"component":  "handler",
		"action":     "delete_product",
		"product_id": id,
	}).Info("Product deleted successfully")

	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
	})
}

// AnalyzeProduct handles GET /api/v1/products/analyze?id=123
// This endpoint demonstrates complex distributed tracing with multiple spans
func (h *ProductHandler) AnalyzeProduct(c *gin.Context) {
	// Parse optional product ID
	var productID *int
	if idStr := c.Query("id"); idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			logger.WithError(err).WithFields(logrus.Fields{
				"component": "handler",
				"action":    "analyze_product",
				"id_param":  idStr,
			}).Error("Invalid product ID")

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid product ID",
			})
			return
		}
		productID = &id
	}

	// Perform complex analysis that creates multiple spans
	result, err := h.analysisService.AnalyzeProduct(c, productID)
	if err != nil {
		logger.WithError(err).WithFields(logrus.Fields{
			"component":  "handler",
			"action":     "analyze_product",
			"product_id": productID,
		}).Error("Failed to analyze product")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to analyze product",
		})
		return
	}

	logger.WithFields(logrus.Fields{
		"component":  "handler",
		"action":     "analyze_product",
		"product_id": productID,
		"duration":   result.TotalDurationMs,
	}).Info("Product analysis completed")

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}

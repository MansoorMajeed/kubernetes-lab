package grpc

import (
	"context"
	"database/sql"
	"fmt"

	"catalog-service/internal/logger"
	"catalog-service/internal/models"

	catalogpb "github.com/mansoormajeed/kubernetes-lab/proto/catalog"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CatalogGRPCServer struct {
	catalogpb.UnimplementedCatalogServiceServer
	db *sql.DB
}

func NewCatalogGRPCServer(db *sql.DB) *CatalogGRPCServer {
	return &CatalogGRPCServer{
		db: db,
	}
}

func (s *CatalogGRPCServer) ValidateProduct(ctx context.Context, req *catalogpb.ProductValidationRequest) (*catalogpb.ProductValidationResponse, error) {
	tracer := otel.Tracer("catalog-grpc")
	ctx, span := tracer.Start(ctx, "ValidateProduct")
	defer span.End()

	span.SetAttributes(
		attribute.String("product.id", req.ProductId),
		attribute.Int("product.requested_quantity", int(req.Quantity)),
	)

	logger.WithFields(logrus.Fields{
		"component":  "grpc",
		"action":     "validate_product",
		"product_id": req.ProductId,
		"quantity":   req.Quantity,
	}).Info("Product validation request received")

	if req.ProductId == "" {
		span.RecordError(fmt.Errorf("product ID is required"))
		span.SetStatus(codes.Error, "product ID is required")
		return nil, status.Error(codes.InvalidArgument, "product ID is required")
	}

	if req.Quantity <= 0 {
		span.RecordError(fmt.Errorf("quantity must be positive"))
		span.SetStatus(codes.Error, "quantity must be positive")
		return nil, status.Error(codes.InvalidArgument, "quantity must be positive")
	}

	// Query product from database
	var product models.Product
	query := `SELECT id, name, description, price, stock_quantity, category_id, image_url, created_at, updated_at 
			  FROM products WHERE id = $1`
	
	err := s.db.QueryRowContext(ctx, query, req.ProductId).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.StockQuantity,
		&product.CategoryID,
		&product.ImageURL,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			span.SetAttributes(attribute.Bool("product.found", false))
			logger.WithFields(logrus.Fields{
				"component":  "grpc",
				"action":     "validate_product",
				"product_id": req.ProductId,
				"error":      "product not found",
			}).Warn("Product not found")

			return &catalogpb.ProductValidationResponse{
				Valid:        false,
				InStock:      false,
				ProductName:  "",
				Price:        0,
				ErrorMessage: "Product not found",
			}, nil
		}

		span.RecordError(err)
		span.SetStatus(codes.Error, "database query failed")
		logger.WithError(err).WithFields(logrus.Fields{
			"component":  "grpc",
			"action":     "validate_product",
			"product_id": req.ProductId,
		}).Error("Failed to query product")

		return nil, status.Error(codes.Internal, "Failed to retrieve product")
	}

	// Check stock availability
	inStock := product.StockQuantity >= int(req.Quantity)
	
	span.SetAttributes(
		attribute.Bool("product.found", true),
		attribute.String("product.name", product.Name),
		attribute.Float64("product.price", product.Price),
		attribute.Int("product.stock_quantity", product.StockQuantity),
		attribute.Bool("product.in_stock", inStock),
	)

	response := &catalogpb.ProductValidationResponse{
		Valid:             true,
		InStock:           inStock,
		AvailableQuantity: int32(product.StockQuantity),
		ProductName:       product.Name,
		Price:             product.Price,
	}

	if !inStock {
		response.ErrorMessage = fmt.Sprintf("Insufficient stock. Available: %d, Requested: %d", product.StockQuantity, req.Quantity)
	}

	logger.WithFields(logrus.Fields{
		"component":          "grpc",
		"action":             "validate_product",
		"product_id":         req.ProductId,
		"product_name":       product.Name,
		"in_stock":           inStock,
		"available_quantity": product.StockQuantity,
		"requested_quantity": req.Quantity,
	}).Info("Product validation completed")

	return response, nil
}

func (s *CatalogGRPCServer) GetProductPrice(ctx context.Context, req *catalogpb.ProductPriceRequest) (*catalogpb.ProductPriceResponse, error) {
	tracer := otel.Tracer("catalog-grpc")
	ctx, span := tracer.Start(ctx, "GetProductPrice")
	defer span.End()

	span.SetAttributes(attribute.String("product.id", req.ProductId))

	if req.ProductId == "" {
		span.RecordError(fmt.Errorf("product ID is required"))
		span.SetStatus(codes.Error, "product ID is required")
		return nil, status.Error(codes.InvalidArgument, "product ID is required")
	}

	var price float64
	query := `SELECT price FROM products WHERE id = $1`
	
	err := s.db.QueryRowContext(ctx, query, req.ProductId).Scan(&price)
	if err != nil {
		if err == sql.ErrNoRows {
			span.SetAttributes(attribute.Bool("product.found", false))
			return &catalogpb.ProductPriceResponse{
				Found:        false,
				ErrorMessage: "Product not found",
			}, nil
		}

		span.RecordError(err)
		span.SetStatus(codes.Error, "database query failed")
		return nil, status.Error(codes.Internal, "Failed to retrieve product price")
	}

	span.SetAttributes(
		attribute.Bool("product.found", true),
		attribute.Float64("product.price", price),
	)

	return &catalogpb.ProductPriceResponse{
		Found:    true,
		Price:    price,
		Currency: "USD",
	}, nil
}

func (s *CatalogGRPCServer) ValidateCartItems(ctx context.Context, req *catalogpb.CartValidationRequest) (*catalogpb.CartValidationResponse, error) {
	tracer := otel.Tracer("catalog-grpc")
	ctx, span := tracer.Start(ctx, "ValidateCartItems")
	defer span.End()

	span.SetAttributes(attribute.Int("cart.items_count", len(req.Items)))

	logger.WithFields(logrus.Fields{
		"component":   "grpc",
		"action":      "validate_cart_items",
		"items_count": len(req.Items),
	}).Info("Cart validation request received")

	var results []*catalogpb.ProductValidationResponse
	var totalPrice float64
	allValid := true

	for i, item := range req.Items {
		span.AddEvent(fmt.Sprintf("validating_item_%d", i), 
			attribute.String("product_id", item.ProductId),
			attribute.Int("quantity", int(item.Quantity)))

		// Reuse the ValidateProduct method for each item
		validationReq := &catalogpb.ProductValidationRequest{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		}

		result, err := s.ValidateProduct(ctx, validationReq)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed to validate cart item")
			return nil, err
		}

		results = append(results, result)

		if !result.Valid || !result.InStock {
			allValid = false
		} else {
			totalPrice += result.Price * float64(item.Quantity)
		}
	}

	span.SetAttributes(
		attribute.Bool("cart.all_valid", allValid),
		attribute.Float64("cart.total_price", totalPrice),
	)

	logger.WithFields(logrus.Fields{
		"component":   "grpc",
		"action":      "validate_cart_items",
		"items_count": len(req.Items),
		"all_valid":   allValid,
		"total_price": totalPrice,
	}).Info("Cart validation completed")

	return &catalogpb.CartValidationResponse{
		Results:    results,
		AllValid:   allValid,
		TotalPrice: totalPrice,
		Currency:   "USD",
	}, nil
}
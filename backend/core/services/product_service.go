package services

import (
	"context"

	"github.com/google/uuid"
	"oms/backend/core/model"
	"oms/backend/core/types"
)

// ProductService defines the interface for product business logic
type ProductService interface {
	GetAll(ctx context.Context) ([]*model.Product, error)
	GetByID(ctx context.Context, productID uuid.UUID) (*model.Product, error)
}

// productService implements ProductService
type productService struct {
	productStore types.ProductStore
}

// NewProductService creates a new ProductService
func NewProductService(productStore types.ProductStore) ProductService {
	return &productService{
		productStore: productStore,
	}
}

// GetAll retrieves all products
func (s *productService) GetAll(ctx context.Context) ([]*model.Product, error) {
	// TODO: Call productStore.GetAll(ctx)
	return nil, nil // Placeholder
}

// GetByID retrieves a product by ID
func (s *productService) GetByID(ctx context.Context, productID uuid.UUID) (*model.Product, error) {
	// TODO: Call productStore.GetByID(ctx, productID)
	return nil, nil // Placeholder
}


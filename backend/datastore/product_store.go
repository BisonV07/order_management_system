package datastore

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"oms/backend/core/model"
	"oms/backend/core/types"
	"gorm.io/gorm"
)

// productStore implements types.ProductStore
type productStore struct {
	db *gorm.DB
}

// NewProductStore creates a new ProductStore
func NewProductStore(db *gorm.DB) types.ProductStore {
	return &productStore{db: db}
}

// GetByID retrieves a product by ID
func (s *productStore) GetByID(ctx context.Context, productID uuid.UUID) (*model.Product, error) {
	var product model.Product
	err := s.db.WithContext(ctx).Where("id = ?", productID).First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &product, nil
}

// GetAll retrieves all products (for dashboard)
// Note: GORM automatically excludes soft-deleted products (where deleted_at IS NOT NULL)
func (s *productStore) GetAll(ctx context.Context) ([]*model.Product, error) {
	var products []*model.Product
	err := s.db.WithContext(ctx).Order("created_at DESC").Find(&products).Error
	if err != nil {
		return nil, err
	}
	// Return empty slice instead of nil if no products found
	if products == nil {
		products = []*model.Product{}
	}
	return products, nil
}

// Create creates a new product
func (s *productStore) Create(ctx context.Context, product *model.Product) error {
	if product.ID == uuid.Nil {
		product.ID = uuid.New()
	}
	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now
	return s.db.WithContext(ctx).Create(product).Error
}

// Update updates an existing product
func (s *productStore) Update(ctx context.Context, productID uuid.UUID, product *model.Product) error {
	product.UpdatedAt = time.Now()
	result := s.db.WithContext(ctx).
		Model(&model.Product{}).
		Where("id = ?", productID).
		Updates(map[string]interface{}{
			"sku":       product.SKU,
			"name":      product.Name,
			"price":     product.Price,
			"metadata":  product.Metadata,
			"updated_at": product.UpdatedAt,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("product not found")
	}
	return nil
}

// Delete deletes a product (soft delete using GORM's DeletedAt)
func (s *productStore) Delete(ctx context.Context, productID uuid.UUID) error {
	result := s.db.WithContext(ctx).Delete(&model.Product{}, "id = ?", productID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("product not found")
	}
	return nil
}


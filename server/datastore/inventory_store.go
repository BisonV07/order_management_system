package datastore

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"oms/server/core/model"
	"oms/server/core/types"
	"gorm.io/gorm"
)

// inventoryStore implements types.InventoryStore
type inventoryStore struct {
	db *gorm.DB
}

// NewInventoryStore creates a new InventoryStore
func NewInventoryStore(db *gorm.DB) types.InventoryStore {
	return &inventoryStore{db: db}
}

// GetByProductID retrieves inventory for a product
func (s *inventoryStore) GetByProductID(ctx context.Context, productID uuid.UUID) (*model.Inventory, error) {
	var inventory model.Inventory
	err := s.db.WithContext(ctx).Where("product_id = ?", productID).First(&inventory).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return zero inventory if not found
			return &model.Inventory{ProductID: productID, Quantity: 0}, nil
		}
		return nil, err
	}
	return &inventory, nil
}

// LockForUpdate locks the inventory row for update (SELECT FOR UPDATE)
// This implements pessimistic locking to prevent overselling
// Note: This method locks the row but doesn't keep the transaction open.
// The caller (order service) should use a transaction for LockForUpdate + DecrementQuantity.
func (s *inventoryStore) LockForUpdate(ctx context.Context, productID uuid.UUID) (*model.Inventory, error) {
	var inventory model.Inventory
	
	// Use SELECT FOR UPDATE to lock the row
	err := s.db.WithContext(ctx).
		Set("gorm:query_option", "FOR UPDATE").
		Where("product_id = ?", productID).
		First(&inventory).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return zero inventory if not found (caller will handle creation)
			return &model.Inventory{ProductID: productID, Quantity: 0}, nil
		}
		return nil, err
	}
	
	return &inventory, nil
}

// DecrementQuantity atomically decrements inventory quantity
// Uses WHERE clause to ensure quantity >= requested quantity
func (s *inventoryStore) DecrementQuantity(ctx context.Context, productID uuid.UUID, quantity int) error {
	result := s.db.WithContext(ctx).
		Model(&model.Inventory{}).
		Where("product_id = ? AND quantity >= ?", productID, quantity).
		Update("quantity", gorm.Expr("quantity - ?", quantity))
	
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("insufficient inventory")
	}
	return nil
}

// IncrementQuantity atomically increments inventory quantity
func (s *inventoryStore) IncrementQuantity(ctx context.Context, productID uuid.UUID, quantity int) error {
	// Use INSERT ... ON CONFLICT (PostgreSQL) or upsert pattern
	result := s.db.WithContext(ctx).
		Model(&model.Inventory{}).
		Where("product_id = ?", productID).
		Update("quantity", gorm.Expr("quantity + ?", quantity))
	
	if result.Error != nil {
		return result.Error
	}
	
	// If no rows were updated, create a new inventory entry
	if result.RowsAffected == 0 {
		inventory := &model.Inventory{
			ProductID: productID,
			Quantity:   quantity,
		}
		return s.db.WithContext(ctx).Create(inventory).Error
	}
	
	return nil
}

// UpdateQuantity sets the inventory quantity for a product (admin only)
func (s *inventoryStore) UpdateQuantity(ctx context.Context, productID uuid.UUID, quantity int) error {
	if quantity < 0 {
		return errors.New("quantity cannot be negative")
	}
	
	// Use upsert pattern: update if exists, create if not
	inventory := &model.Inventory{
		ProductID: productID,
		Quantity:  quantity,
	}
	
	// Try to update first
	result := s.db.WithContext(ctx).
		Model(&model.Inventory{}).
		Where("product_id = ?", productID).
		Update("quantity", quantity)
	
	if result.Error != nil {
		return result.Error
	}
	
	// If no rows were updated, create a new inventory entry
	if result.RowsAffected == 0 {
		return s.db.WithContext(ctx).Create(inventory).Error
	}
	
	return nil
}


package services

import (
	"context"

	"github.com/google/uuid"
	"oms/server/core/model"
	"oms/server/core/types"
)

// InventoryService defines the interface for inventory business logic
type InventoryService interface {
	GetInventory(ctx context.Context, productID uuid.UUID) (*model.Inventory, error)
	LockInventoryForUpdate(ctx context.Context, productID uuid.UUID) (*model.Inventory, error)
	DecrementQuantity(ctx context.Context, productID uuid.UUID, quantity int) error
	IncrementQuantity(ctx context.Context, productID uuid.UUID, quantity int) error
}

// inventoryService implements InventoryService
type inventoryService struct {
	inventoryStore types.InventoryStore
}

// NewInventoryService creates a new InventoryService
func NewInventoryService(inventoryStore types.InventoryStore) InventoryService {
	return &inventoryService{
		inventoryStore: inventoryStore,
	}
}

// GetInventory retrieves inventory for a product
func (s *inventoryService) GetInventory(ctx context.Context, productID uuid.UUID) (*model.Inventory, error) {
	// TODO: Call inventoryStore.GetByProductID(ctx, productID)
	return nil, nil // Placeholder
}

// LockInventoryForUpdate locks the inventory row for update (SELECT FOR UPDATE)
func (s *inventoryService) LockInventoryForUpdate(ctx context.Context, productID uuid.UUID) (*model.Inventory, error) {
	// TODO: Call inventoryStore.LockForUpdate(ctx, productID)
	return nil, nil // Placeholder
}

// DecrementQuantity atomically decrements inventory quantity
func (s *inventoryService) DecrementQuantity(ctx context.Context, productID uuid.UUID, quantity int) error {
	// TODO: Call inventoryStore.DecrementQuantity(ctx, productID, quantity)
	// TODO: Ensure quantity >= 0 (atomic update with WHERE clause)
	return nil // Placeholder
}

// IncrementQuantity atomically increments inventory quantity
func (s *inventoryService) IncrementQuantity(ctx context.Context, productID uuid.UUID, quantity int) error {
	// TODO: Call inventoryStore.IncrementQuantity(ctx, productID, quantity)
	return nil // Placeholder
}


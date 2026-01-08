package fake

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"oms/backend/core/model"
	"oms/backend/core/types"
)

// InventoryStoreFake is a fake implementation of InventoryStore for testing
type InventoryStoreFake struct {
	GetByProductIDFunc     func(ctx context.Context, productID uuid.UUID) (*model.Inventory, error)
	LockForUpdateFunc      func(ctx context.Context, productID uuid.UUID) (*model.Inventory, error)
	DecrementQuantityFunc  func(ctx context.Context, productID uuid.UUID, quantity int) error
	IncrementQuantityFunc  func(ctx context.Context, productID uuid.UUID, quantity int) error
}

// inventoryMap maintains inventory state for fake store with mutex protection for race conditions
var inventoryMap = struct {
	sync.RWMutex
	m map[uuid.UUID]*model.Inventory
}{m: make(map[uuid.UUID]*model.Inventory)}

// productLocks provides per-product locks for pessimistic locking simulation
var productLocks = struct {
	sync.Mutex
	locks map[uuid.UUID]*sync.Mutex
}{locks: make(map[uuid.UUID]*sync.Mutex)}

// getProductLock returns a mutex for a specific product (for pessimistic locking)
func getProductLock(productID uuid.UUID) *sync.Mutex {
	productLocks.Lock()
	defer productLocks.Unlock()
	if productLocks.locks[productID] == nil {
		productLocks.locks[productID] = &sync.Mutex{}
	}
	return productLocks.locks[productID]
}

// GetByProductID implements types.InventoryStore
func (f *InventoryStoreFake) GetByProductID(ctx context.Context, productID uuid.UUID) (*model.Inventory, error) {
	if f.GetByProductIDFunc != nil {
		return f.GetByProductIDFunc(ctx, productID)
	}
	// Thread-safe read
	inventoryMap.RLock()
	defer inventoryMap.RUnlock()
	
	inv, exists := inventoryMap.m[productID]
	if !exists {
		// Need to upgrade to write lock to create
		inventoryMap.RUnlock()
		inventoryMap.Lock()
		// Double-check after acquiring write lock
		inv, exists = inventoryMap.m[productID]
		if !exists {
			defaultQty := getDefaultQuantity(productID)
			inv = &model.Inventory{ProductID: productID, Quantity: defaultQty}
			inventoryMap.m[productID] = inv
		}
		inventoryMap.Unlock()
		inventoryMap.RLock()
	}
	// Return a copy to prevent external modification
	copiedInv := *inv
	return &copiedInv, nil
}

// LockForUpdate implements types.InventoryStore
// Uses pessimistic locking - acquires product-specific lock to prevent race conditions
func (f *InventoryStoreFake) LockForUpdate(ctx context.Context, productID uuid.UUID) (*model.Inventory, error) {
	if f.LockForUpdateFunc != nil {
		return f.LockForUpdateFunc(ctx, productID)
	}
	
	// Acquire product-specific lock (pessimistic locking)
	productLock := getProductLock(productID)
	productLock.Lock()
	// Note: Lock is NOT released here - caller must ensure DecrementQuantity/IncrementQuantity
	// releases it, or use a transaction-like pattern
	
	// Get or create inventory from map (with write lock)
	inventoryMap.Lock()
	inv, exists := inventoryMap.m[productID]
	if !exists {
		// Initialize with default quantity based on product ID
		defaultQty := getDefaultQuantity(productID)
		inv = &model.Inventory{ProductID: productID, Quantity: defaultQty}
		inventoryMap.m[productID] = inv
	}
	// Return a copy to prevent external modification
	copiedInv := *inv
	inventoryMap.Unlock()
	
	// Note: productLock is held and will be released in DecrementQuantity
	// This ensures atomicity: Lock -> Check -> Decrement all happen atomically
	return &copiedInv, nil
}

// DecrementQuantity implements types.InventoryStore
// Must be called after LockForUpdate to ensure thread safety
// Note: LockForUpdate already holds the product lock, so we use a flag to track it
func (f *InventoryStoreFake) DecrementQuantity(ctx context.Context, productID uuid.UUID, quantity int) error {
	if f.DecrementQuantityFunc != nil {
		return f.DecrementQuantityFunc(ctx, productID, quantity)
	}
	
	// Get product lock reference (lock is already held by LockForUpdate)
	productLock := getProductLock(productID)
	
	// Get inventory with write lock
	inventoryMap.Lock()
	inv, exists := inventoryMap.m[productID]
	if !exists {
		// This shouldn't happen if LockForUpdate was called first, but handle it
		defaultQty := getDefaultQuantity(productID)
		inv = &model.Inventory{ProductID: productID, Quantity: defaultQty}
		inventoryMap.m[productID] = inv
	}
	
	// Check availability
	if inv.Quantity < quantity {
		inventoryMap.Unlock()
		productLock.Unlock() // Release lock before returning error
		return fmt.Errorf("insufficient inventory: requested %d, available %d", quantity, inv.Quantity)
	}
	
	// Decrement atomically
	inv.Quantity -= quantity
	inventoryMap.Unlock()
	
	// Release the product lock that was acquired in LockForUpdate
	productLock.Unlock()
	
	return nil
}

// IncrementQuantity implements types.InventoryStore
// Thread-safe increment operation
func (f *InventoryStoreFake) IncrementQuantity(ctx context.Context, productID uuid.UUID, quantity int) error {
	if f.IncrementQuantityFunc != nil {
		return f.IncrementQuantityFunc(ctx, productID, quantity)
	}
	
	// Acquire product-specific lock
	productLock := getProductLock(productID)
	productLock.Lock()
	defer productLock.Unlock()
	
	// Get inventory with write lock
	inventoryMap.Lock()
	defer inventoryMap.Unlock()
	
	inv, exists := inventoryMap.m[productID]
	if !exists {
		inv = &model.Inventory{ProductID: productID, Quantity: 0}
		inventoryMap.m[productID] = inv
	}
	
	// Increment atomically
	inv.Quantity += quantity
	return nil
}

// getDefaultQuantity returns default inventory quantity for a product
// This matches the initial values shown in the products endpoint
func getDefaultQuantity(productID uuid.UUID) int {
	// Map product IDs to their default inventory quantities
	productDefaults := map[string]int{
		"550e8400-e29b-41d4-a716-446655440000": 100, // Laptop Computer
		"550e8400-e29b-41d4-a716-446655440001": 150, // Wireless Mouse
		"550e8400-e29b-41d4-a716-446655440002": 75,  // Mechanical Keyboard
		"550e8400-e29b-41d4-a716-446655440003": 200, // USB-C Hub
		"550e8400-e29b-41d4-a716-446655440004": 50,  // Monitor Stand
	}
	
	if qty, ok := productDefaults[productID.String()]; ok {
		return qty
	}
	return 100 // Default fallback
}

// Ensure InventoryStoreFake implements types.InventoryStore
var _ types.InventoryStore = (*InventoryStoreFake)(nil)


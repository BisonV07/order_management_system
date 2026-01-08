package fake

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"oms/server/core/model"
	"oms/server/core/types"
)

// OrderStoreFake is a fake implementation of OrderStore for testing
type OrderStoreFake struct {
	CreateFunc      func(ctx context.Context, order *model.Order) error
	GetByIDFunc     func(ctx context.Context, orderID uuid.UUID) (*model.Order, error)
	GetByUserIDFunc func(ctx context.Context, userID int) ([]*model.Order, error)
	UpdateStatusFunc func(ctx context.Context, orderID uuid.UUID, status model.OrderStatus) error
}

var orders = struct {
	sync.RWMutex
	m map[uuid.UUID]*model.Order
}{m: make(map[uuid.UUID]*model.Order)}

// Create implements types.OrderStore
func (f *OrderStoreFake) Create(ctx context.Context, order *model.Order) error {
	if f.CreateFunc != nil {
		return f.CreateFunc(ctx, order)
	}
	orders.Lock()
	defer orders.Unlock()
	if order.ID == uuid.Nil {
		order.ID = uuid.New()
	}
	now := time.Now()
	order.CreatedAt = now
	order.UpdatedAt = now
	orders.m[order.ID] = order
	return nil
}

// GetByID implements types.OrderStore
func (f *OrderStoreFake) GetByID(ctx context.Context, orderID uuid.UUID) (*model.Order, error) {
	if f.GetByIDFunc != nil {
		return f.GetByIDFunc(ctx, orderID)
	}
	orders.RLock()
	defer orders.RUnlock()
	order, exists := orders.m[orderID]
	if !exists {
		return nil, fmt.Errorf("order not found")
	}
	// Return a copy
	copiedOrder := *order
	return &copiedOrder, nil
}

// GetByUserID implements types.OrderStore
func (f *OrderStoreFake) GetByUserID(ctx context.Context, userID int) ([]*model.Order, error) {
	if f.GetByUserIDFunc != nil {
		return f.GetByUserIDFunc(ctx, userID)
	}
	orders.RLock()
	defer orders.RUnlock()
	var userOrders []*model.Order
	for _, order := range orders.m {
		if order.UserID == userID {
			copiedOrder := *order
			userOrders = append(userOrders, &copiedOrder)
		}
	}
	return userOrders, nil
}

// GetAll implements types.OrderStore - returns all orders (for admin)
func (f *OrderStoreFake) GetAll(ctx context.Context) ([]*model.Order, error) {
	orders.RLock()
	defer orders.RUnlock()
	var allOrders []*model.Order
	for _, order := range orders.m {
		copiedOrder := *order
		allOrders = append(allOrders, &copiedOrder)
	}
	return allOrders, nil
}

// UpdateStatus implements types.OrderStore
func (f *OrderStoreFake) UpdateStatus(ctx context.Context, orderID uuid.UUID, status model.OrderStatus) error {
	if f.UpdateStatusFunc != nil {
		return f.UpdateStatusFunc(ctx, orderID, status)
	}
	orders.Lock()
	defer orders.Unlock()
	order, exists := orders.m[orderID]
	if !exists {
		return fmt.Errorf("order not found")
	}
	order.CurrentStatus = status
	order.UpdatedAt = time.Now()
	return nil
}

// Ensure OrderStoreFake implements types.OrderStore
var _ types.OrderStore = (*OrderStoreFake)(nil)


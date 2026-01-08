package fake

import (
	"context"

	"github.com/google/uuid"
	"oms/backend/core/model"
	"oms/backend/core/services"
)

// OrderServiceFake is a fake implementation of OrderService for testing
type OrderServiceFake struct {
	CreateOrderFunc        func(ctx context.Context, userID int, productID uuid.UUID, quantity int, metadata model.JSONB) (*model.Order, error)
	UpdateOrderStatusFunc  func(ctx context.Context, orderID uuid.UUID, newStatus model.OrderStatus, updatedBy int) (*model.Order, error)
	GetOrderByIDFunc       func(ctx context.Context, orderID uuid.UUID) (*model.Order, error)
	GetOrdersByUserIDFunc  func(ctx context.Context, userID int) ([]*model.Order, error)
	GetAllOrdersFunc       func(ctx context.Context) ([]*model.Order, error)
	GetOrderHistoryFunc    func(ctx context.Context, orderID uuid.UUID) ([]*model.OrderStateLog, error)
}

// NewOrderServiceFake creates a new fake OrderService
func NewOrderServiceFake() services.OrderService {
	return &OrderServiceFake{}
}

// CreateOrder implements services.OrderService
func (f *OrderServiceFake) CreateOrder(ctx context.Context, userID int, productID uuid.UUID, quantity int, metadata model.JSONB) (*model.Order, error) {
	if f.CreateOrderFunc != nil {
		return f.CreateOrderFunc(ctx, userID, productID, quantity, metadata)
	}
	return nil, nil
}

// UpdateOrderStatus implements services.OrderService
func (f *OrderServiceFake) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, newStatus model.OrderStatus, updatedBy int) (*model.Order, error) {
	if f.UpdateOrderStatusFunc != nil {
		return f.UpdateOrderStatusFunc(ctx, orderID, newStatus, updatedBy)
	}
	return nil, nil
}

// GetOrderByID implements services.OrderService
func (f *OrderServiceFake) GetOrderByID(ctx context.Context, orderID uuid.UUID) (*model.Order, error) {
	if f.GetOrderByIDFunc != nil {
		return f.GetOrderByIDFunc(ctx, orderID)
	}
	return nil, nil
}

// GetOrdersByUserID implements services.OrderService
func (f *OrderServiceFake) GetOrdersByUserID(ctx context.Context, userID int) ([]*model.Order, error) {
	if f.GetOrdersByUserIDFunc != nil {
		return f.GetOrdersByUserIDFunc(ctx, userID)
	}
	return []*model.Order{}, nil
}

// GetAllOrders implements services.OrderService
func (f *OrderServiceFake) GetAllOrders(ctx context.Context) ([]*model.Order, error) {
	if f.GetAllOrdersFunc != nil {
		return f.GetAllOrdersFunc(ctx)
	}
	return []*model.Order{}, nil
}

// GetOrderHistory implements services.OrderService
func (f *OrderServiceFake) GetOrderHistory(ctx context.Context, orderID uuid.UUID) ([]*model.OrderStateLog, error) {
	if f.GetOrderHistoryFunc != nil {
		return f.GetOrderHistoryFunc(ctx, orderID)
	}
	return []*model.OrderStateLog{}, nil
}

// Ensure OrderServiceFake implements services.OrderService
var _ services.OrderService = (*OrderServiceFake)(nil)


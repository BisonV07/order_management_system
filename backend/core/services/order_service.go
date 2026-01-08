package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"oms/backend/core/model"
	"oms/backend/core/types"
)

// OrderService defines the interface for order business logic
type OrderService interface {
	CreateOrder(ctx context.Context, userID int, productID uuid.UUID, quantity int, metadata model.JSONB) (*model.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, newStatus model.OrderStatus, updatedBy int) (*model.Order, error)
	GetOrderByID(ctx context.Context, orderID uuid.UUID) (*model.Order, error)
	GetOrdersByUserID(ctx context.Context, userID int) ([]*model.Order, error)
	GetAllOrders(ctx context.Context) ([]*model.Order, error)
	GetOrderHistory(ctx context.Context, orderID uuid.UUID) ([]*model.OrderStateLog, error)
}

// orderService implements OrderService
type orderService struct {
	orderStore         types.OrderStore
	inventoryStore     types.InventoryStore
	orderStateLogStore types.OrderStateLogStore
	fsmValidator       types.FSMValidator
}

// NewOrderService creates a new OrderService
func NewOrderService(
	orderStore types.OrderStore,
	inventoryStore types.InventoryStore,
	orderStateLogStore types.OrderStateLogStore,
	fsmValidator types.FSMValidator,
) OrderService {
	return &orderService{
		orderStore:         orderStore,
		inventoryStore:     inventoryStore,
		orderStateLogStore: orderStateLogStore,
		fsmValidator:       fsmValidator,
	}
}

// CreateOrder creates a new order with inventory locking
// Uses pessimistic locking (SELECT FOR UPDATE) to prevent overselling
func (s *orderService) CreateOrder(ctx context.Context, userID int, productID uuid.UUID, quantity int, metadata model.JSONB) (*model.Order, error) {
	// Validate inputs
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", userID)
	}
	if quantity <= 0 {
		return nil, fmt.Errorf("invalid quantity: %d", quantity)
	}

	// Lock inventory row for update (pessimistic locking)
	inventory, err := s.inventoryStore.LockForUpdate(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to lock inventory: %w", err)
	}

	// Check stock availability
	if inventory.Quantity < quantity {
		return nil, fmt.Errorf("insufficient inventory: requested %d, available %d", quantity, inventory.Quantity)
	}

	// Deduct inventory atomically
	err = s.inventoryStore.DecrementQuantity(ctx, productID, quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to decrement inventory: %w", err)
	}

	// Create order with status ORDERED and metadata
	order := &model.Order{
		UserID:        userID,
		ProductID:     productID,
		Quantity:      quantity,
		CurrentStatus: model.OrderStatusOrdered,
		Metadata:      metadata,
	}

	err = s.orderStore.Create(ctx, order)
	if err != nil {
		// If order creation fails, restore inventory
		restoreErr := s.inventoryStore.IncrementQuantity(ctx, productID, quantity)
		if restoreErr != nil {
			// Log restore error but return original error
			return nil, fmt.Errorf("failed to create order: %w (inventory restore also failed: %v)", err, restoreErr)
		}
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}

// UpdateOrderStatus updates the order status with FSM validation
func (s *orderService) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, newStatus model.OrderStatus, updatedBy int) (*model.Order, error) {
	// Fetch current order
	order, err := s.orderStore.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	currentStatus := order.CurrentStatus

	// Validate transition using FSM
	err = s.fsmValidator.ValidateTransition(currentStatus, newStatus)
	if err != nil {
		return nil, fmt.Errorf("invalid transition from %s to %s: %w", currentStatus, newStatus, err)
	}

	// Idempotency: If same status, return current order
	if currentStatus == newStatus {
		return order, nil
	}

	// Update order status
	err = s.orderStore.UpdateStatus(ctx, orderID, newStatus)
	if err != nil {
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	// Create audit log entry
	stateLog := &model.OrderStateLog{
		OrderID:        orderID,
		PreviousStatus: currentStatus,
		NewStatus:      newStatus,
		UpdatedBy:      updatedBy,
		UpdatedAt:      time.Now(),
	}
	err = s.orderStateLogStore.Create(ctx, stateLog)
	if err != nil {
		// Log error but don't fail the update
		// TODO: Add proper logging
	}

	// If status is CANCELLED, restore inventory
	if s.fsmValidator.RequiresInventoryRestore(newStatus) {
		err = s.inventoryStore.IncrementQuantity(ctx, order.ProductID, order.Quantity)
		if err != nil {
			// Log error but don't fail the update
			// TODO: Add proper logging and potentially rollback
		}
	}

	// Fetch updated order
	updatedOrder, err := s.orderStore.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated order: %w", err)
	}

	return updatedOrder, nil
}

// GetOrderByID retrieves an order by ID
func (s *orderService) GetOrderByID(ctx context.Context, orderID uuid.UUID) (*model.Order, error) {
	return s.orderStore.GetByID(ctx, orderID)
}

// GetOrdersByUserID retrieves all orders for a user
func (s *orderService) GetOrdersByUserID(ctx context.Context, userID int) ([]*model.Order, error) {
	return s.orderStore.GetByUserID(ctx, userID)
}

// GetAllOrders retrieves all orders (for admin)
func (s *orderService) GetAllOrders(ctx context.Context) ([]*model.Order, error) {
	return s.orderStore.GetAll(ctx)
}

// GetOrderHistory retrieves the state change history for an order
func (s *orderService) GetOrderHistory(ctx context.Context, orderID uuid.UUID) ([]*model.OrderStateLog, error) {
	return s.orderStateLogStore.GetByOrderID(ctx, orderID)
}


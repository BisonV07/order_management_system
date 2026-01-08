package types

import (
	"context"

	"github.com/google/uuid"
	"oms/server/core/model"
)

// OrderStore defines the interface for order data access
type OrderStore interface {
	Create(ctx context.Context, order *model.Order) error
	GetByID(ctx context.Context, orderID uuid.UUID) (*model.Order, error)
	GetByUserID(ctx context.Context, userID int) ([]*model.Order, error)
	GetAll(ctx context.Context) ([]*model.Order, error)
	UpdateStatus(ctx context.Context, orderID uuid.UUID, status model.OrderStatus) error
}

// InventoryStore defines the interface for inventory data access
type InventoryStore interface {
	GetByProductID(ctx context.Context, productID uuid.UUID) (*model.Inventory, error)
	LockForUpdate(ctx context.Context, productID uuid.UUID) (*model.Inventory, error)
	DecrementQuantity(ctx context.Context, productID uuid.UUID, quantity int) error
	IncrementQuantity(ctx context.Context, productID uuid.UUID, quantity int) error
	UpdateQuantity(ctx context.Context, productID uuid.UUID, quantity int) error // Admin: Set inventory quantity
}

// ProductStore defines the interface for product data access
type ProductStore interface {
	GetByID(ctx context.Context, productID uuid.UUID) (*model.Product, error)
	GetAll(ctx context.Context) ([]*model.Product, error)
	Create(ctx context.Context, product *model.Product) error // Admin: Create new product
	Update(ctx context.Context, productID uuid.UUID, product *model.Product) error // Admin: Update product
	Delete(ctx context.Context, productID uuid.UUID) error // Admin: Delete product (soft delete)
}

// OrderStateLogStore defines the interface for order state log data access
type OrderStateLogStore interface {
	Create(ctx context.Context, log *model.OrderStateLog) error
	GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*model.OrderStateLog, error)
}

// FSMValidator defines the interface for FSM validation
type FSMValidator interface {
	ValidateTransition(currentStatus, newStatus model.OrderStatus) error
	IsValidStatus(status model.OrderStatus) bool
	RequiresInventoryRestore(status model.OrderStatus) bool
}

// UserStore defines the interface for user data access
type UserStore interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, userID int) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
}


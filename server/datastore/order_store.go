package datastore

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"oms/server/core/model"
	"oms/server/core/types"
	"gorm.io/gorm"
)

// orderStore implements types.OrderStore
type orderStore struct {
	db *gorm.DB
}

// NewOrderStore creates a new OrderStore
func NewOrderStore(db *gorm.DB) types.OrderStore {
	return &orderStore{db: db}
}

// Create creates a new order
func (s *orderStore) Create(ctx context.Context, order *model.Order) error {
	if order.ID == uuid.Nil {
		order.ID = uuid.New()
	}
	now := time.Now()
	order.CreatedAt = now
	order.UpdatedAt = now
	return s.db.WithContext(ctx).Create(order).Error
}

// GetByID retrieves an order by ID
func (s *orderStore) GetByID(ctx context.Context, orderID uuid.UUID) (*model.Order, error) {
	var order model.Order
	err := s.db.WithContext(ctx).Where("id = ?", orderID).First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}
	return &order, nil
}

// GetByUserID retrieves all orders for a given user ID
func (s *orderStore) GetByUserID(ctx context.Context, userID int) ([]*model.Order, error) {
	var orders []*model.Order
	err := s.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&orders).Error
	return orders, err
}

// GetAll retrieves all orders (for admin)
func (s *orderStore) GetAll(ctx context.Context) ([]*model.Order, error) {
	var orders []*model.Order
	err := s.db.WithContext(ctx).Order("created_at DESC").Find(&orders).Error
	return orders, err
}

// UpdateStatus updates the order status
func (s *orderStore) UpdateStatus(ctx context.Context, orderID uuid.UUID, status model.OrderStatus) error {
	result := s.db.WithContext(ctx).
		Model(&model.Order{}).
		Where("id = ?", orderID).
		Updates(map[string]interface{}{
			"current_status": status,
			"updated_at":     time.Now(),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("order not found")
	}
	return nil
}


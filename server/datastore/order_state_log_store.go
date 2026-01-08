package datastore

import (
	"context"

	"github.com/google/uuid"
	"oms/server/core/model"
	"oms/server/core/types"
	"gorm.io/gorm"
)

// orderStateLogStore implements types.OrderStateLogStore
type orderStateLogStore struct {
	db *gorm.DB
}

// NewOrderStateLogStore creates a new OrderStateLogStore
func NewOrderStateLogStore(db *gorm.DB) types.OrderStateLogStore {
	return &orderStateLogStore{db: db}
}

// Create creates a new order state log entry
func (s *orderStateLogStore) Create(ctx context.Context, log *model.OrderStateLog) error {
	return s.db.WithContext(ctx).Create(log).Error
}

// GetByOrderID retrieves all state logs for an order
func (s *orderStateLogStore) GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*model.OrderStateLog, error) {
	var logs []*model.OrderStateLog
	err := s.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Order("updated_at ASC").
		Find(&logs).Error
	return logs, err
}


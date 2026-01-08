package fake

import (
	"context"

	"github.com/google/uuid"
	"oms/backend/core/model"
	"oms/backend/core/types"
)

// OrderStateLogStoreFake is a fake implementation of OrderStateLogStore for testing
type OrderStateLogStoreFake struct {
	CreateFunc    func(ctx context.Context, log *model.OrderStateLog) error
	GetByOrderIDFunc func(ctx context.Context, orderID uuid.UUID) ([]*model.OrderStateLog, error)
}

var orderStateLogs = make(map[uuid.UUID][]*model.OrderStateLog)

// Create implements types.OrderStateLogStore
func (f *OrderStateLogStoreFake) Create(ctx context.Context, log *model.OrderStateLog) error {
	if f.CreateFunc != nil {
		return f.CreateFunc(ctx, log)
	}
	orderStateLogs[log.OrderID] = append(orderStateLogs[log.OrderID], log)
	return nil
}

// GetByOrderID implements types.OrderStateLogStore
func (f *OrderStateLogStoreFake) GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*model.OrderStateLog, error) {
	if f.GetByOrderIDFunc != nil {
		return f.GetByOrderIDFunc(ctx, orderID)
	}
	logs, exists := orderStateLogs[orderID]
	if !exists {
		return []*model.OrderStateLog{}, nil
	}
	return logs, nil
}

// Ensure OrderStateLogStoreFake implements types.OrderStateLogStore
var _ types.OrderStateLogStore = (*OrderStateLogStoreFake)(nil)


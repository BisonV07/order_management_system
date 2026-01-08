package fsm

import (
	"fmt"

	"oms/server/core/model"
)

// StateTransition defines allowed transitions for order states
var StateTransition = map[model.OrderStatus][]model.OrderStatus{
	model.OrderStatusOrdered: {
		model.OrderStatusShipped,
		model.OrderStatusCancelled,
	},
	model.OrderStatusShipped: {
		model.OrderStatusDelivered,
	},
	model.OrderStatusDelivered: {
		// No transitions allowed from DELIVERED
	},
	model.OrderStatusCancelled: {
		// No transitions allowed from CANCELLED
	},
}

// ValidateTransition checks if a state transition is allowed
// Returns error if transition is invalid, nil if valid
func ValidateTransition(currentStatus, newStatus model.OrderStatus) error {
	// Same state is allowed (idempotency)
	if currentStatus == newStatus {
		return nil
	}

	allowedTransitions, exists := StateTransition[currentStatus]
	if !exists {
		return fmt.Errorf("invalid current status: %s", currentStatus)
	}

	for _, allowed := range allowedTransitions {
		if allowed == newStatus {
			return nil
		}
	}

	return fmt.Errorf("invalid transition from %s to %s", currentStatus, newStatus)
}

// IsValidStatus checks if a status is a valid order status
func IsValidStatus(status model.OrderStatus) bool {
	_, exists := StateTransition[status]
	return exists
}

// RequiresInventoryRestore checks if transitioning to this status requires inventory to be restored
func RequiresInventoryRestore(status model.OrderStatus) bool {
	return status == model.OrderStatusCancelled
}


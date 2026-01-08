package fsm

import (
	"oms/server/core/model"
	"oms/server/core/types"
)

// validator implements types.FSMValidator using the state machine
type validator struct{}

// NewValidator creates a new FSM validator
func NewValidator() types.FSMValidator {
	return &validator{}
}

// ValidateTransition implements types.FSMValidator
func (v *validator) ValidateTransition(currentStatus, newStatus model.OrderStatus) error {
	return ValidateTransition(currentStatus, newStatus)
}

// IsValidStatus implements types.FSMValidator
func (v *validator) IsValidStatus(status model.OrderStatus) bool {
	return IsValidStatus(status)
}

// RequiresInventoryRestore implements types.FSMValidator
func (v *validator) RequiresInventoryRestore(status model.OrderStatus) bool {
	return RequiresInventoryRestore(status)
}

// Ensure validator implements types.FSMValidator
var _ types.FSMValidator = (*validator)(nil)


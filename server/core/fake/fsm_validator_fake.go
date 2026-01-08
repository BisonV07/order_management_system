package fake

import (
	"oms/server/core/model"
	"oms/server/core/types"
)

// FSMValidatorFake is a fake implementation of FSMValidator for testing
type FSMValidatorFake struct {
	ValidateTransitionFunc        func(currentStatus, newStatus model.OrderStatus) error
	IsValidStatusFunc             func(status model.OrderStatus) bool
	RequiresInventoryRestoreFunc  func(status model.OrderStatus) bool
}

// ValidateTransition implements types.FSMValidator
func (f *FSMValidatorFake) ValidateTransition(currentStatus, newStatus model.OrderStatus) error {
	if f.ValidateTransitionFunc != nil {
		return f.ValidateTransitionFunc(currentStatus, newStatus)
	}
	return nil
}

// IsValidStatus implements types.FSMValidator
func (f *FSMValidatorFake) IsValidStatus(status model.OrderStatus) bool {
	if f.IsValidStatusFunc != nil {
		return f.IsValidStatusFunc(status)
	}
	return true
}

// RequiresInventoryRestore implements types.FSMValidator
func (f *FSMValidatorFake) RequiresInventoryRestore(status model.OrderStatus) bool {
	if f.RequiresInventoryRestoreFunc != nil {
		return f.RequiresInventoryRestoreFunc(status)
	}
	return false
}

// Ensure FSMValidatorFake implements types.FSMValidator
var _ types.FSMValidator = (*FSMValidatorFake)(nil)


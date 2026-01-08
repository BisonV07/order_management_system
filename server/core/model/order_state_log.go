package model

import (
	"time"

	"github.com/google/uuid"
)

// OrderStateLog represents an audit trail of order status changes
type OrderStateLog struct {
	ID            uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrderID       uuid.UUID   `gorm:"type:uuid;not null" json:"order_id"`
	PreviousStatus OrderStatus `gorm:"type:varchar(50)" json:"previous_status"`
	NewStatus     OrderStatus `gorm:"type:varchar(50);not null" json:"new_status"`
	UpdatedBy     int         `gorm:"not null" json:"updated_by"` // User ID who made the change
	UpdatedAt     time.Time   `gorm:"autoCreateTime" json:"updated_at"`
	Order         Order       `gorm:"foreignKey:OrderID" json:"order,omitempty"`
}

// TableName specifies the table name for OrderStateLog
func (OrderStateLog) TableName() string {
	return "order_state_logs"
}


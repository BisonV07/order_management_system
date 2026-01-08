package model

import (
	"time"

	"github.com/google/uuid"
)

// OrderStatus represents the possible states of an order
type OrderStatus string

const (
	OrderStatusOrdered   OrderStatus = "ORDERED"
	OrderStatusShipped   OrderStatus = "SHIPPED"
	OrderStatusDelivered OrderStatus = "DELIVERED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

// Order represents an order in the system
type Order struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID       int        `gorm:"not null" json:"user_id"`
	ProductID    uuid.UUID  `gorm:"type:uuid;not null" json:"product_id"`
	Quantity     int        `gorm:"not null" json:"quantity"`
	CurrentStatus OrderStatus `gorm:"type:varchar(50);not null;default:'ORDERED'" json:"current_status"`
	Metadata     JSONB      `gorm:"type:jsonb" json:"metadata"` // For shipping address and other order details
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	Product      Product    `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

// TableName specifies the table name for Order
func (Order) TableName() string {
	return "orders"
}


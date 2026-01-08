package model

import (
	"github.com/google/uuid"
)

// Inventory represents inventory/stock for a product
type Inventory struct {
	ProductID uuid.UUID `gorm:"type:uuid;primary_key" json:"product_id"`
	Quantity  int       `gorm:"not null" json:"quantity"`
}

// TableName specifies the table name for Inventory
func (Inventory) TableName() string {
	return "inventory"
}


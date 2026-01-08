package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Product represents a product in the system
type Product struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	SKU       string         `gorm:"type:varchar(255);unique;not null" json:"sku"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Price     float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	Metadata  JSONB          `gorm:"type:jsonb" json:"metadata"` // For attributes like Color, Size, Weight
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// JSONB is a type alias for jsonb fields with custom scanning support
type JSONB map[string]interface{}

// Scan implements the sql.Scanner interface for JSONB
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = JSONB{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		*j = JSONB{}
		return nil
	}

	if len(bytes) == 0 {
		*j = JSONB{}
		return nil
	}

	return json.Unmarshal(bytes, j)
}

// Value implements the driver.Valuer interface for JSONB
func (j JSONB) Value() (driver.Value, error) {
	if j == nil || len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
}

// TableName specifies the table name for Product
func (Product) TableName() string {
	return "products"
}


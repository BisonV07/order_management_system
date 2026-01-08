package types

import "time"

// LoginResponse represents the response for login
type LoginResponse struct {
	Token  string `json:"token"`
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
}

// SignupResponse represents the response for signup
type SignupResponse struct {
	Message string `json:"message"`
	UserID  int    `json:"user_id"`
}

// CreateOrderResponse represents the response for order creation
type CreateOrderResponse struct {
	OrderID       string `json:"order_id"`
	CurrentStatus string `json:"current_status"`
	Message       string `json:"message"`
}

// UpdateOrderStatusResponse represents the response for order status update
type UpdateOrderStatusResponse struct {
	OrderID       string    `json:"order_id"`
	PreviousStatus string   `json:"previous_status"`
	CurrentStatus string    `json:"current_status"`
	UpdatedBy     int       `json:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
}

// OrderResponse represents an order in the response
type OrderResponse struct {
	ID            string                 `json:"id"`
	UserID        int                    `json:"user_id"`
	ProductID     string                 `json:"product_id"`
	Quantity      int                    `json:"quantity"`
	CurrentStatus string                 `json:"current_status"`
	Metadata      map[string]interface{} `json:"metadata"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

// OrderHistoryResponse represents an order state change in history
type OrderHistoryResponse struct {
	OrderID        string    `json:"order_id"`
	PreviousStatus string    `json:"previous_status"`
	NewStatus      string    `json:"new_status"`
	UpdatedBy      int       `json:"updated_by"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ProductResponse represents a product in the response
type ProductResponse struct {
	ID        string                 `json:"id"`
	SKU       string                 `json:"sku"`
	Name      string                 `json:"name"`
	Price     float64                `json:"price"`
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// InventoryResponse represents inventory in the response
type InventoryResponse struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

// CreateProductResponse represents the response for product creation
type CreateProductResponse struct {
	ProductID string `json:"product_id"`
	Message   string `json:"message"`
}

// UpdateInventoryResponse represents the response for inventory update
type UpdateInventoryResponse struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Message   string `json:"message"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}


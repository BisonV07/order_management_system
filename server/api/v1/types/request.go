package types

// LoginRequest represents the request body for login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// SignupRequest represents the request body for signup
type SignupRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// CreateOrderRequest represents the request body for creating an order
type CreateOrderRequest struct {
	ProductID       string                 `json:"product_id" binding:"required"` // UUID as string
	Quantity        int                    `json:"quantity" binding:"required,min=1"`
	ShippingAddress map[string]interface{} `json:"shipping_address"` // Shipping address metadata
}

// UpdateOrderStatusRequest represents the request body for updating order status
type UpdateOrderStatusRequest struct {
	CurrentStatus string `json:"current_status" binding:"required"`
}

// CreateProductRequest represents the request body for creating a product (admin only)
type CreateProductRequest struct {
	SKU      string                 `json:"sku" binding:"required"`
	Name     string                 `json:"name" binding:"required"`
	Price    float64                `json:"price" binding:"required,min=0"`
	Metadata map[string]interface{} `json:"metadata"`
}

// UpdateProductRequest represents the request body for updating a product (admin only)
type UpdateProductRequest struct {
	SKU      string                 `json:"sku"`
	Name     string                 `json:"name"`
	Price    float64                `json:"price"`
	Metadata map[string]interface{} `json:"metadata"`
}

// UpdateInventoryRequest represents the request body for updating inventory (admin only)
type UpdateInventoryRequest struct {
	ProductID string `json:"product_id" binding:"required"` // UUID as string
	Quantity  int    `json:"quantity" binding:"required,min=0"`
}


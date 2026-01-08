package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"oms/backend/api/v1/helpers"
	"oms/backend/api/v1/types"
	"oms/backend/core/model"
	"oms/backend/core/services"
)

// OrderController handles order-related HTTP requests
type OrderController struct {
	orderService services.OrderService
}

// NewOrderController creates a new OrderController
func NewOrderController(orderService services.OrderService) *OrderController {
	return &OrderController{
		orderService: orderService,
	}
}

// CreateOrder handles POST /api/v1/orders
// User ID is extracted from JWT token in middleware (req.user)
// Admin users cannot place orders
func (oc *OrderController) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract user_id and role from JWT token (set by auth middleware)
	userID := getUserIDFromContext(ctx)
	role := getUserRoleFromContext(ctx)
	if userID == 0 {
		helpers.WriteErrorResponse(w, http.StatusUnauthorized, "unauthorized", "User ID not found in context")
		return
	}

	// Admin cannot place orders
	if role == "admin" {
		helpers.WriteErrorResponse(w, http.StatusForbidden, "forbidden", "Admin users cannot place orders")
		return
	}

	var req types.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}

	// Validate request
	if req.Quantity <= 0 {
		helpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid_request", "Quantity must be greater than 0")
		return
	}

	// Parse product_id as UUID
	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		helpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid_request", "Invalid product ID format")
		return
	}

	// Convert shipping address to metadata JSONB
	metadata := model.JSONB{}
	if req.ShippingAddress != nil {
		metadata = model.JSONB(req.ShippingAddress)
	}

	// Call orderService.CreateOrder with user_id from JWT token
	order, err := oc.orderService.CreateOrder(ctx, userID, productID, req.Quantity, metadata)
	if err != nil {
		// Check for specific error types
		errMsg := err.Error()
		if len(errMsg) > 20 && errMsg[:20] == "insufficient inventory" {
			helpers.WriteErrorResponse(w, http.StatusBadRequest, "insufficient_inventory", errMsg)
			return
		}
		helpers.WriteErrorResponse(w, http.StatusInternalServerError, "internal_error", "Failed to create order: "+errMsg)
		return
	}

	helpers.WriteJSONResponse(w, http.StatusCreated, types.CreateOrderResponse{
		OrderID:       order.ID.String(),
		CurrentStatus: string(order.CurrentStatus),
		Message:       "Order placed successfully",
	})
}

// UpdateOrderStatus handles PATCH /api/v1/orders/{orderId}
func (oc *OrderController) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	orderID := vars["orderId"]

	// Extract user_id from JWT token (set by auth middleware)
	userID := getUserIDFromContext(ctx)

	var req types.UpdateOrderStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}

	// Validate orderID format (UUID)
	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		helpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid_request", "Invalid order ID format")
		return
	}

	// Validate status
	newStatus := model.OrderStatus(req.CurrentStatus)
	if !isValidOrderStatus(newStatus) {
		helpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid_request", "Invalid order status")
		return
	}

	// Fetch current order to get previous status and validate
	currentOrder, err := oc.orderService.GetOrderByID(ctx, orderUUID)
	if err != nil {
		helpers.WriteErrorResponse(w, http.StatusNotFound, "not_found", "Order not found")
		return
	}
	previousStatus := string(currentOrder.CurrentStatus)

	// Check role-based restrictions
	role := getUserRoleFromContext(ctx)
	if role == "admin" {
		// Admin can only update to SHIPPED or DELIVERED
		if newStatus != model.OrderStatusShipped && newStatus != model.OrderStatusDelivered {
			helpers.WriteErrorResponse(w, http.StatusForbidden, "forbidden", "Admin can only update status to SHIPPED or DELIVERED")
			return
		}
	} else {
		// Regular users can only cancel ORDERED orders (cannot update to SHIPPED or DELIVERED)
		if newStatus != model.OrderStatusCancelled {
			helpers.WriteErrorResponse(w, http.StatusForbidden, "forbidden", "Regular users can only cancel ORDERED orders. Only admin can update status to SHIPPED or DELIVERED")
			return
		}
		// Verify the order is in ORDERED status before allowing cancellation
		if currentOrder.CurrentStatus != model.OrderStatusOrdered {
			helpers.WriteErrorResponse(w, http.StatusForbidden, "forbidden", "Only ORDERED orders can be cancelled by regular users")
			return
		}
	}

	// Call orderService.UpdateOrderStatus with user_id from JWT token
	order, err := oc.orderService.UpdateOrderStatus(ctx, orderUUID, newStatus, userID)
	if err != nil {
		// Check for FSM validation error (409 Conflict)
		errMsg := err.Error()
		if len(errMsg) > 15 && errMsg[:15] == "invalid transition" {
			helpers.WriteErrorResponse(w, http.StatusConflict, "invalid_transition", errMsg)
			return
		}
		if errMsg == "order not found" || len(errMsg) > 13 && errMsg[:13] == "order not found" {
			helpers.WriteErrorResponse(w, http.StatusNotFound, "not_found", errMsg)
			return
		}
		helpers.WriteErrorResponse(w, http.StatusInternalServerError, "internal_error", "Failed to update order status")
		return
	}

	helpers.WriteJSONResponse(w, http.StatusOK, types.UpdateOrderStatusResponse{
		OrderID:        order.ID.String(),
		PreviousStatus: previousStatus,
		CurrentStatus:  string(order.CurrentStatus),
		UpdatedBy:      userID, // Retrieved from JWT token via context
		UpdatedAt:      &order.UpdatedAt,
	})
}

// GetOrders handles GET /api/v1/orders
// Admin sees all orders, regular users see only their own orders
func (oc *OrderController) GetOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract user_id and role from JWT token (set by auth middleware)
	userID := getUserIDFromContext(ctx)
	role := getUserRoleFromContext(ctx)
	if userID == 0 {
		helpers.WriteErrorResponse(w, http.StatusUnauthorized, "unauthorized", "User ID not found in context")
		return
	}

	var orders []*model.Order
	var err error

	// Admin sees all orders, regular users see only their own
	if role == "admin" {
		orders, err = oc.orderService.GetAllOrders(ctx)
	} else {
		orders, err = oc.orderService.GetOrdersByUserID(ctx, userID)
	}

	if err != nil {
		helpers.WriteErrorResponse(w, http.StatusInternalServerError, "internal_error", "Failed to fetch orders")
		return
	}

	// Convert to response format
	orderResponses := make([]types.OrderResponse, len(orders))
	for i, order := range orders {
		// Handle metadata conversion safely
		metadata := map[string]interface{}{}
		if order.Metadata != nil {
			metadata = map[string]interface{}(order.Metadata)
		}
		
		orderResponses[i] = types.OrderResponse{
			ID:            order.ID.String(),
			UserID:        order.UserID,
			ProductID:     order.ProductID.String(),
			Quantity:      order.Quantity,
			CurrentStatus: string(order.CurrentStatus),
			Metadata:      metadata,
			CreatedAt:     order.CreatedAt,
			UpdatedAt:     order.UpdatedAt,
		}
	}

	helpers.WriteJSONResponse(w, http.StatusOK, orderResponses)
}

// GetOrderHistory handles GET /api/v1/orders/{orderId}/history - Get order state change history
func (oc *OrderController) GetOrderHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	orderID := vars["orderId"]

	// Extract user_id from JWT token (set by auth middleware)
	userID := getUserIDFromContext(ctx)
	if userID == 0 {
		helpers.WriteErrorResponse(w, http.StatusUnauthorized, "unauthorized", "User ID not found in context")
		return
	}

	// Validate orderID format (UUID)
	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		helpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid_request", "Invalid order ID format")
		return
	}

	// Verify order exists and belongs to user
	order, err := oc.orderService.GetOrderByID(ctx, orderUUID)
	if err != nil {
		helpers.WriteErrorResponse(w, http.StatusNotFound, "not_found", "Order not found")
		return
	}

	// Verify order belongs to user (admin can access any order)
	role := getUserRoleFromContext(ctx)
	if role != "admin" && order.UserID != userID {
		helpers.WriteErrorResponse(w, http.StatusForbidden, "forbidden", "You don't have access to this order")
		return
	}

	// Get order history
	history, err := oc.orderService.GetOrderHistory(ctx, orderUUID)
	if err != nil {
		helpers.WriteErrorResponse(w, http.StatusInternalServerError, "internal_error", "Failed to fetch order history")
		return
	}

	// Convert to response format
	historyResponses := make([]types.OrderHistoryResponse, len(history))
	for i, log := range history {
		historyResponses[i] = types.OrderHistoryResponse{
			OrderID:        log.OrderID.String(),
			PreviousStatus: string(log.PreviousStatus),
			NewStatus:      string(log.NewStatus),
			UpdatedBy:      log.UpdatedBy,
			UpdatedAt:      log.UpdatedAt,
		}
	}

	helpers.WriteJSONResponse(w, http.StatusOK, historyResponses)
}

// Helper function to validate order status
func isValidOrderStatus(status model.OrderStatus) bool {
	validStatuses := []model.OrderStatus{
		model.OrderStatusOrdered,
		model.OrderStatusShipped,
		model.OrderStatusDelivered,
		model.OrderStatusCancelled,
	}
	for _, valid := range validStatuses {
		if status == valid {
			return true
		}
	}
	return false
}

// Helper function to extract user ID from context (set by auth middleware)
func getUserIDFromContext(ctx context.Context) int {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		return 0
	}
	return userID
}

// Helper function to extract user role from context (set by auth middleware)
func getUserRoleFromContext(ctx context.Context) string {
	role, ok := ctx.Value("user_role").(string)
	if !ok {
		return "user" // Default to regular user
	}
	return role
}


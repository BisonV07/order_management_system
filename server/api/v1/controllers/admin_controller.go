package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"oms/server/api/v1/helpers"
	apitypes "oms/server/api/v1/types"
	"oms/server/core/model"
	"oms/server/core/types"
)

// AdminController handles admin-only operations (products and inventory management)
type AdminController struct {
	productStore   types.ProductStore
	inventoryStore types.InventoryStore
}

// NewAdminController creates a new AdminController
func NewAdminController(productStore types.ProductStore, inventoryStore types.InventoryStore) *AdminController {
	return &AdminController{
		productStore:   productStore,
		inventoryStore: inventoryStore,
	}
}

// CreateProduct handles POST /api/v1/admin/products - Create a new product (admin only)
func (ac *AdminController) CreateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Verify admin role
	role := getUserRoleFromContext(ctx)
	if role != "admin" {
		helpers.WriteErrorResponse(w, http.StatusForbidden, "forbidden", "Admin access required")
		return
	}

	var req apitypes.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}

	// Validate request
	if req.SKU == "" || req.Name == "" {
		helpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid_request", "SKU and Name are required")
		return
	}
	if req.Price < 0 {
		helpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid_request", "Price cannot be negative")
		return
	}

	// Create product
	product := &model.Product{
		SKU:      req.SKU,
		Name:     req.Name,
		Price:    req.Price,
		Metadata: model.JSONB(req.Metadata),
	}

	err := ac.productStore.Create(ctx, product)
	if err != nil {
		// Check for duplicate SKU error
		if err.Error() == "duplicate key value violates unique constraint" || 
		   err.Error() == "UNIQUE constraint failed" {
			helpers.WriteErrorResponse(w, http.StatusConflict, "conflict", "Product with this SKU already exists")
			return
		}
		helpers.WriteErrorResponse(w, http.StatusInternalServerError, "internal_error", "Failed to create product: "+err.Error())
		return
	}

	// Create initial inventory entry (default to 0)
	_ = ac.inventoryStore.UpdateQuantity(ctx, product.ID, 0) // Ignore error if inventory already exists

	helpers.WriteJSONResponse(w, http.StatusCreated, apitypes.CreateProductResponse{
		ProductID: product.ID.String(),
		Message:   "Product created successfully",
	})
}

// UpdateProduct handles PUT /api/v1/admin/products/{productId} - Update a product (admin only)
func (ac *AdminController) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	productIDStr := vars["productId"]

	// Verify admin role
	role := getUserRoleFromContext(ctx)
	if role != "admin" {
		helpers.WriteErrorResponse(w, http.StatusForbidden, "forbidden", "Admin access required")
		return
	}

	// Parse product ID
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		helpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid_request", "Invalid product ID format")
		return
	}

	var req apitypes.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}

	// Get existing product to preserve fields not being updated
	existingProduct, err := ac.productStore.GetByID(ctx, productID)
	if err != nil {
		helpers.WriteErrorResponse(w, http.StatusNotFound, "not_found", "Product not found")
		return
	}

	// Update fields if provided
	if req.SKU != "" {
		existingProduct.SKU = req.SKU
	}
	if req.Name != "" {
		existingProduct.Name = req.Name
	}
	if req.Price >= 0 {
		existingProduct.Price = req.Price
	}
	if req.Metadata != nil {
		existingProduct.Metadata = model.JSONB(req.Metadata)
	}

	err = ac.productStore.Update(ctx, productID, existingProduct)
	if err != nil {
		helpers.WriteErrorResponse(w, http.StatusInternalServerError, "internal_error", "Failed to update product: "+err.Error())
		return
	}

	helpers.WriteJSONResponse(w, http.StatusOK, apitypes.ProductResponse{
		ID:        existingProduct.ID.String(),
		SKU:       existingProduct.SKU,
		Name:      existingProduct.Name,
		Price:     existingProduct.Price,
		Metadata:  map[string]interface{}(existingProduct.Metadata),
		CreatedAt: existingProduct.CreatedAt,
		UpdatedAt: existingProduct.UpdatedAt,
	})
}

// UpdateInventory handles PUT /api/v1/admin/inventory - Update inventory quantity (admin only)
func (ac *AdminController) UpdateInventory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Verify admin role
	role := getUserRoleFromContext(ctx)
	if role != "admin" {
		helpers.WriteErrorResponse(w, http.StatusForbidden, "forbidden", "Admin access required")
		return
	}

	var req apitypes.UpdateInventoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}

	// Validate request
	if req.Quantity < 0 {
		helpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid_request", "Quantity cannot be negative")
		return
	}

	// Parse product ID
	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		helpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid_request", "Invalid product ID format")
		return
	}

	// Verify product exists
	_, err = ac.productStore.GetByID(ctx, productID)
	if err != nil {
		helpers.WriteErrorResponse(w, http.StatusNotFound, "not_found", "Product not found")
		return
	}

	// Update inventory
	err = ac.inventoryStore.UpdateQuantity(ctx, productID, req.Quantity)
	if err != nil {
		helpers.WriteErrorResponse(w, http.StatusInternalServerError, "internal_error", "Failed to update inventory: "+err.Error())
		return
	}

	helpers.WriteJSONResponse(w, http.StatusOK, apitypes.UpdateInventoryResponse{
		ProductID: productID.String(),
		Quantity:  req.Quantity,
		Message:   "Inventory updated successfully",
	})
}

// DeleteProduct handles DELETE /api/v1/admin/products/{productId} - Delete a product (admin only)
func (ac *AdminController) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	productIDStr := vars["productId"]

	// Verify admin role
	role := getUserRoleFromContext(ctx)
	if role != "admin" {
		helpers.WriteErrorResponse(w, http.StatusForbidden, "forbidden", "Admin access required")
		return
	}

	// Parse product ID
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		helpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid_request", "Invalid product ID format")
		return
	}

	// Verify product exists before deleting
	_, err = ac.productStore.GetByID(ctx, productID)
	if err != nil {
		helpers.WriteErrorResponse(w, http.StatusNotFound, "not_found", "Product not found")
		return
	}

	// Delete product (soft delete)
	err = ac.productStore.Delete(ctx, productID)
	if err != nil {
		helpers.WriteErrorResponse(w, http.StatusInternalServerError, "internal_error", "Failed to delete product: "+err.Error())
		return
	}

	helpers.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message":   "Product deleted successfully",
		"product_id": productID.String(),
	})
}


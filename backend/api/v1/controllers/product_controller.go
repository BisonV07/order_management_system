package controllers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"oms/backend/api/v1/helpers"
	"oms/backend/core/services"
)

// ProductController handles product-related HTTP requests
type ProductController struct {
	productService services.ProductService
}

// NewProductController creates a new ProductController
func NewProductController(productService services.ProductService) *ProductController {
	return &ProductController{
		productService: productService,
	}
}

// GetProducts handles GET /api/v1/products
func (pc *ProductController) GetProducts(w http.ResponseWriter, r *http.Request) {
	_ = r.Context() // TODO: Use context for logging/request tracking

	// TODO: Call productService.GetAll(ctx)
	// ctx := r.Context()
	// products, err := pc.productService.GetAll(ctx)
	// if err != nil {
	//     helpers.WriteErrorResponse(w, http.StatusInternalServerError, "internal_error", "Failed to fetch products")
	//     return
	// }

	// Placeholder response - return empty array for now
	helpers.WriteJSONResponse(w, http.StatusOK, []interface{}{})
}

// GetProduct handles GET /api/v1/products/{productId}
func (pc *ProductController) GetProduct(w http.ResponseWriter, r *http.Request) {
	_ = r.Context() // TODO: Use context for logging/request tracking
	vars := mux.Vars(r)
	productIDStr := vars["productId"]

	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		helpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid_request", "Invalid product ID format")
		return
	}

	// TODO: Call productService.GetByID(ctx, productID)
	// ctx := r.Context()
	// product, err := pc.productService.GetByID(ctx, productID)
	// if err != nil {
	//     helpers.WriteErrorResponse(w, http.StatusNotFound, "not_found", "Product not found")
	//     return
	// }

	// Placeholder response
	helpers.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"id":   productID,
		"name": "Placeholder Product",
	})
}


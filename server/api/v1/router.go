package v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"oms/server/api/v1/controllers"
	"oms/server/api/v1/helpers"
	"oms/server/core/services"
	"oms/server/core/types"
	"oms/server/middleware"
)

// SetupRouter configures and returns the API v1 router
func SetupRouter(orderService services.OrderService) *mux.Router {
	return SetupRouterWithInventory(orderService, nil, nil)
}

// SetupRouterWithInventory configures and returns the API v1 router with inventory store
func SetupRouterWithInventory(orderService services.OrderService, inventoryStore types.InventoryStore, userStore types.UserStore) *mux.Router {
	return SetupRouterWithStores(orderService, inventoryStore, userStore, nil)
}

// SetupRouterWithStores configures and returns the API v1 router with all stores
func SetupRouterWithStores(orderService services.OrderService, inventoryStore types.InventoryStore, userStore types.UserStore, productStore types.ProductStore) *mux.Router {
	return SetupRouterWithStoresAndDB(orderService, inventoryStore, userStore, productStore, nil)
}

// SetupRouterWithStoresAndDB configures and returns the API v1 router with all stores and database
func SetupRouterWithStoresAndDB(orderService services.OrderService, inventoryStore types.InventoryStore, userStore types.UserStore, productStore types.ProductStore, db *gorm.DB) *mux.Router {
	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()

	// Apply middleware (CORS must be first)
	router.Use(middleware.CORSMiddleware)
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.PanicRecoveryMiddleware)
	router.Use(middleware.AuthMiddleware) // JWT authentication

	// Initialize controllers
	authController := controllers.NewAuthController(userStore)
	orderController := controllers.NewOrderController(orderService)
	
	// Initialize admin controller if stores are available
	var adminController *controllers.AdminController
	if productStore != nil && inventoryStore != nil {
		adminController = controllers.NewAdminController(productStore, inventoryStore)
	}
	
	// Initialize metrics controller if database is available
	var metricsController *controllers.MetricsController
	if db != nil {
		metricsController = controllers.NewMetricsController(db)
	}

	// Auth routes (no auth required)
	router.HandleFunc("/auth/login", authController.Login).Methods("POST")
	router.HandleFunc("/auth/signup", authController.Signup).Methods("POST")

	// Order routes (require authentication)
	router.HandleFunc("/orders", orderController.CreateOrder).Methods("POST")
	router.HandleFunc("/orders", orderController.GetOrders).Methods("GET")
	router.HandleFunc("/orders/{orderId}", orderController.UpdateOrderStatus).Methods("PATCH")
	router.HandleFunc("/orders/{orderId}/history", orderController.GetOrderHistory).Methods("GET")
	
	// Product routes (public, no auth required for GET)
	router.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		
		// If product store is available, fetch real products from database
		if productStore != nil {
			products, err := productStore.GetAll(ctx)
			if err != nil {
				helpers.WriteErrorResponse(w, http.StatusInternalServerError, "internal_error", "Failed to fetch products: "+err.Error())
				return
			}
			
			// Convert to response format with inventory
			productResponses := make([]map[string]interface{}, len(products))
			for i, product := range products {
				// Handle metadata conversion safely
				metadata := map[string]interface{}{}
				if product.Metadata != nil {
					metadata = map[string]interface{}(product.Metadata)
				}
				
				productResponses[i] = map[string]interface{}{
					"id":       product.ID.String(),
					"sku":      product.SKU,
					"name":     product.Name,
					"price":    product.Price,
					"metadata": metadata,
				}
				
				// Fetch inventory if inventory store is available
				if inventoryStore != nil {
					inv, err := inventoryStore.GetByProductID(ctx, product.ID)
					if err == nil {
						productResponses[i]["inventory"] = inv.Quantity
					} else {
						productResponses[i]["inventory"] = 0
					}
				} else {
					productResponses[i]["inventory"] = 0
				}
			}
			
			helpers.WriteJSONResponse(w, http.StatusOK, productResponses)
			return
		}
		
		// Fallback to empty array if no product store
		helpers.WriteJSONResponse(w, http.StatusOK, []interface{}{})
	}).Methods("GET")
	router.HandleFunc("/products/{productId}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		helpers.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
			"id":   vars["productId"],
			"name": "Placeholder Product",
		})
	}).Methods("GET")

	// Admin routes (require admin role)
	if adminController != nil {
		router.HandleFunc("/admin/products", adminController.CreateProduct).Methods("POST")
		router.HandleFunc("/admin/products/{productId}", adminController.UpdateProduct).Methods("PUT")
		router.HandleFunc("/admin/products/{productId}", adminController.DeleteProduct).Methods("DELETE")
		router.HandleFunc("/admin/inventory", adminController.UpdateInventory).Methods("PUT")
	}

	// Metrics routes (require admin role)
	if metricsController != nil {
		router.HandleFunc("/admin/metrics", metricsController.GetMetrics).Methods("GET")
		router.HandleFunc("/admin/metrics/docker", metricsController.GetDockerMetrics).Methods("GET")
		router.HandleFunc("/admin/metrics/postgresql", metricsController.GetPostgreSQLMetrics).Methods("GET")
	} else {
		// Register routes with error handler if metrics controller is not available
		router.HandleFunc("/admin/metrics", func(w http.ResponseWriter, r *http.Request) {
			helpers.WriteErrorResponse(w, http.StatusServiceUnavailable, "service_unavailable", "Metrics service is not available. Database connection may be missing.")
		}).Methods("GET")
		router.HandleFunc("/admin/metrics/docker", func(w http.ResponseWriter, r *http.Request) {
			helpers.WriteErrorResponse(w, http.StatusServiceUnavailable, "service_unavailable", "Metrics service is not available. Database connection may be missing.")
		}).Methods("GET")
		router.HandleFunc("/admin/metrics/postgresql", func(w http.ResponseWriter, r *http.Request) {
			helpers.WriteErrorResponse(w, http.StatusServiceUnavailable, "service_unavailable", "Metrics service is not available. Database connection may be missing.")
		}).Methods("GET")
	}

	// Health check endpoint (no auth required)
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	return router
}


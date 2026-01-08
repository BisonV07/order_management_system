package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"oms/backend/api/v1"
	"oms/backend/config"
	"oms/backend/database"
	"oms/backend/datastore"
	"oms/backend/core/fsm"
	"oms/backend/core/model"
	"oms/backend/core/services"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func main() {
	apiFlag := flag.Bool("api", false, "Start the API server")
	migrateFlag := flag.Bool("migrate", false, "Run database migrations")
	port := flag.String("port", "8080", "Port to run the API server on")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if *migrateFlag {
		runMigrations(db)
		return
	}

	if *apiFlag {
		startAPIServer(*port, db)
		return
	}

	flag.Usage()
	os.Exit(1)
}

func runMigrations(db *gorm.DB) {
	fmt.Println("Running database migrations...")
	
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	
	// Seed admin user if it doesn't exist
	seedAdminUser(db)
	
	// Seed dummy products with 0 stock
	seedDummyProducts(db)
}

func seedAdminUser(db *gorm.DB) {
	var adminUser model.User
	result := db.Where("username = ?", "admin").First(&adminUser)
	if result.Error != nil {
		// Admin user doesn't exist, create it
		hashedPassword, err := model.HashPassword("1234")
		if err != nil {
			log.Printf("Warning: Failed to hash admin password: %v", err)
			return
		}
		
		adminUser = model.User{
			Username: "admin",
			Password: hashedPassword,
			Role:     model.UserRoleAdmin,
		}
		if err := db.Create(&adminUser).Error; err != nil {
			log.Printf("Warning: Failed to create admin user: %v", err)
		} else {
			log.Println("✅ Admin user created (username: admin, password: 1234)")
		}
	} else {
		log.Println("✅ Admin user already exists")
	}
}

func seedDummyProducts(db *gorm.DB) {
	fmt.Println("Seeding dummy products...")
	
	dummyProducts := []model.Product{
		{
			ID:       uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			SKU:      "PROD-001",
			Name:     "Laptop Computer",
			Price:    1299.99,
			Metadata: model.JSONB{"brand": "TechCorp", "color": "Silver"},
		},
		{
			ID:       uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
			SKU:      "PROD-002",
			Name:     "Wireless Mouse",
			Price:    29.99,
			Metadata: model.JSONB{"brand": "TechCorp", "color": "Black"},
		},
		{
			ID:       uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"),
			SKU:      "PROD-003",
			Name:     "Mechanical Keyboard",
			Price:    149.99,
			Metadata: model.JSONB{"brand": "TechCorp", "color": "RGB"},
		},
		{
			ID:       uuid.MustParse("550e8400-e29b-41d4-a716-446655440003"),
			SKU:      "PROD-004",
			Name:     "USB-C Hub",
			Price:    79.99,
			Metadata: model.JSONB{"brand": "TechCorp", "ports": 7},
		},
		{
			ID:       uuid.MustParse("550e8400-e29b-41d4-a716-446655440004"),
			SKU:      "PROD-005",
			Name:     "Monitor Stand",
			Price:    89.99,
			Metadata: model.JSONB{"brand": "TechCorp", "material": "Aluminum"},
		},
	}
	
	createdCount := 0
	for _, product := range dummyProducts {
		var existingProduct model.Product
		result := db.Where("id = ? OR sku = ?", product.ID, product.SKU).First(&existingProduct)
		
		if result.Error != nil {
			// Product doesn't exist, create it
			if err := db.Create(&product).Error; err != nil {
				log.Printf("Warning: Failed to create product %s: %v", product.SKU, err)
			} else {
				createdCount++
				// Create inventory entry with 0 stock
				inventory := model.Inventory{
					ProductID: product.ID,
					Quantity:  0,
				}
				if err := db.Create(&inventory).Error; err != nil {
					log.Printf("Warning: Failed to create inventory for product %s: %v", product.SKU, err)
				}
			}
		}
	}
	
	if createdCount > 0 {
		log.Printf("✅ Created %d new dummy products with 0 stock", createdCount)
	} else {
		log.Println("✅ All dummy products already exist")
	}
}

func startAPIServer(port string, db *gorm.DB) {
	fmt.Printf("Starting API server on port %s...\n", port)
	
	// Seed admin user and products if they don't exist (idempotent)
	seedAdminUser(db)
	seedDummyProducts(db)
	
	// Initialize real database stores
	orderStore := datastore.NewOrderStore(db)
	inventoryStore := datastore.NewInventoryStore(db)
	orderStateLogStore := datastore.NewOrderStateLogStore(db)
	userStore := datastore.NewUserStore(db)
	productStore := datastore.NewProductStore(db)
	fsmValidator := fsm.NewValidator()
	
	orderService := services.NewOrderService(
		orderStore,
		inventoryStore,
		orderStateLogStore,
		fsmValidator,
	)
	
	// Setup router with all stores including product store and database for admin features and metrics
	router := v1.SetupRouterWithStoresAndDB(orderService, inventoryStore, userStore, productStore, db)
	
	// Start server - bind to all interfaces to ensure browser connectivity
	serverAddr := "0.0.0.0:" + port
	fmt.Printf("Server listening on http://localhost:%s\n", port)
	fmt.Printf("Health check: http://localhost:%s/api/v1/health\n", port)
	fmt.Printf("Products: http://localhost:%s/api/v1/products\n", port)
	fmt.Printf("✅ Connected to PostgreSQL database\n")
	
	if err := http.ListenAndServe(serverAddr, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}


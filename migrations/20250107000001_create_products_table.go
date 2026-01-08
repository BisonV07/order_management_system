package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upCreateProductsTable, downCreateProductsTable)
}

func upCreateProductsTable(tx *sql.Tx) error {
	query := `
	CREATE TABLE IF NOT EXISTS products (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		sku VARCHAR(255) UNIQUE NOT NULL,
		name VARCHAR(255) NOT NULL,
		price DECIMAL(10,2) NOT NULL,
		metadata JSONB,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_products_sku ON products(sku);
	CREATE INDEX IF NOT EXISTS idx_products_deleted_at ON products(deleted_at);
	`
	_, err := tx.Exec(query)
	return err
}

func downCreateProductsTable(tx *sql.Tx) error {
	query := `DROP TABLE IF EXISTS products;`
	_, err := tx.Exec(query)
	return err
}


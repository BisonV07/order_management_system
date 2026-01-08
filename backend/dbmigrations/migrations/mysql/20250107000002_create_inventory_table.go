package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upCreateInventoryTable, downCreateInventoryTable)
}

func upCreateInventoryTable(tx *sql.Tx) error {
	query := `
	CREATE TABLE IF NOT EXISTS inventory (
		product_id UUID PRIMARY KEY,
		quantity INTEGER NOT NULL DEFAULT 0,
		last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT fk_inventory_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_inventory_product_id ON inventory(product_id);
	`
	_, err := tx.Exec(query)
	return err
}

func downCreateInventoryTable(tx *sql.Tx) error {
	query := `DROP TABLE IF EXISTS inventory;`
	_, err := tx.Exec(query)
	return err
}


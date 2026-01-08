package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upCreateOrdersTable, downCreateOrdersTable)
}

func upCreateOrdersTable(tx *sql.Tx) error {
	query := `
	CREATE TYPE order_status AS ENUM ('ORDERED', 'SHIPPED', 'DELIVERED', 'CANCELLED');

	CREATE TABLE IF NOT EXISTS orders (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id INTEGER NOT NULL,
		product_id UUID NOT NULL,
		quantity INTEGER NOT NULL,
		current_status VARCHAR(50) NOT NULL DEFAULT 'ORDERED',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT fk_orders_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE RESTRICT
	);

	CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
	CREATE INDEX IF NOT EXISTS idx_orders_product_id ON orders(product_id);
	CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(current_status);
	CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at);
	`
	_, err := tx.Exec(query)
	return err
}

func downCreateOrdersTable(tx *sql.Tx) error {
	query := `
	DROP TABLE IF EXISTS orders;
	DROP TYPE IF EXISTS order_status;
	`
	_, err := tx.Exec(query)
	return err
}


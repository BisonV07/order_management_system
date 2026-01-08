package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upAddMetadataToOrders, downAddMetadataToOrders)
}

func upAddMetadataToOrders(tx *sql.Tx) error {
	query := `
	ALTER TABLE orders ADD COLUMN IF NOT EXISTS metadata JSONB;
	`
	_, err := tx.Exec(query)
	return err
}

func downAddMetadataToOrders(tx *sql.Tx) error {
	query := `
	ALTER TABLE orders DROP COLUMN IF EXISTS metadata;
	`
	_, err := tx.Exec(query)
	return err
}


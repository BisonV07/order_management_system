package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upCreateOrderStateLogsTable, downCreateOrderStateLogsTable)
}

func upCreateOrderStateLogsTable(tx *sql.Tx) error {
	query := `
	CREATE TABLE IF NOT EXISTS order_state_logs (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		order_id UUID NOT NULL,
		previous_status VARCHAR(50),
		new_status VARCHAR(50) NOT NULL,
		updated_by INTEGER NOT NULL,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT fk_order_state_logs_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_order_state_logs_order_id ON order_state_logs(order_id);
	CREATE INDEX IF NOT EXISTS idx_order_state_logs_updated_at ON order_state_logs(updated_at);
	`
	_, err := tx.Exec(query)
	return err
}

func downCreateOrderStateLogsTable(tx *sql.Tx) error {
	query := `DROP TABLE IF EXISTS order_state_logs;`
	_, err := tx.Exec(query)
	return err
}


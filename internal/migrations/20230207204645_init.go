package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upInit, downInit)
}

func upInit(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`select 1`)
	return err
}

func downInit(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.Exec(`select 1`)
	return err
}

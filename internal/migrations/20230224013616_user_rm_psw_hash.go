package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upUserRmPswHash, downUserRmPswHash)
}

func upUserRmPswHash(tx *sql.Tx) error {
	q := `alter table if exists user_auth drop column if exists psw_hash;`

	_, err := tx.Exec(q)
	return err
}

func downUserRmPswHash(tx *sql.Tx) error {
	q := `alter table if exists user_auth add column if not exists psw_hash bytea not null default ''::bytea;`

	_, err := tx.Exec(q)
	return err
}

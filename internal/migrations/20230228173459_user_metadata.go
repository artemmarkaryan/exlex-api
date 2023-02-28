package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upUserMetadata, downUserMetadata)
}

func upUserMetadata(tx *sql.Tx) error {
	q := `
create table if not exists user_metadata
(
    user_uuid  uuid unique references user_auth (id),
    fullName   text,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);
`

	_, err := tx.Exec(q)
	return err
}

func downUserMetadata(tx *sql.Tx) error {
	q := `drop table if exists user_metadata;`

	_, err := tx.Exec(q)
	return err
}

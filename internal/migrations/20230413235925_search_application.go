package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upSearchApplication, downSearchApplication)
}

func upSearchApplication(tx *sql.Tx) error {
	q := `
	create table if not exists search_application (
		id       	uuid default gen_random_uuid() primary key,
		search_id 	uuid 		not null references search(id),
		user_id 	uuid 		not null references user_auth(id),
		created_at 	timestamp 	not null default current_timestamp,
		comment text,
		unique (search_id, user_id)
	);
`

	_, err := tx.Exec(q)
	return err
}

func downSearchApplication(tx *sql.Tx) error {
	q := `drop table if exists search_application`
	_, err := tx.Exec(q)
	return err
}

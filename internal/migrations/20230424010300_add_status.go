package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upAddStatus, downAddStatus)
}

func upAddStatus(tx *sql.Tx) (err error) {
	q := `
alter table if exists 
	search_application 
add column if not exists 
	status text not null default 'new'::text;
`

	if _, err = tx.Exec(q); err != nil {
		return err
	}

	q = `
alter table if exists
	search
add column if not exists
	status text not null default 'new'::text;
`

	if _, err = tx.Exec(q); err != nil {
		return err
	}

	return nil
}

func downAddStatus(tx *sql.Tx) error {
	q := `
alter table if exists search_application drop column if exists status;
alter table if exists search drop column if exists status;
`
	_, err := tx.Exec(q)
	return err
}

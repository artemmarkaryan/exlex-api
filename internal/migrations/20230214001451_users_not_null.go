package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upUsersNotNull, downUsersNotNull)
}

func upUsersNotNull(tx *sql.Tx) error {
	q := `
alter table if exists user_auth  
    alter column email set not null,
	alter column psw_hash set not null;
`
	_, err := tx.Exec(q)
	return err
}

func downUsersNotNull(tx *sql.Tx) error {
	q := `
alter table if exists user_auth  
    alter column email drop not null,
	alter column psw_hash drop not null;
`
	_, err := tx.Exec(q)
	return err
}

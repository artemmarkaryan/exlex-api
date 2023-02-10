package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upUsers, downUsers)
}

func upUsers(tx *sql.Tx) error {
	q := `
create table if not exists user_auth
(
    id       uuid default gen_random_uuid() primary key,
    email    text,
    psw_hash bytea
);

create index if not exists user_auth_email_idx on user_auth (email);

create table if not exists user_role
(
    user_id uuid references user_auth(id),
    role text
)
`
	_, err := tx.Exec(q)
	return err
}

func downUsers(tx *sql.Tx) error {
	q := `
drop table if exists user_auth;
drop table if exists user_role;
`
	_, err := tx.Exec(q)
	return err
}

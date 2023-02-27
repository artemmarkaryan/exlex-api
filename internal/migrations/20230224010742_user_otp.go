package migrations

import (
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upUserOtp, downUserOtp)
}

func upUserOtp(tx *sql.Tx) error {
	q := `
create table if not exists user_otp
(
    user_uuid  uuid references user_auth (id) not null,
    otp        text                           not null,
    created_at timestamp                      not null default current_timestamp
);

create index if not exists user_otp_uuid on user_otp(user_uuid);
`
	_, err := tx.Exec(q)
	return err
}

func downUserOtp(tx *sql.Tx) error {
	q := `drop table if exists user_otp`

	_, err := tx.Exec(q)
	return err
}

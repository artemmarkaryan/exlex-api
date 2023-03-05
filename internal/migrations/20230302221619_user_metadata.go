package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upUserMetadata2, downUserMetadata2)
}

func upUserMetadata2(tx *sql.Tx) error {
	q := `
drop table if exists user_education;
drop table if exists user_metadata;

create table if not exists customer_metadata
(
    user_uuid uuid unique not null references user_auth (id),
    name      text
);

create index if not exists customer_metadata_idx on customer_metadata(user_uuid);

create table if not exists executor_metadata
(
    user_uuid        uuid unique not null references user_auth (id),
    education        text references education_type(id),
    full_name        text,
    experience_years integer
);

create index if not exists executor_metadata_idx on executor_metadata(user_uuid);
`
	_, err := tx.Exec(q)
	return err
}

func downUserMetadata2(tx *sql.Tx) error {
	q := `drop table if exists user_metadata;`
	_, err := tx.Exec(q)
	return err
}

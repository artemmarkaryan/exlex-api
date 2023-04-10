package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upSearchEducationFix, downSearchEducationFix)
}

func upSearchEducationFix(tx *sql.Tx) error {
	q := `
drop table if exists search_requirement_education;
create table if not exists search_requirement_education
(
    search_uuid uuid not null references search (id),
    education   text not null references education_type (id),
    unique (search_uuid, education)
)`

	_, err := tx.Exec(q)
	return err
}

func downSearchEducationFix(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}

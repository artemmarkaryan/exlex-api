package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upEducationType, downEducationType)
}

func upEducationType(tx *sql.Tx) error {
	q := `
create table if not exists education_type
(
    id    text primary key,
    title text unique not null
);

insert into education_type (id, title)
values ('secondary_vocational', 'Среднее профессиональное'),
       ('incomplete_higher', 'Неоконченное высшее'),
       ('higher', 'Высшее');

create table if not exists user_education
(
    user_uuid uuid references user_auth (id),
    education text references education_type (id),
    unique (user_uuid, education)
);
`
	_, err := tx.Exec(q)
	return err
}

func downEducationType(tx *sql.Tx) error {
	q := `
drop table if exists user_education;
drop table if exists education_type;
`
	_, err := tx.Exec(q)
	return err
}

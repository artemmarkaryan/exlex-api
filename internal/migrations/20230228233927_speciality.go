package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upSpeciality, downSpeciality)
}

func upSpeciality(tx *sql.Tx) error {
	q := `
create table if not exists speciality
(
    id    text primary key,
    title text not null
);

insert into speciality (id, title)
values ('public_law', 'Публичное право'),
       ('civil_law', 'Частное право');

create table if not exists user_speciality
(
    user_uuid  uuid references user_auth (id),
    speciality text references speciality (id),
    unique (user_uuid, speciality)
);
`
	_, err := tx.Exec(q)
	return err
}

func downSpeciality(tx *sql.Tx) error {
	q := `
drop table if exists user_speciality;
drop table if exists speciality;
`
	_, err := tx.Exec(q)
	return err
}

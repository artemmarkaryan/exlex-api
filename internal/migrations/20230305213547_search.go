package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upSearch, downSearch)
}

func upSearch(tx *sql.Tx) error {
	q := `
create table if not exists search
(

    id                       uuid primary key   default gen_random_uuid(),
    creator                  uuid      not null references user_auth (id),
    name                     text      not null,
    description              text      not null,
    price                    float4    not null,
    required_work_experience int       not null,
    created_at               timestamp not null default current_timestamp,
    deadline                 date               default null
);

create table if not exists search_requirement_speciality
(
    search_uuid uuid not null references search (id),
    speciality  text not null references speciality (id),
	unique (search_uuid, speciality)
);

create table if not exists search_requirement_education
(
    search_uuid uuid not null references search (id),
    education   text not null references education_type (id),
    unique (search_uuid)
)`
	_, err := tx.Exec(q)
	return err
}

func downSearch(tx *sql.Tx) error {
	q := `
drop table if exists search_requirement_education;
drop table if exists search_requirement_speciality;
drop table if exists search;
`
	_, err := tx.Exec(q)
	return err
}

package schema

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type SearchApplication struct {
	ID        uuid.UUID      `db:"id"`
	SearchID  uuid.UUID      `db:"search_id"`
	UserID    uuid.UUID      `db:"user_id"` // applicant
	CreatedAt time.Time      `db:"created_at"`
	Comment   sql.NullString `db:"comment"`
	Status    string         `db:"status"`
}

func (a SearchApplication) TableName() string { return "search_application" }

type SearchApplicationRaw struct {
	ApplicationID uuid.UUID      `db:"application_id"`
	UserID        uuid.UUID      `db:"user_id"`
	CreatedAt     time.Time      `db:"created_at"`
	Comment       sql.NullString `db:"comment"`
	Education     string         `db:"education"`
	FullName      string         `db:"full_name"`
	Experience    int            `db:"experience_years"`
	Specialities  []byte         `db:"speciality"`
	Status        string         `db:"status"`
}

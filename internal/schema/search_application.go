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
}

func (a SearchApplication) TableName() string { return "search_application" }

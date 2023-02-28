package schema

import "github.com/google/uuid"

type UserAuth struct {
	ID    uuid.UUID `db:"id"`
	Email string    `db:"email"`
}

func (u UserAuth) TableName() string { return "user_auth" }

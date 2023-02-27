package schema

import (
	"time"

	"github.com/google/uuid"
)

type UserAuth struct {
	ID    uuid.UUID `db:"id"`
	Email string    `db:"email"`
}

func (u UserAuth) TableName() string { return "user_auth" }

type UserOTP struct {
	userUUID  uuid.UUID `db:"user_uuid"`
	otp       string    `db:"otp"`
	createdAt time.Time `db:"created_at"`
}

func (u UserOTP) TableName() string { return "user_otp" }

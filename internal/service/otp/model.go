package otp

import (
	"time"

	"github.com/google/uuid"
)

type UserOTP struct {
	userUUID  uuid.UUID `db:"user_uuid"`
	otp       string    `db:"otp"`
	createdAt time.Time `db:"created_at"`
}

func (u UserOTP) TableName() string { return "user_otp" }

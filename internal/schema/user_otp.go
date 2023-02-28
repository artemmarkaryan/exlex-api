package schema

import (
	"time"

	"github.com/google/uuid"
)

type UserOTP struct {
	UserUUID  uuid.UUID `db:"user_uuid"`
	OTP       string    `db:"otp"`
	CreatedAt time.Time `db:"created_at"`
}

func (u UserOTP) TableName() string { return "user_otp" }

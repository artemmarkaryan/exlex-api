package schema

import (
	"errors"
	"time"

	"github.com/artemmarkaryan/exlex-backend/graph/model"
	"github.com/google/uuid"
)

var ErrUnknownRole = errors.New("unknown role")

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

type Role string

const RoleCustomer = "customer"
const RoleExecutor = "executor"

func MapRole(role model.Role) (Role, error) {
	switch role {
	case model.RoleCustomer:
		return RoleCustomer, nil
	case model.RoleExecutor:
		return RoleExecutor, nil
	default:
		return "", ErrUnknownRole
	}
}

type UserRole struct {
	userID uuid.UUID `db:"user_id"`
	role   Role      `db:"role"`
}

func (u UserRole) TableName() string { return "user_role" }

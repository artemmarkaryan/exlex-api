package schema

import (
	"errors"

	"github.com/artemmarkaryan/exlex-backend/graph/model"
	"github.com/google/uuid"
)

var ErrUnknownRole = errors.New("unknown role")

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
	UserID uuid.UUID `db:"user_id"`
	Role   Role      `db:"role"`
}

func (u UserRole) TableName() string { return "user_role" }

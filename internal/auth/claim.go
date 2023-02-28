package auth

import (
	"github.com/artemmarkaryan/exlex-backend/internal/schema"
	"github.com/cristalhq/jwt/v5"
	"github.com/google/uuid"
)

type Claim struct {
	jwt.RegisteredClaims

	UserID uuid.UUID
	Email  string
	Role   schema.Role
}

func MakeClaim(userID uuid.UUID, email string, role schema.Role) Claim {
	return Claim{
		UserID: userID,
		Email:  email,
		Role:   role,
	}
}

package auth

import (
	"github.com/cristalhq/jwt/v5"
	"github.com/google/uuid"
)

type Claim struct {
	jwt.RegisteredClaims

	UserID uuid.UUID
	Email  string
}

func MakeClaim(userID uuid.UUID, email string) Claim {
	return Claim{
		UserID: userID,
		Email:  email,
	}
}

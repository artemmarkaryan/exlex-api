package graph

import (
	"github.com/artemmarkaryan/exlex-backend/internal/service/authentication"
	user_profile "github.com/artemmarkaryan/exlex-backend/internal/service/user-profile"
)

//go:generate go run github.com/99designs/gqlgen generate

type serviceContainer interface {
	Authentication() authentication.Service
	UserProfile() user_profile.Service
}

type Resolver struct {
	ServiceContainer serviceContainer
}

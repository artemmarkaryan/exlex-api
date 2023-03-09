package graph

import (
	"github.com/artemmarkaryan/exlex-backend/internal/service/authentication"
	"github.com/artemmarkaryan/exlex-backend/internal/service/search"
	user_profile "github.com/artemmarkaryan/exlex-backend/internal/service/user-profile"
)

//go:generate go run github.com/99designs/gqlgen generate

type serviceContainer interface {
	Authentication() authentication.Service
	UserProfile() user_profile.Service
	Search() search.Service
}

type Resolver struct {
	ServiceContainer serviceContainer
}

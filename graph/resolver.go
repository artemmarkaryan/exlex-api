package graph

import "github.com/artemmarkaryan/exlex-backend/internal/service/authentication"

//go:generate go run github.com/99designs/gqlgen generate

type serviceContainer interface {
	Authentication() authentication.Service
}

type Resolver struct {
	ServiceContainer serviceContainer
}

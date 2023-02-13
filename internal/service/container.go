package service

import (
	"context"

	"github.com/artemmarkaryan/exlex-backend/internal/service/authentication"
)

type Container struct {
	authentication authentication.Service
}

func MakeContainer(ctx context.Context) (c Container, err error) {
	c.authentication = authentication.Make(ctx)

	return
}

func (c Container) Authentication() authentication.Service { return c.authentication }

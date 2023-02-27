package service

import (
	"context"
	"errors"
	"os"

	"github.com/artemmarkaryan/exlex-backend/internal/service/authentication"
	"github.com/artemmarkaryan/exlex-backend/internal/service/otp"
)

type Container struct {
	authentication authentication.Service
	otp            otp.Service
}

func MakeContainer(ctx context.Context) (c Container, err error) {
	{
		token := os.Getenv("UNIONE_TOKEN")
		if token == "" {
			err = errors.New("UNIONE_TOKEN is empty")
			return Container{}, err
		}

		c.otp = otp.Make(otp.Config{UnioneToken: token})
	}

	c.authentication = authentication.Make(ctx, c)

	return
}

func (c Container) Authentication() authentication.Service { return c.authentication }
func (c Container) OTP() otp.Service                       { return c.otp }

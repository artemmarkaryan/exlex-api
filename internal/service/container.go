package service

import (
	"context"
	"errors"
	"os"

	"github.com/artemmarkaryan/exlex-backend/internal/service/authentication"
	"github.com/artemmarkaryan/exlex-backend/internal/service/otp"
	"github.com/artemmarkaryan/exlex-backend/pkg/tokenizer"
	"github.com/cristalhq/jwt/v5"
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

	{
		cfg := authentication.Config{
			TokenizerConfig: tokenizer.Config{
				Algorithm: jwt.HS512,
				SecretKey: os.Getenv("UNIONE_TOKEN"),
			},
		}

		c.authentication, err = authentication.Make(ctx, cfg, c)
		if err != nil {
			return
		}
	}

	return
}

func (c Container) Authentication() authentication.Service { return c.authentication }
func (c Container) OTP() otp.Service                       { return c.otp }

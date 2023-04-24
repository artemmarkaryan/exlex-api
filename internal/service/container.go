package service

import (
	"context"
	"errors"
	"os"

	"github.com/artemmarkaryan/exlex-backend/internal/service/authentication"
	"github.com/artemmarkaryan/exlex-backend/internal/service/otp"
	"github.com/artemmarkaryan/exlex-backend/internal/service/search"
	user_profile "github.com/artemmarkaryan/exlex-backend/internal/service/user-profile"
)

type Container struct {
	authentication authentication.Service
	otp            otp.Service
	userProfile    user_profile.Service
	search         search.Service
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
	c.userProfile = user_profile.MakeService()
	c.search = search.MakeService(c)

	return
}

func (c Container) Authentication() authentication.Service { return c.authentication }
func (c Container) OTP() otp.Service                       { return c.otp }
func (c Container) UserProfile() user_profile.Service      { return c.userProfile }
func (c Container) Search() search.Service                 { return c.search }

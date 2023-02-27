package authentication

import (
	"context"

	"github.com/artemmarkaryan/exlex-backend/internal/service/otp"
	"github.com/google/uuid"
)

type serviceContainer interface {
	OTP() otp.Service
}

type otpService interface {
	GenerateAndSend(ctx context.Context, uuid uuid.UUID, email string) error
}

type Service struct {
	repo
	otpService otpService
}

func Make(_ context.Context, container serviceContainer) (s Service) {
	s.repo = repo{}
	s.otpService = container.OTP()
	return
}

func (s Service) RequestOTP(ctx context.Context, email string) (err error) {
	id, err := s.getOrCreateUser(ctx, email)
	if err != nil {
		return err
	}

	if err = s.otpService.GenerateAndSend(ctx, id, email); err != nil {
		return err
	}

	return
}

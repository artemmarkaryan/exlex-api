package authentication

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/artemmarkaryan/exlex-backend/internal/auth"
	"github.com/artemmarkaryan/exlex-backend/internal/schema"
	"github.com/artemmarkaryan/exlex-backend/internal/service/otp"
	"github.com/artemmarkaryan/exlex-backend/pkg/database"
	"github.com/artemmarkaryan/exlex-backend/pkg/tokenizer"
	"github.com/cristalhq/jwt/v5"
	"github.com/google/uuid"
)

var ErrAlreadyExists = errors.New("user already exists")
var ErrNotFound = errors.New("user not found")

type serviceContainer interface {
	OTP() otp.Service
}

type otpService interface {
	GenerateAndSend(ctx context.Context, uuid uuid.UUID, email string, debug bool) error
	Verify(ctx context.Context, email string, o string) error
}

type tokenFactory interface {
	NewToken(claims any) (*jwt.Token, error)
	VerifyToken(token *jwt.Token) (err error)
	Parse(raw []byte, claims *any) error
}

type Service struct {
	repo
	otpService   otpService
	tokenFactory tokenFactory
}

func Make(ctx context.Context, container serviceContainer) (s Service) {
	s.repo = repo{}
	s.otpService = container.OTP()
	s.tokenFactory = tokenizer.FromContext(ctx)
	return
}

func (s Service) Signup(ctx context.Context, email string, role schema.Role, debug bool) (err error) {
	id, err := s.createUser(ctx, email, role)
	if err != nil {
		return err
	}

	if err = s.otpService.GenerateAndSend(ctx, id, email, debug); err != nil {
		return err
	}

	return
}

func (s Service) Login(ctx context.Context, email string, debug bool) (err error) {
	id, err := s.getUser(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return err
	}

	if err = s.otpService.GenerateAndSend(ctx, id, email, debug); err != nil {
		return err
	}

	return
}

func (s Service) VerifyOTP(ctx context.Context, email string, o string) (token string, err error) {
	if err = s.otpService.
		Verify(ctx, email, o); err != nil {
		return
	}

	q := sq.
		Select("*").
		From(new(schema.UserAuth).TableName()).
		Where(sq.Eq{"email": email})

	user, err := database.GetX[schema.UserAuth](ctx, q)
	if err != nil {
		return
	}

	q = sq.
		Select("role").
		From(new(schema.UserRole).TableName()).
		Where(sq.Eq{"user_id": user.ID})
	role, err := database.GetX[schema.UserRole](ctx, q)
	if err != nil {
		return
	}

	claim := auth.MakeClaim(user.ID, user.Email, role.Role)
	t, err := s.tokenFactory.NewToken(claim)
	if err != nil {
		return
	}

	return t.String(), nil
}

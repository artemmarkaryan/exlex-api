package otp

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"

	"github.com/artemmarkaryan/exlex-backend/pkg/telegram"
	"github.com/artemmarkaryan/exlex-backend/pkg/unione"
	"github.com/google/uuid"
)

var (
	ErrNotFound = errors.New("no OTP found")
	ErrWrongOTP = errors.New("wrong OTP")
)

type Config struct {
	UnioneToken string
}

type Service struct {
	repo
	unione unione.Client
}

func Make(cfg Config) (s Service) {
	s.unione = unione.NewClient(cfg.UnioneToken)
	s.repo = repo{}
	return
}

func (s Service) GenerateAndSend(ctx context.Context, id uuid.UUID, email string, debug bool) error {
	code := generateOTP()

	if err := s.repo.insert(ctx, id, code); err != nil {
		return err
	}

	if debug {
		log.Printf("[debug] confirmation otp: %s", code)
		telegram.Report(ctx, fmt.Sprintf("code: %s", code))
		return nil
	}

	err := s.unione.Send(
		unione.Message{
			Recipients: []unione.Recipient{{Email: email}},
			Body: unione.Body{
				Plaintext: "код подтверждения: " + code,
			},
			Subject:   "Код подтверждения ExLex",
			FromEmail: "no-reply@exlex.site",
			FromName:  "служебная почта ExLex",
		},
	)
	if err != nil {
		log.Printf("unione: %v", err)
		return errors.New("sending message failed")
	}

	return nil
}

func generateOTP() string {
	randomNumber := rand.Intn(9000) + 1000
	return strconv.Itoa(randomNumber)
}

func (s Service) Verify(ctx context.Context, email string, input string) error {
	otp, err := s.repo.get(ctx, email)
	if err == sql.ErrNoRows {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if input != otp {
		return ErrWrongOTP
	}

	if err = s.repo.delete(ctx, email); err != nil {
		log.Printf("error deleting otp: %v", err)
	}

	return nil
}

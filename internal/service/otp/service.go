package otp

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"strconv"

	"github.com/artemmarkaryan/exlex-backend/pkg/unione"
	"github.com/google/uuid"
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

func (s Service) GenerateAndSend(ctx context.Context, id uuid.UUID, email string) error {
	var code string
	{
		randomNumber := rand.Intn(10_000) + 10_000
		randomString := strconv.Itoa(randomNumber)
		code = randomString[1:]
	}

	err := s.repo.insert(ctx, id, code)
	if err != nil {
		return err
	}

	err = s.unione.Send(
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
		log.Print("unione: " + err.Error())
		return errors.New("sending message failed")
	}

	return nil
}

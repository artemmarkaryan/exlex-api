package user_profile

import (
	"context"

	"github.com/artemmarkaryan/exlex-backend/internal/schema"
)

type Service struct {
	repo repo
}

func MakeService() (s Service) {
	s.repo = repo{}
	return
}

func (s Service) Specialities(ctx context.Context) ([]schema.Speciality, error) {
	return s.repo.specialities(ctx)
}

func (s Service) EducationTypes(ctx context.Context) ([]schema.EducationType, error) {
	return s.repo.educationTypes(ctx)
}

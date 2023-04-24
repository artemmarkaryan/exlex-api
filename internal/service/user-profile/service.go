package user_profile

import (
	"context"
	"errors"

	"github.com/artemmarkaryan/exlex-backend/internal/schema"
	"github.com/google/uuid"
)

var ErrUserNotFound = errors.New("user not found")

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

func (s Service) UpdateCustomerProfile(ctx context.Context, d UpdateCustomerProfileData) error {
	if d.UserUUID == *new(uuid.UUID) {
		return ErrNoUserIDProvided
	}

	return s.repo.setCustomer(ctx, d)
}

func (s Service) UpdateExecutorProfile(ctx context.Context, d UpdateExecutorProfileData) error {
	if d.UserUUID == *new(uuid.UUID) {
		return ErrNoUserIDProvided
	}

	if err := s.repo.setExecutor(ctx, d); err != nil {
		return err
	}

	return nil
}

func (s Service) GetCustomerProfile(ctx context.Context, id uuid.UUID) (CustomerProfile, error) {
	cp, err := s.repo.getCustomerProfile(ctx, id)
	if err != nil {
		return CustomerProfile{}, err
	}

	return CustomerProfile{
		FullName: cp.Name,
	}, nil
}

func (s Service) GetUserEmail(ctx context.Context, id uuid.UUID) (string, error) {
	return s.repo.getEmail(ctx, id)
}

func (s Service) GetExecutorProfile(ctx context.Context, id uuid.UUID) (ExecutorProfile, error) {
	ep, err := s.repo.getExecutorProfile(ctx, id)
	if err != nil {
		return ExecutorProfile{}, err
	}

	return ExecutorProfile{
		FullName:        ep.FullName,
		WorkExperience:  ep.ExperienceYears,
		EducationTypeID: ep.Education,
		Specialization:  ep.Specialities,
	}, nil
}

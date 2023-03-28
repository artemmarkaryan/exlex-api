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

type UpdateUserProfileData struct {
	UserUUID uuid.UUID
}

type UpdateExecutorProfileData struct {
	UpdateUserProfileData
	FullName        *string
	ExperienceYears *int
	Specialities    []string
	Education       *string
}

type UpdateCustomerProfileData struct {
	UpdateUserProfileData
	Name *string
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

type CustomerProfile struct {
	FullName string
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

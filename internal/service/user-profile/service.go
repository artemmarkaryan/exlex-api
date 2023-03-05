package user_profile

import (
	"context"

	"github.com/artemmarkaryan/exlex-backend/internal/schema"
	"github.com/google/uuid"
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

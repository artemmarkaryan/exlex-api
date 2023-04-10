package search

import (
	"context"
	"time"

	"github.com/artemmarkaryan/exlex-backend/internal/schema"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type Service struct {
	repo repo
}

func MakeService() (s Service) {
	s.repo = repo{}
	return
}

type CreateSearch struct {
	Creator                uuid.UUID
	Name                   string
	Description            string
	Price                  float64
	RequiredWorkExperience int
	Deadline               *time.Time
	RequiredSpecialities   []string
	RequiredEducation      []string
}

type Search struct {
	ID                     uuid.UUID
	Creator                uuid.UUID
	Name                   string
	Description            string
	Price                  float64
	RequiredWorkExperience int
	Deadline               *time.Time
	CreatedAt              time.Time
	RequiredSpecialities   []string
	RequiredEducation      []string
}

func (s Service) Create(ctx context.Context, d CreateSearch) (uuid.UUID, error) {
	id, err := s.repo.create(ctx, d)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (s Service) Delete(ctx context.Context, user, search uuid.UUID) error {
	if err := s.repo.checkCreator(ctx, user, search); err != nil {
		return nil
	}

	if err := s.repo.delete(ctx, search); err != nil {
		return err
	}

	return nil
}

func (s Service) Get(ctx context.Context, user, search uuid.UUID) (se Search, err error) {
	if err = s.repo.checkCreator(ctx, user, search); err != nil {
		return
	}

	dbo, err := s.repo.get(ctx, search)
	if err != nil {
		return
	}

	return dboToSearch(dbo), nil
}

func (s Service) List(ctx context.Context, user uuid.UUID) (se []Search, err error) {
	dbos, err := s.repo.list(ctx, user)
	if err != nil {
		return
	}

	se = lo.Map(dbos, func(o schema.SearchFullData, _ int) Search {
		return dboToSearch(o)
	})

	return
}

func dboToSearch(dbo schema.SearchFullData) (s Search) {
	return Search{
		ID:                     dbo.ID,
		Name:                   dbo.Name,
		Description:            dbo.Description,
		Price:                  dbo.Price,
		RequiredWorkExperience: dbo.RequiredWorkExperience,
		Deadline:               dbo.Deadline,
		CreatedAt:              dbo.CreatedAt,
		RequiredSpecialities:   dbo.Speciality,
		RequiredEducation:      dbo.Education,
	}
}

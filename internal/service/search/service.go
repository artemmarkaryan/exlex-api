package search

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/artemmarkaryan/exlex-backend/internal/schema"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
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
	Name                   string
	Description            string
	Price                  float64
	RequiredWorkExperience int
	Deadline               *time.Time
	CreatedAt              time.Time
	RequiredSpecialities   []string
	RequiredEducation      []string
}

func (s *Search) fillFromRaw(dbo schema.SearchFullDataRaw) (err error) {
	var educations []string
	err = json.Unmarshal(dbo.Education, &educations)
	if err != nil {
		return
	}

	var specialities []string
	err = json.Unmarshal(dbo.Speciality, &specialities)
	if err != nil {
		return
	}

	filter := func(s []string) []string {
		slices.Sort(s)
		s = slices.Compact(s)
		s = lo.Filter(s, func(obj string, _ int) bool { return obj != "" })
		return s
	}

	*s = Search{
		ID:                     dbo.ID,
		Name:                   dbo.Name,
		Description:            dbo.Description,
		Price:                  dbo.Price,
		RequiredWorkExperience: dbo.RequiredWorkExperience,
		Deadline:               dbo.Deadline,
		CreatedAt:              dbo.CreatedAt,
		RequiredEducation:      filter(educations),
		RequiredSpecialities:   filter(specialities),
	}

	return nil
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

	err = se.fillFromRaw(dbo)
	return
}

func (s Service) ListByAuthor(ctx context.Context, user uuid.UUID) (se []Search, err error) {
	dbos, err := s.repo.listByAuthor(ctx, user)
	if err != nil {
		return
	}

	se = make([]Search, len(dbos))
	for i := range dbos {
		err = se[i].fillFromRaw(dbos[i])
		if err != nil {
			return nil, err
		}
	}

	return
}

// todo: apply fiters, based on requirements. now all are available
func (s Service) ListAvailableForApplication(ctx context.Context, user uuid.UUID) (se []Search, err error) {
	dbos, err := s.repo.listAvailableForApplication(ctx, user)
	if err != nil {
		return
	}

	se = make([]Search, len(dbos))
	for i := range dbos {
		err = se[i].fillFromRaw(dbos[i])
		if err != nil {
			return nil, err
		}
	}

	return
}

type SearchApplicationRequest struct {
	SearchID uuid.UUID
	UserID   uuid.UUID
	Comment  *string
}

func (s Service) Apply(ctx context.Context, r SearchApplicationRequest) (applicationID uuid.UUID, err error) {
	applicationID, err = s.repo.apply(ctx, r)
	if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		err = ErrApplicationAlreadyExists
	}

	return
}

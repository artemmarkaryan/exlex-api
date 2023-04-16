package search

import (
	"context"
	"strings"

	"github.com/google/uuid"
)

type Service struct {
	repo repo
}

func MakeService() (s Service) {
	s.repo = repo{}
	return
}

func (s Service) Create(ctx context.Context, d CreateSearchRequest) (uuid.UUID, error) {
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

func (s Service) Apply(ctx context.Context, r SearchApplicationRequest) (applicationID uuid.UUID, err error) {
	applicationID, err = s.repo.apply(ctx, r)
	if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		err = ErrApplicationAlreadyExists
	}

	return
}

func (s Service) ListApplicants(ctx context.Context, r ListApplicantsRequest) ([]Application, error) {
	err := s.repo.checkCreator(ctx, r.UserID, r.SearchID)
	if err != nil {
		return nil, err
	}

	dbos, err := s.repo.listApplications(ctx, r.SearchID)
	if err != nil {
		return nil, err
	}

	apps := make([]Application, len(dbos))
	for i := range dbos {
		err := apps[i].FillFromRaw(dbos[i])
		if err != nil {
			return nil, err
		}
	}

	return apps, nil
}

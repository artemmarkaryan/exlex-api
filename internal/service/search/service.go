package search

import (
	"context"
	"strings"

	user_profile "github.com/artemmarkaryan/exlex-backend/internal/service/user-profile"
	"github.com/google/uuid"
)

type ServiceContainer interface {
	UserProfile() user_profile.Service
}

type Service struct {
	repo      repo
	container ServiceContainer
}

func MakeService(sc ServiceContainer) (s Service) {
	s.repo = repo{}
	s.container = sc
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

func (s Service) ApproveApplication(ctx context.Context, r ApproveApplicationRequest) error {
	return s.repo.approveApplication(ctx, r)
}

func (s Service) GetSearchAssignee(ctx context.Context, searchID uuid.UUID) (a Assignee, err error) {
	assigneeID, err := s.repo.getAssigneeID(ctx, searchID)
	if err != nil {
		return
	}

	profile, err := s.container.UserProfile().GetExecutorProfile(ctx, assigneeID)
	if err != nil {
		return
	}

	email, err := s.container.UserProfile().GetUserEmail(ctx, assigneeID)
	if err != nil {
		return
	}

	return Assignee{
		Email:          email,
		FullName:       profile.FullName,
		WorkExperience: profile.WorkExperience,
		Education:      profile.EducationTypeID,
		Specialities:   profile.Specialization,
	}, nil
}

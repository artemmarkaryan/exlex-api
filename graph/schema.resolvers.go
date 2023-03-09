package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.24

import (
	"context"
	"fmt"
	"net/mail"
	"time"

	"github.com/artemmarkaryan/exlex-backend/graph/model"
	"github.com/artemmarkaryan/exlex-backend/internal/auth"
	"github.com/artemmarkaryan/exlex-backend/internal/schema"
	"github.com/artemmarkaryan/exlex-backend/internal/service/search"
	user_profile "github.com/artemmarkaryan/exlex-backend/internal/service/user-profile"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, data model.LoginData) (ok model.Ok, err error) {
	err = r.ServiceContainer.Authentication().Login(ctx, data.Email, data.Debug)
	ok.Ok = err == nil
	return
}

// Signup is the resolver for the signup field.
func (r *mutationResolver) Signup(ctx context.Context, data model.SignupData) (ok model.Ok, err error) {
	role, err := schema.MapRole(data.Role)
	if err != nil {
		return
	}

	err = r.ServiceContainer.Authentication().Signup(ctx, data.Email, role, data.Debug)
	ok.Ok = err == nil
	return
}

// VerifyOtp is the resolver for the verifyOTP field.
func (r *mutationResolver) VerifyOtp(ctx context.Context, email string, otp string) (token model.Token, err error) {
	if _, err = mail.ParseAddress(email); err != nil {
		return
	}

	t, err := r.ServiceContainer.Authentication().VerifyOTP(ctx, email, otp)
	if err != nil {
		return
	}

	token.Access = t
	return
}

// SetCustomerProfile is the resolver for the SetCustomerProfile field.
func (r *mutationResolver) SetCustomerProfile(ctx context.Context, data model.SetCustomerProfileData) (model.Ok, error) {
	c, err := auth.FromContext(ctx)
	if err != nil {
		return model.Ok{}, err
	}

	updateData := user_profile.UpdateCustomerProfileData{}
	updateData.UserUUID = c.UserID
	updateData.Name = data.FullName

	err = r.ServiceContainer.
		UserProfile().
		UpdateCustomerProfile(ctx, updateData)

	return model.Ok{Ok: err == nil}, err
}

// SetExecutorProfile is the resolver for the SetExecutorProfile field.
func (r *mutationResolver) SetExecutorProfile(ctx context.Context, data model.SetExecutorProfileData) (model.Ok, error) {
	c, err := auth.FromContext(ctx)
	if err != nil {
		return model.Ok{}, err
	}

	updateData := user_profile.UpdateExecutorProfileData{}
	updateData.UserUUID = c.UserID
	updateData.FullName = data.FullName
	updateData.Education = data.EducationTypeID
	updateData.Specialities = data.Specialization
	updateData.ExperienceYears = data.WorkExperience

	err = r.ServiceContainer.UserProfile().UpdateExecutorProfile(ctx, updateData)
	return model.Ok{Ok: err == nil}, err
}

// CreateSearch is the resolver for the createSearch field.
func (r *mutationResolver) CreateSearch(ctx context.Context, data model.CreateSearchInput) (string, error) {
	claims, err := auth.FromContext(ctx)
	if err != nil {
		return "", err
	}

	d := search.CreateSearch{
		Creator:     claims.UserID,
		Name:        data.Title,
		Description: data.Description,
		Price:       data.Price,
	}

	if data.Deadline != nil {
		dl := *data.Deadline
		dt := time.Date(dl.Year, time.Month(dl.Month), dl.Day, 0, 0, 0, 0, time.UTC)
		d.Deadline = &dt
	}

	if req := data.Requirements; req != nil {
		d.RequiredSpecialities = req.Speciality
		d.RequiredEducation = req.EducationType
		d.RequiredWorkExperience = req.WorkExperience
	}

	id, err := r.ServiceContainer.Search().Create(ctx, d)
	return id.String(), err
}

// DeleteSearch is the resolver for the deleteSearch field.
func (r *mutationResolver) DeleteSearch(ctx context.Context, id string) (model.Ok, error) {
	claims, err := auth.FromContext(ctx)
	if err != nil {
		return model.Ok{Ok: false}, err
	}

	searchID, err := uuid.Parse(id)
	if err != nil {
		return model.Ok{Ok: false}, ErrBadUUID
	}

	err = r.ServiceContainer.Search().Delete(ctx, claims.UserID, searchID)

	return model.Ok{Ok: err == nil}, err
}

// Live is the resolver for the live field.
func (r *queryResolver) Live(ctx context.Context) (bool, error) {
	return true, nil
}

// Authenticated is the resolver for the authenticated field.
func (r *queryResolver) Authenticated(ctx context.Context) (bool, error) {
	_, err := auth.FromContext(ctx)
	return err != auth.ErrUnauthenticated, err
}

// Specialities is the resolver for the specialities field.
func (r *queryResolver) Specialities(ctx context.Context) ([]model.Speciality, error) {
	s, err := r.ServiceContainer.UserProfile().Specialities(ctx)
	if err != nil {
		return nil, err
	}

	f := func(o schema.Speciality, _ int) model.Speciality {
		return model.Speciality{
			ID:    o.ID,
			Title: o.Title,
		}
	}

	return lo.Map(s, f), nil
}

// EducationTypes is the resolver for the educationTypes field.
func (r *queryResolver) EducationTypes(ctx context.Context) ([]model.EducationType, error) {
	s, err := r.ServiceContainer.UserProfile().EducationTypes(ctx)
	if err != nil {
		return nil, err
	}

	f := func(o schema.EducationType, _ int) model.EducationType {
		return model.EducationType{
			ID:    o.ID,
			Title: o.Title,
		}
	}

	return lo.Map(s, f), nil
}

// Search is the resolver for the search field.
func (r *queryResolver) Search(ctx context.Context, id string) (model.Search, error) {
	claims, err := auth.FromContext(ctx)
	if err != nil {
		return model.Search{}, err
	}

	searchID, err := uuid.Parse(id)
	if err != nil {
		return model.Search{}, ErrBadUUID
	}

	s, err := r.ServiceContainer.Search().Get(ctx, claims.UserID, searchID)

	return model.Search{
		Title:       s.Name,
		Description: s.Description,
		Price:       s.Price,
		Deadline: &model.Date{
			Year:  s.Deadline.Year(),
			Month: int(s.Deadline.Month()),
			Day:   s.Deadline.Day(),
		},
		Requirements: &model.SearchRequirements{
			Speciality:     s.RequiredSpecialities,
			EducationType:  s.RequiredEducation,
			WorkExperience: s.RequiredWorkExperience,
		},
	}, nil
}

// Searches is the resolver for the searches field.
func (r *queryResolver) Searches(ctx context.Context) ([]*model.Search, error) {
	panic(fmt.Errorf("not implemented: Searches - searches"))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

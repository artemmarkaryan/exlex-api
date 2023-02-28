package user_profile

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/artemmarkaryan/exlex-backend/internal/schema"
	"github.com/artemmarkaryan/exlex-backend/pkg/database"
)

type repo struct{}

func (repo) specialities(ctx context.Context) ([]schema.Speciality, error) {
	q := sq.
		Select("*").
		From(new(schema.Speciality).TableName())

	return database.SelectX[schema.Speciality](ctx, q)
}

func (repo) educationTypes(ctx context.Context) ([]schema.EducationType, error) {
	q := sq.
		Select("*").
		From(new(schema.EducationType).TableName())

	return database.SelectX[schema.EducationType](ctx, q)
}

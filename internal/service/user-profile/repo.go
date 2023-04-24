package user_profile

import (
	"context"
	"database/sql"
	"encoding/json"

	sq "github.com/Masterminds/squirrel"
	"github.com/artemmarkaryan/exlex-backend/internal/schema"
	"github.com/artemmarkaryan/exlex-backend/pkg/database"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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

func (repo) setCustomer(ctx context.Context, d UpdateCustomerProfileData) error {
	if d.Name == nil {
		return nil
	}

	q := sq.
		Insert(new(schema.CustomerMetadata).TableName()).
		SetMap(map[string]any{"user_uuid": d.UserUUID, "name": *d.Name}).
		Suffix("ON CONFLICT (user_uuid) DO UPDATE SET name = EXCLUDED.name")

	_, err := database.InsertX(ctx, q)
	return err
}

// todo: optimise this
func (repo) setExecutor(ctx context.Context, d UpdateExecutorProfileData) error {
	executorMetadataExists := false

	return database.Tx(
		ctx, database.DefaultTxOptions(),
		func(tx *sqlx.Tx) (errTx error) {
			q := sq.
				Select("true").
				From(new(schema.ExecutorMetadata).TableName()).
				Where(sq.Eq{"user_uuid": d.UserUUID})

			executorMetadataExists, errTx = database.GetTxX[bool](ctx, tx, q)
			if errTx == sql.ErrNoRows {
				executorMetadataExists = false
				return nil
			}

			return
		},

		// update if exists
		func(tx *sqlx.Tx) error {
			if !executorMetadataExists {
				return nil
			}

			m := make(map[string]any)

			if d.FullName != nil {
				m["full_name"] = d.FullName
			}

			if d.ExperienceYears != nil {
				m["experience_years"] = d.ExperienceYears
			}

			if d.Education != nil {
				m["education"] = d.Education
			}

			if len(m) == 0 {
				return nil
			}

			q := sq.
				Update(new(schema.ExecutorMetadata).TableName()).
				Where(sq.Eq{"user_uuid": d.UserUUID}).
				SetMap(m)

			_, err := database.UpdateTxX(ctx, tx, q)
			return err
		},

		// insert if not existing
		func(tx *sqlx.Tx) error {
			if executorMetadataExists {
				return nil
			}
			m := make(map[string]any)

			if d.FullName != nil {
				m["full_name"] = *d.FullName
			}

			if d.ExperienceYears != nil {
				m["experience_years"] = *d.ExperienceYears
			}

			if d.Education != nil {
				m["education"] = *d.Education
			}

			if len(m) == 0 {
				return nil
			}

			m["user_uuid"] = d.UserUUID

			q := sq.
				Insert(new(schema.ExecutorMetadata).TableName()).
				SetMap(m)

			_, err := database.InsertTxX(ctx, tx, q)
			return err
		},

		// insert specialities
		func(tx *sqlx.Tx) error {
			if len(d.Specialities) == 0 {
				return nil
			}

			q := sq.
				Insert(new(schema.UserSpeciality).TableName()).
				Columns("user_uuid", "speciality")

			for _, s := range d.Specialities {
				q = q.Values(d.UserUUID, s)
			}

			q = q.Suffix(`on conflict (user_uuid, speciality) do nothing`)

			_, err := database.InsertTxX(ctx, tx, q)
			return err
		},

		func(tx *sqlx.Tx) error {
			q := sq.
				Delete(new(schema.UserSpeciality).TableName()).
				Where(
					sq.And{
						sq.Eq{"user_uuid": d.UserUUID},
						sq.NotEq{"speciality": d.Specialities},
					},
				)

			_, err := database.DeleteTxX(ctx, tx, q)
			return err
		},
	)
}

func (repo) getCustomerProfile(ctx context.Context, id uuid.UUID) (c schema.CustomerMetadata, err error) {
	err = database.Tx(ctx, database.DefaultTxOptions(),
		func(tx *sqlx.Tx) error {
			q := sq.
				Select(`true`).
				From(new(schema.UserAuth).TableName()).
				Where(sq.Eq{"id": id})

			_, err := database.GetTxX[bool](ctx, tx, q)
			if err == sql.ErrNoRows {
				return ErrUserNotFound
			}
			if err != nil {
				return err
			}
			return nil
		},
		func(tx *sqlx.Tx) error {
			q := sq.
				Select(`*`).
				From(new(schema.CustomerMetadata).TableName()).
				Where(sq.Eq{"user_uuid": id})

			c, err = database.GetTxX[schema.CustomerMetadata](ctx, tx, q)
			if err == sql.ErrNoRows {
				c.UserUUID = id
			}

			return nil
		},
	)

	return
}

func (repo) getExecutorProfile(ctx context.Context, id uuid.UUID) (m schema.FullExecutorMetadata, err error) {
	err = database.Tx(ctx, database.DefaultTxOptions(),
		func(tx *sqlx.Tx) error {
			q := sq.
				Select(`true`).
				From(new(schema.UserAuth).TableName()).
				Where(sq.Eq{"id": id})

			_, err := database.GetTxX[bool](ctx, tx, q)
			if err == sql.ErrNoRows {
				return ErrUserNotFound
			}
			if err != nil {
				return err
			}
			return nil
		},

		func(tx *sqlx.Tx) error {
			q := sq.
				Select(`*`).
				From(new(schema.ExecutorMetadata).TableName()).
				Where(sq.Eq{"user_uuid": id})

			m.ExecutorMetadata, err = database.GetTxX[schema.ExecutorMetadata](ctx, tx, q)
			if err == sql.ErrNoRows {
				m.UserUUID = id
			}

			return nil
		},

		func(tx *sqlx.Tx) error {
			q := sq.
				Select("jsonb_agg(speciality)").
				From(new(schema.UserSpeciality).TableName()).
				Where(sq.Eq{"user_uuid": id})

			raw, err := database.GetTxX[[]byte](ctx, tx, q)
			if err != nil {
				return err
			}

			if len(raw) == 0 {
				m.Specialities = []string{}
				return nil
			}

			err = json.Unmarshal(raw, &m.Specialities)
			return err
		},
	)

	return
}

func (repo) getEmail(ctx context.Context, id uuid.UUID) (string, error) {
	q := sq.
		Select("email").
		From(new(schema.UserAuth).TableName()).
		Where(sq.Eq{
			"id": id,
		})
	return database.GetX[string](ctx, q)
}

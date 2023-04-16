package search

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/artemmarkaryan/exlex-backend/internal/schema"
	"github.com/artemmarkaryan/exlex-backend/pkg/database"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type repo struct{}

func (repo) create(ctx context.Context, r CreateSearchRequest) (id uuid.UUID, err error) {
	err = database.Tx(
		ctx, database.DefaultTxOptions(),

		func(tx *sqlx.Tx) error {
			m := map[string]any{
				"creator":                  r.Creator,
				"name":                     r.Name,
				"description":              r.Description,
				"price":                    r.Price,
				"required_work_experience": r.RequiredWorkExperience,
				"deadline":                 r.Deadline,
			}

			q := sq.
				Insert(new(schema.Search).TableName()).
				SetMap(m).
				Suffix("RETURNING id")

			var errTx error
			id, errTx = database.GetFromInsertTxX[uuid.UUID](ctx, tx, q)
			return errTx
		},

		func(tx *sqlx.Tx) error {
			if len(r.RequiredSpecialities) == 0 {
				return nil
			}

			q := sq.
				Insert(new(schema.SearchRequirementSpeciality).TableName()).
				Columns("search_uuid", "speciality")

			for _, s := range r.RequiredSpecialities {
				q = q.Values(id, s)
			}

			_, errTx := database.InsertTxX(ctx, tx, q)
			return errTx
		},

		func(tx *sqlx.Tx) error {
			if len(r.RequiredEducation) == 0 {
				return nil
			}

			q := sq.
				Insert(new(schema.SearchRequirementEducation).TableName()).
				Columns("search_uuid", "education")

			for _, e := range r.RequiredEducation {
				q = q.Values(id, e)
			}

			_, errTx := database.InsertTxX(ctx, tx, q)
			return errTx
		},
	)

	return
}

func (repo) delete(ctx context.Context, id uuid.UUID) error {
	err := database.Tx(
		ctx, database.DefaultTxOptions(),

		func(tx *sqlx.Tx) error {
			q := sq.
				Delete(new(schema.SearchApplication).TableName()).
				Where(sq.Eq{"search_id": id})

			_, errTx := database.DeleteTxX(ctx, tx, q)
			return errTx
		},

		func(tx *sqlx.Tx) error {
			q := sq.
				Delete(new(schema.SearchRequirementSpeciality).TableName()).
				Where(sq.Eq{"search_uuid": id})

			_, errTx := database.DeleteTxX(ctx, tx, q)
			return errTx
		},

		func(tx *sqlx.Tx) error {
			q := sq.
				Delete(new(schema.SearchRequirementEducation).TableName()).
				Where(sq.Eq{"search_uuid": id})

			_, errTx := database.DeleteTxX(ctx, tx, q)
			return errTx
		},

		func(tx *sqlx.Tx) error {
			q := sq.
				Delete(new(schema.Search).TableName()).
				Where(sq.Eq{"id": id})

			_, errTx := database.DeleteTxX(ctx, tx, q)
			return errTx
		},
	)

	return err
}

func (repo) checkCreator(ctx context.Context, user uuid.UUID, search uuid.UUID) error {
	q := sq.
		Select("creator").
		From(new(schema.Search).TableName()).
		Where(sq.Eq{"id": search})

	creator, err := database.GetX[uuid.UUID](ctx, q)
	if err != nil {
		return err
	}

	if creator != user {
		return ErrUnauthorized
	}

	return nil
}

func (r repo) get(ctx context.Context, search uuid.UUID) (d schema.SearchFullDataRaw, err error) {
	q := sq.
		Select(
			"s.id",
			"s.name",
			"s.description",
			"s.price",
			"s.required_work_experience",
			"s.created_at",
			"s.deadline",
			"jsonb_agg(e.education) as education",
			"jsonb_agg(sp.speciality) as speciality",
		).
		From(new(schema.Search).TableName() + " s").
		LeftJoin(new(schema.SearchRequirementEducation).TableName() + " e on s.id = e.search_uuid").
		LeftJoin(new(schema.SearchRequirementSpeciality).TableName() + " sp on s.id = sp.search_uuid").
		LeftJoin(new(schema.SearchApplication).TableName() + " sa on s.id = sa.search_id").
		Where(sq.Eq{"s.id": search}).
		OrderBy("s.created_at desc").
		GroupBy("s.id")

	return database.GetX[schema.SearchFullDataRaw](ctx, q)
}

func (r repo) listAvailableForApplication(ctx context.Context, user uuid.UUID) (d []schema.SearchFullDataRaw, err error) {
	q := sq.
		Select(
			"s.id",
			"s.name",
			"s.description",
			"s.price",
			"s.required_work_experience",
			"s.created_at",
			"s.deadline",
			"jsonb_agg(e.education) as education",
			"jsonb_agg(sp.speciality) as speciality",
		).
		From(new(schema.Search).TableName() + " s").
		LeftJoin(new(schema.SearchRequirementEducation).TableName() + " e on s.id = e.search_uuid").
		LeftJoin(new(schema.SearchRequirementSpeciality).TableName() + " sp on s.id = sp.search_uuid").
		LeftJoin(new(schema.SearchApplication).TableName() + " sa on s.id = sa.search_id").
		Where(sq.Or{
			sq.Eq{"sa.id": nil},          // doesn't have application
			sq.NotEq{"sa.user_id": user}, // is applied by another user
		}).
		OrderBy("s.created_at desc").
		GroupBy("s.id")

	return database.SelectX[schema.SearchFullDataRaw](ctx, q)
}

func (r repo) listByAuthor(ctx context.Context, user uuid.UUID) ([]schema.SearchFullDataRaw, error) {
	q := sq.
		Select(
			"s.id",
			"s.name",
			"s.description",
			"s.price",
			"s.required_work_experience",
			"s.created_at",
			"s.deadline",
			"jsonb_agg(e.education) as education",
			"jsonb_agg(sp.speciality) as speciality",
		).
		From(new(schema.Search).TableName() + " s").
		LeftJoin(new(schema.SearchRequirementEducation).TableName() + " e on s.id = e.search_uuid").
		LeftJoin(new(schema.SearchRequirementSpeciality).TableName() + " sp on s.id = sp.search_uuid").
		OrderBy("s.created_at desc").
		Where(sq.Eq{"creator": user}).
		GroupBy("s.id")

	return database.SelectX[schema.SearchFullDataRaw](ctx, q)
}

func (repo) apply(ctx context.Context, r SearchApplicationRequest) (applicationID uuid.UUID, err error) {
	values := map[string]interface{}{
		"search_id": r.SearchID,
		"user_id":   r.UserID,
	}

	if r.Comment != nil {
		values["comment"] = *r.Comment
	}

	q := sq.
		Insert(new(schema.SearchApplication).TableName()).
		SetMap(values).
		Suffix("returning id")

	applicationID, err = database.GetFromInsertX[uuid.UUID](ctx, q)
	return
}

func (r repo) listApplications(ctx context.Context, searchID uuid.UUID) (a []schema.SearchApplicationRaw, err error) {
	q := `
select 
	a.id as application_id,
	a.user_id,
	a.created_at,
	a.comment,
	em.education,
	em.full_name,
	em.experience_years,
	jsonb_agg(us.speciality) as speciality
from 
	search_application a 
		inner	join user_speciality us on   	a.user_id = us.user_uuid
		left  	join executor_metadata em on 	a.user_id = em.user_uuid
where 
	search_id = ?
group by 
	application_id,
	a.user_id, 
	a.created_at,
	a.comment,
	em.education,
	em.full_name,
	em.experience_years
order by created_at desc
;`

	return database.SelectRawX[schema.SearchApplicationRaw](ctx, q, searchID)
}

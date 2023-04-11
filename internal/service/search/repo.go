package search

import (
	"context"
	"encoding/json"

	sq "github.com/Masterminds/squirrel"
	"github.com/artemmarkaryan/exlex-backend/internal/schema"
	"github.com/artemmarkaryan/exlex-backend/pkg/database"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
)

type repo struct{}

func (repo) create(ctx context.Context, d CreateSearch) (id uuid.UUID, err error) {
	err = database.Tx(
		ctx, database.DefaultTxOptions(),

		func(tx *sqlx.Tx) error {
			m := map[string]any{
				"creator":                  d.Creator,
				"name":                     d.Name,
				"description":              d.Description,
				"price":                    d.Price,
				"required_work_experience": d.RequiredWorkExperience,
				"deadline":                 d.Deadline,
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
			if len(d.RequiredSpecialities) == 0 {
				return nil
			}

			q := sq.
				Insert(new(schema.SearchRequirementSpeciality).TableName()).
				Columns("search_uuid", "speciality")

			for _, s := range d.RequiredSpecialities {
				q = q.Values(id, s)
			}

			_, errTx := database.InsertTxX(ctx, tx, q)
			return errTx
		},

		func(tx *sqlx.Tx) error {
			if len(d.RequiredEducation) == 0 {
				return nil
			}

			q := sq.
				Insert(new(schema.SearchRequirementEducation).TableName()).
				Columns("search_uuid", "education")

			for _, e := range d.RequiredEducation {
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

func (r repo) get(ctx context.Context, search uuid.UUID) (d schema.SearchFullData, err error) {
	// err = database.Tx(
	// 	ctx, database.DefaultTxOptions(),
	// 	func(tx *sqlx.Tx) (errTx error) {
	// 		q := sq.
	// 			Select("*").
	// 			From(new(schema.Search).TableName()).
	// 			Where(sq.Eq{"id": search})

	// 		d.Search, errTx = database.GetTxX[schema.Search](ctx, tx, q)
	// 		return
	// 	},
	// 	func(tx *sqlx.Tx) (errTx error) {
	// 		q := sq.
	// 			Select("*").
	// 			From(new(schema.SearchRequirementEducation).TableName()).
	// 			Where(sq.Eq{"search_uuid": search})

	// 		d.Education, errTx = database.SelectTxX[schema.SearchRequirementEducation](ctx, tx, q)
	// 		return
	// 	},
	// 	func(tx *sqlx.Tx) (errTx error) {
	// 		q := sq.
	// 			Select("*").
	// 			From(new(schema.SearchRequirementSpeciality).TableName()).
	// 			Where(sq.Eq{"search_uuid": search})

	// 		d.Speciality, errTx = database.SelectTxX[schema.SearchRequirementSpeciality](ctx, tx, q)
	// 		return
	// 	},
	// )

	return
}

func (r repo) list(ctx context.Context, user uuid.UUID) (d []schema.SearchFullData, err error) {
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
		OrderBy("s.created_at").
		Where(sq.Eq{"creator": user}).
		GroupBy("s.id")

	dbos, err := database.SelectX[schema.SearchFullDataRaw](ctx, q)

	for _, dbo := range dbos {
		var e []string
		err = json.Unmarshal(dbo.Education, &e)
		if err != nil {
			return
		}

		e = lo.Filter(e, func(obj string, _ int) bool { return obj != "" })

		var s []string
		err = json.Unmarshal(dbo.Speciality, &s)
		if err != nil {
			return
		}

		s = lo.Filter(s, func(obj string, _ int) bool { return obj != "" })

		search := schema.SearchFullData{
			ID:                     dbo.ID,
			Name:                   dbo.Name,
			Description:            dbo.Description,
			Price:                  dbo.Price,
			RequiredWorkExperience: dbo.RequiredWorkExperience,
			Deadline:               dbo.Deadline,
			CreatedAt:              dbo.CreatedAt,
			Education:              e,
			Speciality:             s,
		}

		d = append(d, search)
	}

	return
}

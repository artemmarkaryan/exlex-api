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
	"golang.org/x/exp/slices"
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

func (r repo) listAvailableForApplication(ctx context.Context, user uuid.UUID) (d []schema.SearchFullData, err error) {
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

	dbos, err := database.SelectX[schema.SearchFullDataRaw](ctx, q)

	filter := func(s []string) []string {
		slices.Sort(s)
		s = slices.Compact(s)
		s = lo.Filter(s, func(obj string, _ int) bool { return obj != "" })
		return s
	}

	for _, dbo := range dbos {
		var e []string
		err = json.Unmarshal(dbo.Education, &e)
		if err != nil {
			return
		}

		var s []string
		err = json.Unmarshal(dbo.Speciality, &s)
		if err != nil {
			return
		}

		search := schema.SearchFullData{
			ID:                     dbo.ID,
			Name:                   dbo.Name,
			Description:            dbo.Description,
			Price:                  dbo.Price,
			RequiredWorkExperience: dbo.RequiredWorkExperience,
			Deadline:               dbo.Deadline,
			CreatedAt:              dbo.CreatedAt,
			Education:              filter(e),
			Speciality:             filter(s),
		}

		d = append(d, search)
	}

	return
}

func (r repo) listByAuthor(ctx context.Context, user uuid.UUID) (d []schema.SearchFullData, err error) {
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

	dbos, err := database.SelectX[schema.SearchFullDataRaw](ctx, q)

	filter := func(s []string) []string {
		slices.Sort(s)
		s = slices.Compact(s)
		s = lo.Filter(s, func(obj string, _ int) bool { return obj != "" })
		return s
	}

	for _, dbo := range dbos {
		var e []string
		err = json.Unmarshal(dbo.Education, &e)
		if err != nil {
			return
		}

		var s []string
		err = json.Unmarshal(dbo.Speciality, &s)
		if err != nil {
			return
		}

		search := schema.SearchFullData{
			ID:                     dbo.ID,
			Name:                   dbo.Name,
			Description:            dbo.Description,
			Price:                  dbo.Price,
			RequiredWorkExperience: dbo.RequiredWorkExperience,
			Deadline:               dbo.Deadline,
			CreatedAt:              dbo.CreatedAt,
			Education:              filter(e),
			Speciality:             filter(s),
		}

		d = append(d, search)
	}

	return
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

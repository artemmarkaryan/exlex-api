package authentication

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/artemmarkaryan/exlex-backend/internal/schema"
	"github.com/artemmarkaryan/exlex-backend/pkg/database"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type repo struct{}

func (r repo) getUser(ctx context.Context, email string) (id uuid.UUID, err error) {
	return database.GetX[uuid.UUID](
		ctx, sq.
			Select("id").
			From(new(schema.UserAuth).TableName()).
			Where(sq.Eq{"email": email}),
	)
}

func (r repo) createUser(ctx context.Context, email string, role schema.Role) (id uuid.UUID, err error) {
	err = database.Tx(
		ctx, database.DefaultTxOptions(),
		func(tx *sqlx.Tx) (errTx error) {
			id, errTx = database.GetX[uuid.UUID](
				ctx, sq.
					Select("id").
					From(new(schema.UserAuth).TableName()).
					Where(sq.Eq{"email": email}),
			)
			if errTx == nil {
				errTx = ErrAlreadyExists
				return
			}

			if errTx == sql.ErrNoRows {
				errTx = nil
				return
			}

			return
		},
		func(tx *sqlx.Tx) (errTx error) {
			q := sq.
				Insert(new(schema.UserAuth).TableName()).
				SetMap(map[string]any{"email": email}).
				Suffix("returning id")

			id, errTx = database.GetFromInsertTxX[uuid.UUID](ctx, tx, q)
			return
		},
		func(tx *sqlx.Tx) (errTx error) {
			q := sq.
				Insert(new(schema.UserRole).TableName()).
				SetMap(
					map[string]any{
						"user_id": id,
						"role":    role,
					},
				)

			_, errTx = database.InsertTxX(ctx, tx, q)
			return
		},
	)

	return
}

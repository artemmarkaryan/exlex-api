package otp

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/artemmarkaryan/exlex-backend/pkg/database"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type repo struct{}

func (repo) insert(ctx context.Context, id uuid.UUID, otp string) error {
	return database.Tx(
		ctx, database.DefaultTxOptions(),

		func(tx *sqlx.Tx) error {
			q := sq.
				Delete(new(UserOTP).TableName()).
				Where(sq.Eq{"user_uuid": id})

			_, err := database.DeleteTxX(ctx, tx, q)
			return err
		},

		func(tx *sqlx.Tx) error {
			q := sq.
				Insert(new(UserOTP).TableName()).
				Columns("user_uuid", "otp").
				Values(id, otp)

			_, err := database.InsertTxX(ctx, tx, q)
			return err
		},
	)
}

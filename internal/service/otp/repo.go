package otp

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/artemmarkaryan/exlex-backend/internal/schema"
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
				Delete(new(schema.UserOTP).TableName()).
				Where(sq.Eq{"user_uuid": id})

			_, err := database.DeleteTxX(ctx, tx, q)
			return err
		},

		func(tx *sqlx.Tx) error {
			q := sq.
				Insert(new(schema.UserOTP).TableName()).
				Columns("user_uuid", "otp").
				Values(id, otp)

			_, err := database.InsertTxX(ctx, tx, q)
			return err
		},
	)
}

func (r repo) get(ctx context.Context, email string) (string, error) {
	u := new(schema.UserAuth).TableName()
	q := sq.
		Select("otp.otp").
		From(new(schema.UserOTP).TableName() + " otp").
		InnerJoin(u + " u on otp.user_uuid = u.id").
		Where(sq.Eq{"u.email": email})

	return database.GetX[string](ctx, q)
}

func (r repo) delete(ctx context.Context, email string) error {
	u := new(schema.UserAuth).TableName()
	e, args, err := sq.Expr("user_uuid = (select id from "+u+" where email = ?)", email).ToSql()
	if err != nil {
		return err
	}

	q := sq.
		Delete(new(schema.UserOTP).TableName()).
		Where(e, args...)

	_, err = database.DeleteX(ctx, q)
	return err
}

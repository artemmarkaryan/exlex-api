package database

import (
	"context"
	"database/sql"
	"errors"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

const key = "database-connection"

func Connect(ctx context.Context, dsn string) *sqlx.DB {
	if db, ok := ctx.Value(key).(*sqlx.DB); ok {
		return db
	}

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				_ = db.Close()
			}
		}
	}()

	return db
}

// C is short for "Connect", cause it use it often
func C(ctx context.Context) *sqlx.DB {
	return ctx.Value(key).(*sqlx.DB)
}

func Propagate(ctx context.Context, db *sqlx.DB) context.Context {
	return context.WithValue(ctx, key, db)
}

type Model interface{}

var NotFound = errors.New("not found")

func Getx[m Model](ctx context.Context, builder sq.SelectBuilder) (result m, err error) {
	query, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return
	}

	dst := new(m)
	err = C(ctx).GetContext(ctx, dst, query, args...)
	return *dst, err
}

func Execx(ctx context.Context, builder sq.InsertBuilder) (sql.Result, error) {
	query, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	return C(ctx).ExecContext(ctx, query, args...)
}

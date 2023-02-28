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

type contextGetter interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type contextSelecter interface {
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type contextExecer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

func selectX[m Model](ctx context.Context, c contextSelecter, builder sq.SelectBuilder) (result []m, err error) {
	query, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return
	}

	err = c.SelectContext(ctx, &result, query, args...)
	return
}

func SelectX[m Model](ctx context.Context, b sq.SelectBuilder) (result []m, err error) {
	return selectX[m](ctx, C(ctx), b)
}

func SelectTxX[m Model](ctx context.Context, tx *sqlx.Tx, b sq.SelectBuilder) (result []m, err error) {
	return selectX[m](ctx, tx, b)
}

func getX[m Model](ctx context.Context, c contextGetter, builder sq.SelectBuilder) (result m, err error) {
	query, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return
	}

	err = c.GetContext(ctx, &result, query, args...)
	return
}

func GetX[m Model](ctx context.Context, builder sq.SelectBuilder) (result m, err error) {
	return getX[m](ctx, C(ctx), builder)
}

func GetTxX[m Model](ctx context.Context, tx *sqlx.Tx, builder sq.SelectBuilder) (result m, err error) {
	return getX[m](ctx, tx, builder)
}

func getFromInsertX[m Model](ctx context.Context, c contextGetter, builder sq.InsertBuilder) (result m, err error) {
	query, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return
	}

	dst := new(m)
	err = c.GetContext(ctx, dst, query, args...)
	return *dst, err
}

func GetFromInsertX[m Model](ctx context.Context, builder sq.InsertBuilder) (result m, err error) {
	return getFromInsertX[m](ctx, C(ctx), builder)
}

func GetFromInsertTxX[m Model](ctx context.Context, tx *sqlx.Tx, builder sq.InsertBuilder) (result m, err error) {
	return getFromInsertX[m](ctx, tx, builder)
}

func insertX(ctx context.Context, c contextExecer, builder sq.InsertBuilder) (sql.Result, error) {
	query, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	return c.ExecContext(ctx, query, args...)
}

func InsertX(ctx context.Context, builder sq.InsertBuilder) (sql.Result, error) {
	return insertX(ctx, C(ctx), builder)
}

func InsertTxX(ctx context.Context, tx *sqlx.Tx, builder sq.InsertBuilder) (sql.Result, error) {
	return insertX(ctx, tx, builder)
}

func deleteX(ctx context.Context, c contextExecer, builder sq.DeleteBuilder) (sql.Result, error) {
	query, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	return c.ExecContext(ctx, query, args...)
}

func DeleteX(ctx context.Context, builder sq.DeleteBuilder) (sql.Result, error) {
	return deleteX(ctx, C(ctx), builder)
}

func DeleteTxX(ctx context.Context, tx *sqlx.Tx, builder sq.DeleteBuilder) (sql.Result, error) {
	return deleteX(ctx, tx, builder)
}

func DefaultTxOptions() *sql.TxOptions {
	return &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	}
}

func Tx(ctx context.Context, opts *sql.TxOptions, f ...func(tx *sqlx.Tx) error) error {
	tx, err := C(ctx).BeginTxx(ctx, opts)
	if err != nil {
		return err
	}

	defer func() { _ = tx.Rollback() }()

	for _, ff := range f {
		err = ff(tx)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

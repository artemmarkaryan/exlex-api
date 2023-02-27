package authentication

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/artemmarkaryan/exlex-backend/pkg/database"
	"github.com/google/uuid"
)

type repo struct{}

func (r repo) getOrCreateUser(ctx context.Context, email string) (id uuid.UUID, err error) {
	id, err = database.GetX[uuid.UUID](
		ctx, sq.
			Select("id").
			From(new(User).TableName()).
			Where(sq.Eq{"email": email}),
	)

	if err != sql.ErrNoRows {
		return
	}

	q := sq.
		Insert(new(User).TableName()).
		SetMap(map[string]any{"email": email}).
		Suffix("returning id")

	id, err = database.GetFromInsertX[uuid.UUID](ctx, q)
	return
}

package authentication

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/artemmarkaryan/exlex-backend/pkg/database"
)

type Service struct{}

func Make(_ context.Context) (s Service) { return }

func (Service) CreateUser(ctx context.Context, login string, password string) (int64, error) {
	user, err := database.Getx[User](ctx, sq.
		Select("*").
		From(new(User).TableName()).
		Where(sq.Eq{"email": login}),
	)
	if err != nil {
		if err != database.NotFound {
			return 0, err
		} else {
			log.Println("not found")
		}
	}

	log.Println(user)

	return 0, nil
	//hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	//if err != nil {
	//	return 0, fmt.Errorf("generating password hash: %w", err)
	//}
}

package authentication

type User struct {
	ID           string `db:"id"`
	Email        string `db:"email"`
	PasswordHash string `db:"psw_hash"`
}

func (u User) TableName() string { return "user_auth" }

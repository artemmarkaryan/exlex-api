package authentication

type User struct {
	ID    string `db:"id"`
	Email string `db:"email"`
}

func (u User) TableName() string { return "user_auth" }

type Token struct {
	Refresh string
	Access  string
}

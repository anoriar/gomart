package entity

type User struct {
	ID       string `db:"id"`
	Login    string `db:"login"`
	Password string `db:"password"`
	Salt     string `db:"salt"`
}

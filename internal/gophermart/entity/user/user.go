package user

type User struct {
	Id       string `db:"id"`
	Login    string `db:"login"`
	Password string `db:"password"`
	Salt     string `db:"salt"`
}

package register

type RegisterUserRequestDto struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

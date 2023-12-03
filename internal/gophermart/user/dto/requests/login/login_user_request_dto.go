package login

type LoginUserRequestDto struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

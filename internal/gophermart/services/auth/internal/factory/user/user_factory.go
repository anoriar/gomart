package user

import (
	"github.com/anoriar/gophermart/internal/gophermart/entity/user"
	"github.com/google/uuid"
)

type UserFactory struct {
}

func NewUserFactory() *UserFactory {
	return &UserFactory{}
}

func (factory *UserFactory) Create(login string, password string, salt string) user.User {
	return user.User{
		ID:       uuid.NewString(),
		Login:    login,
		Password: password,
		Salt:     salt,
	}
}

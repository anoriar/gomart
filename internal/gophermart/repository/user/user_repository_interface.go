package user

import (
	"context"
	"github.com/anoriar/gophermart/internal/gophermart/entity/user"
)

type UserRepositoryInterface interface {
	AddUser(ctx context.Context, user user.User) error
}

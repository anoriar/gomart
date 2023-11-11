package user

import (
	"context"
	"github.com/anoriar/gophermart/internal/gophermart/entity/user"
)

//go:generate mockgen -source=user_repository_interface.go -destination=mock_user_repository/user_repository.go -package=mock_user_repository
type UserRepositoryInterface interface {
	AddUser(ctx context.Context, user user.User) error
	GetUserByLogin(ctx context.Context, login string) (user.User, error)
}

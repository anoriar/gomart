package user

import (
	"github.com/anoriar/gophermart/internal/gophermart/user/entity"
)

//go:generate mockgen -source=user_factory_interface.go -destination=mock_user_factory/user_factory.go -package=mock_user_factory
type UserFactoryInterface interface {
	Create(login string, password string, salt string) entity.User
}

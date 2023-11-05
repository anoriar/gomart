package auth

import (
	"context"
	"github.com/anoriar/gophermart/internal/gophermart/dto/requests/register"
)

type RegisterServiceInterface interface {
	RegisterUser(ctx context.Context, dto register.RegisterUserRequestDto) error
}

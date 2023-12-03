package token

import (
	"github.com/anoriar/gophermart/internal/gophermart/user/dto/auth"
)

//go:generate mockgen -source=token_service_interface.go -destination=mock_token_service/token_service.go -package=mock_token_service
type TokenSerivceInterface interface {
	BuildTokenString(userClaims auth.UserClaims) (string, error)
	GetUserClaims(tokenString string) (auth.UserClaims, error)
}

package token

import "github.com/anoriar/gophermart/internal/gophermart/dto/auth"

//go:generate mockgen -source=token_service_interface.go -destination=mock/token_service.go -package=mock
type TokenSerivceInterface interface {
	BuildTokenString(userClaims auth.UserClaims) (string, error)
	GetUserClaims(tokenString string) (auth.UserClaims, error)
}

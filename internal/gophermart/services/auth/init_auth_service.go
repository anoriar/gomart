package auth

import (
	"github.com/anoriar/gophermart/internal/gophermart/config"
	"github.com/anoriar/gophermart/internal/gophermart/repository/user"
	"github.com/anoriar/gophermart/internal/gophermart/services/auth/internal/factory/salt"
	user2 "github.com/anoriar/gophermart/internal/gophermart/services/auth/internal/factory/user"
	"github.com/anoriar/gophermart/internal/gophermart/services/auth/internal/services/password"
	"github.com/anoriar/gophermart/internal/gophermart/services/auth/internal/services/token/jwt"
	"go.uber.org/zap"
)

func InitializeAuthService(config *config.Config, userRepository user.UserRepositoryInterface, logger *zap.Logger) *AuthService {
	return NewAuthService(
		userRepository,
		password.NewArgonPasswordService(),
		jwt.NewJWTTokenService(config.JwtSecretKey),
		user2.NewUserFactory(),
		salt.NewSaltFactory(),
		logger,
	)
}

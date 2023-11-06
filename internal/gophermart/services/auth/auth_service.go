package auth

import (
	"context"
	"encoding/hex"
	"github.com/anoriar/gophermart/internal/gophermart/dto/auth"
	"github.com/anoriar/gophermart/internal/gophermart/dto/requests/register"
	"github.com/anoriar/gophermart/internal/gophermart/repository/user"
	"github.com/anoriar/gophermart/internal/gophermart/services/auth/internal/factory/salt"
	user2 "github.com/anoriar/gophermart/internal/gophermart/services/auth/internal/factory/user"
	"github.com/anoriar/gophermart/internal/gophermart/services/auth/internal/services/password"
	"github.com/anoriar/gophermart/internal/gophermart/services/auth/internal/services/token"
	"go.uber.org/zap"
)

type AuthService struct {
	userRepository  user.UserRepositoryInterface
	passwordService password.PasswordServiceInterface
	tokenService    token.TokenSerivceInterface
	userFactory     user2.UserFactoryInterface
	saltFactory     salt.SaltFactoryInterface
	logger          *zap.Logger
}

func NewAuthService(
	userRepository user.UserRepositoryInterface,
	passwordService password.PasswordServiceInterface,
	tokenService token.TokenSerivceInterface,
	userFactory user2.UserFactoryInterface,
	saltFactory salt.SaltFactoryInterface,
	logger *zap.Logger,
) *AuthService {
	return &AuthService{
		userRepository:  userRepository,
		passwordService: passwordService,
		tokenService:    tokenService,
		userFactory:     userFactory,
		saltFactory:     saltFactory,
		logger:          logger,
	}
}

func (service *AuthService) RegisterUser(ctx context.Context, registerUserDto register.RegisterUserRequestDto) (string, error) {
	salt, err := service.saltFactory.GenerateSalt()
	if err != nil {
		service.logger.Error(err.Error())
		return "", err
	}
	hashedPassword := service.passwordService.GenerateHashedPassword(registerUserDto.Password, salt)

	newUser := service.userFactory.Create(registerUserDto.Login, hashedPassword, hex.EncodeToString(salt))
	err = service.userRepository.AddUser(ctx, newUser)
	if err != nil {
		service.logger.Error(err.Error())
		return "", err
	}

	tokenString, err := service.tokenService.BuildTokenString(auth.UserClaims{UserID: newUser.Id})
	if err != nil {
		service.logger.Error(err.Error())
		return "", err
	}
	return tokenString, nil
}

package auth

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/anoriar/gophermart/internal/gophermart/domain_errors"
	"github.com/anoriar/gophermart/internal/gophermart/dto/auth"
	"github.com/anoriar/gophermart/internal/gophermart/dto/requests/login"
	"github.com/anoriar/gophermart/internal/gophermart/dto/requests/register"
	"github.com/anoriar/gophermart/internal/gophermart/repository/user"
	"github.com/anoriar/gophermart/internal/gophermart/services/auth/internal/factory/salt"
	user2 "github.com/anoriar/gophermart/internal/gophermart/services/auth/internal/factory/user"
	"github.com/anoriar/gophermart/internal/gophermart/services/auth/internal/services/password"
	"github.com/anoriar/gophermart/internal/gophermart/services/auth/internal/services/token"
	"github.com/anoriar/gophermart/internal/gophermart/services/auth/internal/services/token/token_errors"
	"go.uber.org/zap"
)

var ErrUnauthorized = errors.New("user is unauthorized")
var ErrUserAlreadyExists = errors.New("user already exists")

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
		if errors.Is(err, domain_errors.ErrConflict) {
			return "", ErrUserAlreadyExists
		}
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

func (service *AuthService) LoginUser(ctx context.Context, dto login.LoginUserRequestDto) (string, error) {
	existedUser, err := service.userRepository.GetUserByLogin(ctx, dto.Login)
	if err != nil {
		if errors.Is(err, domain_errors.ErrNotFound) {
			return "", ErrUnauthorized
		}
		service.logger.Error(err.Error())
		return "", err
	}
	saltInBytes, err := hex.DecodeString(existedUser.Salt)
	if err != nil {
		service.logger.Error(err.Error())
		return "", err
	}
	hashedPasswordFromRequest := service.passwordService.GenerateHashedPassword(dto.Password, saltInBytes)

	if hashedPasswordFromRequest != existedUser.Password {
		return "", ErrUnauthorized
	}

	tokenString, err := service.tokenService.BuildTokenString(auth.UserClaims{UserID: existedUser.Id})
	if err != nil {
		service.logger.Error(err.Error())
		return "", err
	}
	return tokenString, nil
}

func (service *AuthService) ValidateToken(token string) (auth.UserClaims, error) {
	claims, err := service.tokenService.GetUserClaims(token)
	if err != nil {
		if errors.Is(err, token_errors.ErrTokenNotValid) {
			return auth.UserClaims{}, ErrUnauthorized
		}
		service.logger.Error(err.Error())
		return auth.UserClaims{}, err
	}

	return claims, nil
}
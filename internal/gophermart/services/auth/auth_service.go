package auth

import (
	"context"
	"github.com/anoriar/gophermart/internal/gophermart/dto/requests/register"
	"github.com/anoriar/gophermart/internal/gophermart/repository/user"
	"github.com/anoriar/gophermart/internal/gophermart/services/auth/internal/factory"
	"github.com/anoriar/gophermart/internal/gophermart/services/auth/internal/services"
)

type AuthService struct {
	userRepository  user.UserRepositoryInterface
	passwordService services.PasswordServiceInterface
	userFactory     *factory.UserFactory
	saltFactory     *factory.SaltFactory
}

func NewAuthService(userRepository user.UserRepositoryInterface, passwordService services.PasswordServiceInterface, userFactory *factory.UserFactory) *AuthService {
	return &AuthService{userRepository: userRepository, passwordService: passwordService, userFactory: userFactory}
}

func (service *AuthService) RegisterUser(ctx context.Context, registerUserDto register.RegisterUserRequestDto) error {
	salt, err := service.saltFactory.GenerateSalt()
	if err != nil {
		return err
	}
	hashedPassword := service.passwordService.GenerateHashedPassword(registerUserDto.Password, salt)

	newUser := service.userFactory.Create(registerUserDto.Login, hashedPassword, string(salt))
	err = service.userRepository.AddUser(ctx, newUser)
	if err != nil {
		return err
	}
	//TODO: return jwt token
	return nil
}

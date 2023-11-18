package auth

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/anoriar/gophermart/internal/gophermart/app/logger"
	"github.com/anoriar/gophermart/internal/gophermart/domainerrors"
	"github.com/anoriar/gophermart/internal/gophermart/dto/auth"
	"github.com/anoriar/gophermart/internal/gophermart/dto/requests/register"
	"github.com/anoriar/gophermart/internal/gophermart/entity/user"
	"github.com/anoriar/gophermart/internal/gophermart/repository/user/mock_user_repository"
	"github.com/anoriar/gophermart/internal/gophermart/services/auth/internal/factory/salt/mock_salt_factory"
	"github.com/anoriar/gophermart/internal/gophermart/services/auth/internal/factory/user/mock_user_factory"
	"github.com/anoriar/gophermart/internal/gophermart/services/auth/internal/services/password/mock_password_service"
	"github.com/anoriar/gophermart/internal/gophermart/services/auth/internal/services/token/mock_token_service"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAuthService_RegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepositoryMock := mock_user_repository.NewMockUserRepositoryInterface(ctrl)
	saltFactoryMock := mock_salt_factory.NewMockSaltFactoryInterface(ctrl)
	tokenServiceMock := mock_token_service.NewMockTokenSerivceInterface(ctrl)
	passwordServiceMock := mock_password_service.NewMockPasswordServiceInterface(ctrl)
	userFactoryMock := mock_user_factory.NewMockUserFactoryInterface(ctrl)
	loggerMock, err := logger.Initialize("info")
	require.NoError(t, err)

	type args struct {
		ctx             context.Context
		registerUserDto register.RegisterUserRequestDto
	}
	tests := []struct {
		name          string
		mockBehaviour func()
		args          args
		want          string
		wantErr       bool
	}{
		{
			name: "success",
			mockBehaviour: func() {
				salt := []byte{2, 3, 24, 32, 11}
				hashedPassword := "dashdiasdss"
				encodedSalt := hex.EncodeToString(salt)
				userID := "ab161651-ec2c-4cbb-a6c6-f8cf414e503d"
				login := "test"
				userPassword := "user-password"
				userMock := user.User{
					ID:       userID,
					Login:    login,
					Password: hashedPassword,
					Salt:     encodedSalt,
				}
				saltFactoryMock.EXPECT().GenerateSalt().Return(salt, nil).Times(1)
				passwordServiceMock.EXPECT().GenerateHashedPassword(userPassword, salt).Return(hashedPassword).Times(1)
				userFactoryMock.EXPECT().Create(login, hashedPassword, encodedSalt).Return(userMock).Times(1)
				userRepositoryMock.EXPECT().AddUser(gomock.Any(), userMock).Return(nil).Times(1)
				tokenServiceMock.EXPECT().BuildTokenString(auth.UserClaims{UserID: userID}).Return("result-token", nil).Times(1)
			},
			args: args{
				ctx: context.Background(),
				registerUserDto: register.RegisterUserRequestDto{
					Login:    "test",
					Password: "user-password",
				},
			},
			want:    "result-token",
			wantErr: false,
		},
		{
			name: "conflict",
			mockBehaviour: func() {
				salt := []byte{2, 3, 24, 32, 11}
				hashedPassword := "dashdiasdss"
				encodedSalt := hex.EncodeToString(salt)
				userID := "ab161651-ec2c-4cbb-a6c6-f8cf414e503d"
				login := "test"
				userPassword := "user-password"
				userMock := user.User{
					ID:       userID,
					Login:    login,
					Password: hashedPassword,
					Salt:     encodedSalt,
				}
				saltFactoryMock.EXPECT().GenerateSalt().Return(salt, nil).Times(1)
				passwordServiceMock.EXPECT().GenerateHashedPassword(userPassword, salt).Return(hashedPassword).Times(1)
				userFactoryMock.EXPECT().Create(login, hashedPassword, encodedSalt).Return(userMock).Times(1)
				userRepositoryMock.EXPECT().AddUser(gomock.Any(), userMock).Return(domainerrors.ErrConflict).Times(1)
				tokenServiceMock.EXPECT().BuildTokenString(auth.UserClaims{UserID: userID}).Times(0)
			},
			args: args{
				ctx: context.Background(),
				registerUserDto: register.RegisterUserRequestDto{
					Login:    "test",
					Password: "user-password",
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "internal server error",
			mockBehaviour: func() {
				saltFactoryMock.EXPECT().GenerateSalt().Return([]byte{}, errors.New("internal server error")).Times(1)
			},
			args: args{
				ctx: context.Background(),
				registerUserDto: register.RegisterUserRequestDto{
					Login:    "test",
					Password: "user-password",
				},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour()

			service := &AuthService{
				userRepository:  userRepositoryMock,
				passwordService: passwordServiceMock,
				tokenService:    tokenServiceMock,
				userFactory:     userFactoryMock,
				saltFactory:     saltFactoryMock,
				logger:          loggerMock,
			}
			got, err := service.RegisterUser(tt.args.ctx, tt.args.registerUserDto)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RegisterUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

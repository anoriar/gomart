package register

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/anoriar/gophermart/internal/gophermart/user/dto/requests/register"
	"github.com/anoriar/gophermart/internal/gophermart/user/services/auth"
	"github.com/anoriar/gophermart/internal/gophermart/user/services/auth/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterHandler_Register(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	successAuthHeader := "f39tj804jgt5j9gr"

	successRequestDto := register.RegisterUserRequestDto{
		Login:    "test",
		Password: "fihh8hgt9g",
	}
	successRequestBody, err := json.Marshal(successRequestDto)
	require.NoError(t, err)

	badRequestDto := register.RegisterUserRequestDto{}
	badRequestBody, err := json.Marshal(badRequestDto)
	require.NoError(t, err)

	authServiceMock := mock.NewMockAuthServiceInterface(ctrl)

	type want struct {
		status     int
		authHeader string
	}
	tests := []struct {
		name          string
		requestBody   []byte
		mockBehaviour func()
		want          want
	}{
		{
			name:        "success",
			requestBody: successRequestBody,
			mockBehaviour: func() {
				authServiceMock.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Return(successAuthHeader, nil).Times(1)
			},
			want: want{
				status:     http.StatusOK,
				authHeader: successAuthHeader,
			},
		},
		{
			name:        "bad request",
			requestBody: badRequestBody,
			mockBehaviour: func() {
				authServiceMock.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Return(successAuthHeader, nil).Times(0)
			},
			want: want{
				status:     http.StatusBadRequest,
				authHeader: "",
			},
		},
		{
			name:        "invalid request",
			requestBody: []byte("dasdsadsa"),
			mockBehaviour: func() {
				authServiceMock.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Return(successAuthHeader, nil).Times(0)
			},
			want: want{
				status:     http.StatusBadRequest,
				authHeader: "",
			},
		},
		{
			name:        "fail register",
			requestBody: successRequestBody,
			mockBehaviour: func() {
				authServiceMock.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Return("", errors.New("internal server error")).Times(1)
			},
			want: want{
				status:     http.StatusInternalServerError,
				authHeader: "",
			},
		},
		{
			name:        "conflict",
			requestBody: successRequestBody,
			mockBehaviour: func() {
				authServiceMock.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Return("", auth.ErrUserAlreadyExists).Times(1)
			},
			want: want{
				status:     http.StatusConflict,
				authHeader: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.mockBehaviour()

			r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(tt.requestBody))
			w := httptest.NewRecorder()

			NewRegisterHandler(authServiceMock).Register(w, r)

			assert.Equal(t, tt.want.status, w.Code)
			assert.Equal(t, tt.want.authHeader, w.Header().Get("Authorization"))
		})
	}
}

package login

import (
	"encoding/json"
	"errors"
	"github.com/anoriar/gophermart/internal/gophermart/dto/requests/login"
	"github.com/anoriar/gophermart/internal/gophermart/handlers/login/internal"
	"github.com/anoriar/gophermart/internal/gophermart/services/auth"
	"io"
	"net/http"
)

type LoginHandler struct {
	authService auth.AuthServiceInterface
	validator   *internal.LoginValidator
}

func NewLoginHandler(authService auth.AuthServiceInterface) *LoginHandler {
	return &LoginHandler{authService: authService, validator: internal.NewLoginValidator()}
}

func (handler *LoginHandler) Login(w http.ResponseWriter, req *http.Request) {
	requestBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	requestDto := &login.LoginUserRequestDto{}
	err = json.Unmarshal(requestBody, requestDto)
	if err != nil {
		if _, ok := err.(*json.SyntaxError); ok {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}
	validationErrors := handler.validator.Validate(*requestDto)
	if len(validationErrors) > 0 {
		http.Error(w, validationErrors.String(), http.StatusBadRequest)
		return
	}

	tokenString, err := handler.authService.LoginUser(req.Context(), *requestDto)
	if err != nil {
		if errors.Is(err, auth.ErrUnauthorized) {
			http.Error(w, "user unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Authorization", tokenString)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

package register

import (
	"encoding/json"
	"errors"
	"github.com/anoriar/gophermart/internal/gophermart/user/dto/requests/register"
	"github.com/anoriar/gophermart/internal/gophermart/user/handlers/register/internal"
	auth2 "github.com/anoriar/gophermart/internal/gophermart/user/services/auth"
	"io"
	"net/http"
)

type RegisterHandler struct {
	registerService auth2.AuthServiceInterface
	validator       *internal.RegisterValidator
}

func NewRegisterHandler(registerService auth2.AuthServiceInterface) *RegisterHandler {
	return &RegisterHandler{registerService: registerService, validator: internal.NewRegisterValidator()}
}

func (handler *RegisterHandler) Register(w http.ResponseWriter, req *http.Request) {
	requestBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	requestDto := &register.RegisterUserRequestDto{}
	err = json.Unmarshal(requestBody, requestDto)
	if err != nil {
		if _, ok := err.(*json.SyntaxError); ok {
			http.Error(w, "invalid json", http.StatusBadRequest)
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

	tokenString, err := handler.registerService.RegisterUser(req.Context(), *requestDto)
	if err != nil {
		if errors.Is(err, auth2.ErrUserAlreadyExists) {
			http.Error(w, "user already exists", http.StatusConflict)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Authorization", tokenString)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

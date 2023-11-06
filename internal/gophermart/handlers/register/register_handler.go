package register

import (
	"encoding/json"
	"errors"
	"github.com/anoriar/gophermart/internal/gophermart/dto/requests/register"
	"github.com/anoriar/gophermart/internal/gophermart/handlers/register/internal"
	"github.com/anoriar/gophermart/internal/gophermart/repository/repository_error"
	"github.com/anoriar/gophermart/internal/gophermart/services/auth"
	"io"
	"net/http"
)

type RegisterHandler struct {
	registerService auth.AuthServiceInterface
	validator       *internal.RegisterValidator
}

func NewRegisterHandler(registerService auth.AuthServiceInterface) *RegisterHandler {
	return &RegisterHandler{registerService: registerService, validator: internal.NewRegisterValidator()}
}

func (handler *RegisterHandler) Register(w http.ResponseWriter, req *http.Request) {
	requestBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	requestDto := &register.RegisterUserRequestDto{}
	err = json.Unmarshal(requestBody, requestDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	validationErrors := handler.validator.Validate(*requestDto)
	if len(validationErrors) > 0 {
		http.Error(w, validationErrors.String(), http.StatusBadRequest)
		return
	}

	tokenString, err := handler.registerService.RegisterUser(req.Context(), *requestDto)
	if err != nil {
		if errors.Is(err, repository_error.ErrConflict) {
			http.Error(w, "Json marshal Error", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Authorization", tokenString)
	w.WriteHeader(http.StatusOK)
}

package register

import (
	"encoding/json"
	"errors"
	"github.com/anoriar/gophermart/internal/gophermart/dto/requests/register"
	"github.com/anoriar/gophermart/internal/gophermart/repository/repository_error"
	"github.com/anoriar/gophermart/internal/gophermart/services/auth"
	"io"
	"net/http"
)

type RegisterHandler struct {
	registerService auth.RegisterServiceInterface
}

func NewRegisterHandler(registerService auth.RegisterServiceInterface) *RegisterHandler {
	return &RegisterHandler{registerService: registerService}
}

func (handler *RegisterHandler) Ping(w http.ResponseWriter, req *http.Request) {
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
	//TODO: validate

	err = handler.registerService.RegisterUser(req.Context(), *requestDto)
	if err != nil {
		if errors.Is(err, repository_error.ErrConflict) {
			http.Error(w, "Json marshal Error", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

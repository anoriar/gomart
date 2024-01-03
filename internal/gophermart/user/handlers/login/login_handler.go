package login

import (
	"encoding/json"
	"errors"
	"github.com/anoriar/gophermart/internal/gophermart/user/dto/requests/login"
	"github.com/anoriar/gophermart/internal/gophermart/user/handlers/login/internal"
	auth2 "github.com/anoriar/gophermart/internal/gophermart/user/services/auth"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"io"
	"net/http"
)

type LoginHandler struct {
	authService auth2.AuthServiceInterface
	validator   *internal.LoginValidator
}

func NewLoginHandler(authService auth2.AuthServiceInterface) *LoginHandler {
	return &LoginHandler{authService: authService, validator: internal.NewLoginValidator()}
}

func (handler *LoginHandler) Login(w http.ResponseWriter, req *http.Request) {
	reqCtx := req.Context()
	span, reqCtx := opentracing.StartSpanFromContext(reqCtx, "LoginHandler::Login")
	defer span.Finish()

	requestBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		span.LogFields(log.Error(err))
		return
	}

	requestDto := &login.LoginUserRequestDto{}
	err = json.Unmarshal(requestBody, requestDto)
	if err != nil {
		if _, ok := err.(*json.SyntaxError); ok {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			span.LogFields(log.Error(err))
		}
		return
	}
	validationErrors := handler.validator.Validate(*requestDto)
	if len(validationErrors) > 0 {
		http.Error(w, validationErrors.String(), http.StatusBadRequest)
		return
	}

	tokenString, err := handler.authService.LoginUser(reqCtx, *requestDto)
	if err != nil {
		if errors.Is(err, auth2.ErrUnauthorized) {
			http.Error(w, "user unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		span.LogFields(log.Error(err))
		return
	}

	w.Header().Add("Authorization", tokenString)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

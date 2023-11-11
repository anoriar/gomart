package internal

import (
	"errors"
	"github.com/anoriar/gophermart/internal/gophermart/dto/requests/login"
	"github.com/anoriar/gophermart/internal/gophermart/dto/validation"
)

type LoginValidator struct {
}

func NewLoginValidator() *LoginValidator {
	return &LoginValidator{}
}

func (validator *LoginValidator) Validate(requestDto login.LoginUserRequestDto) validation.ValidationErrors {
	var validationErrors validation.ValidationErrors

	if requestDto.Login == "" {
		validationErrors = append(validationErrors, errors.New("login required"))
	}

	if requestDto.Password == "" {
		validationErrors = append(validationErrors, errors.New("password required"))
	}

	return validationErrors
}

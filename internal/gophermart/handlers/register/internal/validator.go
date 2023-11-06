package internal

import (
	"errors"
	"fmt"
	"github.com/anoriar/gophermart/internal/gophermart/dto/requests/register"
	"github.com/anoriar/gophermart/internal/gophermart/dto/validation"
)

const (
	loginLength    = 60
	passwordLength = 60
)

type RegisterValidator struct {
}

func NewRegisterValidator() *RegisterValidator {
	return &RegisterValidator{}
}

func (validator *RegisterValidator) Validate(requestDto register.RegisterUserRequestDto) validation.ValidationErrors {
	var validationErrors validation.ValidationErrors

	if requestDto.Login == "" {
		validationErrors = append(validationErrors, errors.New("login required"))
	}

	if requestDto.Password == "" {
		validationErrors = append(validationErrors, errors.New("password required"))
	}

	if len(requestDto.Login) > loginLength {
		validationErrors = append(validationErrors, fmt.Errorf("login must be less than %d symbols", loginLength))
	}

	if len(requestDto.Password) > passwordLength {
		validationErrors = append(validationErrors, fmt.Errorf("login must be less than %d symbols", passwordLength))
	}
	return validationErrors
}

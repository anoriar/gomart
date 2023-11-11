package internal

import (
	"errors"
	"fmt"
	"github.com/anoriar/gophermart/internal/gophermart/dto/requests/register"
	"github.com/anoriar/gophermart/internal/gophermart/dto/validation"
)

const (
	loginMaxLength    = 60
	passwordMaxLength = 60
	loginMinLength    = 4
	passwordMinLength = 6
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

	if len(requestDto.Login) < loginMinLength {
		validationErrors = append(validationErrors, fmt.Errorf("login must be more than %d symbols", loginMinLength))
	}

	if len(requestDto.Password) < passwordMinLength {
		validationErrors = append(validationErrors, fmt.Errorf("pssword must be more than %d symbols", passwordMinLength))
	}

	if len(requestDto.Login) > loginMaxLength {
		validationErrors = append(validationErrors, fmt.Errorf("login must be less than %d symbols", loginMaxLength))
	}

	if len(requestDto.Password) > passwordMaxLength {
		validationErrors = append(validationErrors, fmt.Errorf("password must be less than %d symbols", passwordMaxLength))
	}
	return validationErrors
}

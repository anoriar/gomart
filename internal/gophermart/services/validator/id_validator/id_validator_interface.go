package id_validator

type IdValidatorInterface interface {
	Validate(number string) bool
}

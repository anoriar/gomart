package idvalidator

type IDValidatorInterface interface {
	Validate(number string) bool
}

package password

//go:generate mockgen -source=password_service_interface.go -destination=mock/password_service.go -package=mock
type PasswordServiceInterface interface {
	GenerateHashedPassword(password string, salt []byte) string
	ComparePasswords(password string, hashedPassword string, salt []byte) bool
}

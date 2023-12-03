package password

//go:generate mockgen -source=password_service_interface.go -destination=mock_password_service/password_service.go -package=mock_password_service
type PasswordServiceInterface interface {
	GenerateHashedPassword(password string, salt []byte) string
}

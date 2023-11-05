package services

type PasswordServiceInterface interface {
	GenerateHashedPassword(password string, salt []byte) string
	ComparePasswords(password string, hashedPassword string, salt []byte) bool
}

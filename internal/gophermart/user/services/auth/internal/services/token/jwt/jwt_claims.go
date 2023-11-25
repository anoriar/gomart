package jwt

import "github.com/golang-jwt/jwt/v4"

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID string
}

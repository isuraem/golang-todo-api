package ports

import "github.com/dgrijalva/jwt-go"

type JWTService interface {
	GenerateToken(userID uint) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

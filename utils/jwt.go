package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("your_secret_key") // change in production

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

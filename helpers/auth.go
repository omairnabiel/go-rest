package helpers

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("secret_key")

// HashPassword hashes the the given password
func HashPassword(password string) ([]byte, error) {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), 1)
	if err != nil {
		return nil, err
	}
	return encrypted, nil
}

// GenerateToken takes user info as argument and generates the token
func GenerateToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"email":     email,
		"ExpiresAt": time.Now().Unix(),
	})
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

package helpers

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// JWTClaims inherited with standard claims
type JWTClaims struct {
	Email  string `json:"email"`
	UserID int    `json:"userId"`
	jwt.StandardClaims
}

var jwtKey = []byte("secret_key")

// HashPassword hashes the the given password
func HashPassword(password string) ([]byte, error) {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), 1)
	if err != nil {
		return nil, err
	}
	return encrypted, nil
}

// GenerateToken returns an access and refresh token (accessToken, refreshToken, error)
func GenerateToken(email string) (accessToken string, refreshToken string, err error) {
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"email":     email,
		"ExpiresAt": (time.Now().Add(time.Minute * 15).Unix()),
	}).SignedString(jwtKey)

	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"email":     email,
		"ExpiresAt": time.Now().Add(time.Hour * 48).Unix(),
	}).SignedString(jwtKey)
	return
}

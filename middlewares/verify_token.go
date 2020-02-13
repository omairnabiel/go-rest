package middlewares

import (
	"net/http"
	"strings"

	"github.com/omairnabiel/go-lang-starter/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWTClaims inherited with standard claims
type JWTClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// Will get this from environment variables in real world application
var jwtKey = []byte("secret_key")

// VerifyToken is a middleware verfies the token user sends in the request header
func VerifyToken(ctx *gin.Context) {
	// tokenStr gets the extracted token string from the Request Header
	claims := &JWTClaims{}
	tokenStr := strings.Fields(ctx.Request.Header["Authorization"][0])[1]
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	// if token in not valid throw UnAuthorized Error
	if err != nil || !token.Valid {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorMessage(http.StatusUnauthorized, utils.ErrTokenNotValid))
		ctx.Abort()
		return
	}

	// if no error move to next function
	ctx.Next()
}

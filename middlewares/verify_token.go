package middlewares

import (
	"net/http"

	"github.com/omairnabiel/go-lang-starter/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWTClaims inherited with standard claims
type JWTClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

var jwtKey = []byte("secret_key")

// VerifyToken is a middleware verfies the token user sends in the request header
func VerifyToken(ctx *gin.Context) {

	authorization := ctx.Request.Header["Authorization"]

	// if Authorization key has length not equal to 2 throw an error. Because token format is Authorization: [Bearer, "token"]
	if len(authorization) != 2 {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorMessage(http.StatusUnauthorized, utils.ErrTokenNotValid))
		ctx.Abort()
		return
	}
	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(authorization[1], claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorMessage(http.StatusUnauthorized, utils.ErrTokenNotValid))
		ctx.Abort()
		return
	}

	ctx.Next()
}

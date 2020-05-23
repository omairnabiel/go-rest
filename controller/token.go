package controller

import (
	"context"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/omairnabiel/go-lang-starter/cache"
	"github.com/omairnabiel/go-lang-starter/helpers"
)

type refreshTokenRequest struct {
	RefreshToken string
	AccessToken  string
}

// RefreshToken controller returns user a new access_token
func RefreshToken(ctx *gin.Context) {

	accessToken, aErr := ctx.Request.Cookie("access_token")
	refreshToken, rErr := ctx.Request.Cookie("refresh_token")

	if aErr != nil && rErr != nil {
		ctx.AbortWithStatus(401)
	}

	accessTokenClaims, err := verifyToken(accessToken.Value, "secret")
	if err != nil {
		ctx.AbortWithError(401, err)
		return
	}

	_, err = verifyToken(refreshToken.Value, "secret")
	if err != nil {
		ctx.AbortWithError(401, err)
		return
	}

	// check if the refreshToken sent by user has not been revoked and exists in array of active user tokens
	userTokens := cache.Redis.LRange(context.Background(), accessTokenClaims.Id, 0, -1).Val()
	isFound := helpers.Contains(userTokens, refreshToken.Name)
	if !isFound {
		ctx.AbortWithStatus(401)
		return
	}

	accessToken.Value, refreshToken.Value, err = helpers.GenerateToken("email")
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	ctx.Request.AddCookie(accessToken)
	ctx.Request.AddCookie(refreshToken)
	ctx.JSON(http.StatusOK, helpers.SuccessMessage(http.StatusOK, "Succesfully Set Tokens", nil))
}

// verifyToken verifies if the token is valid and has expired. If the token is valid and has expired it returns true
func verifyToken(requestToken string, tokenSecret string) (helpers.JWTClaims, error) {
	claims := &helpers.JWTClaims{}
	_, err := jwt.ParseWithClaims(requestToken, claims, func(token *jwt.Token) (interface{}, error) {
		return tokenSecret, nil
	})

	if err != nil {
		return *claims, err
	}

	return *claims, nil
}

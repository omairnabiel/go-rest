package controller

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/omairnabiel/go-lang-starter/cache"
	"github.com/omairnabiel/go-lang-starter/helpers"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
	"github.com/dgrijalva/jwt-go"
)

// LoginRequest maps to login api call params
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=50"`
}

// LoginResponse object to send after successful login of user
type LoginResponse struct {
	Email        string `json:"email"`
	Name         string `json:"name"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// Login function to execute on signup route. Checks if user doesn't exist then add the user in the DB
func Login(ctx *gin.Context) {
	var creds LoginRequest

	if err := ctx.ShouldBindJSON(&creds); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.ErrorMessage(http.StatusBadRequest, err.Error()))
		return
	}

	v := validator.New()

	if err := v.Struct(creds); err != nil {
		var errors []string
		for _, e := range err.(validator.ValidationErrors) {
			log.Println("Errors", e)
			errors = append(errors, helpers.ValidationMessage(e.Field(), e.Tag(), e.Param()))
		}
		var status []int
		status = append(status, http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, helpers.ErrorMessages(status, errors))
		return
	}

	redisErr := cache.Redis.LPush(context.Background(), "user", "Value", 0).Err()
	if redisErr != nil {
		ctx.JSON(http.StatusNotFound, helpers.ErrorMessage(http.StatusNotFound, helpers.ErrUserDoesntExist))
		return
	}

	cachedUser, found := cache.Redis.Get(context.Background(), "user").Result()

	if found != nil {
		ctx.JSON(http.StatusNotFound, helpers.ErrorMessage(http.StatusNotFound, helpers.ErrUserDoesntExist))
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(cachedUser), []byte(creds.Password))

	if err != nil {
		ctx.JSON(http.StatusForbidden, helpers.ErrorMessage(http.StatusForbidden, helpers.ErrIncorrectPassword))
		return
	}

	accessToken, refreshToken, err := helpers.GenerateToken(creds.Email)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.ErrorMessage(http.StatusInternalServerError, helpers.ErrInternalServerError))
		return
	}

	var resp interface{}
	resp = &LoginResponse{Email: creds.Email, AccessToken: accessToken, RefreshToken: refreshToken, Name: "user.Name"}

	ctx.JSON(http.StatusOK, helpers.SuccessMessage(http.StatusOK, helpers.SuccessUserLogin, resp))
}



// LogoutRequest maps to logout api call params
type LogoutRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// Logout route
func Logout(ctx *gin.Context) {
	var body LogoutRequest

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.ErrorMessage(http.StatusBadRequest, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, helpers.SuccessMessage(http.StatusOK, helpers.SuccessUserLoggedOut, nil))

}

// SignUpRequest , request Object to be recieved from client
type SignUpRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=2,max=50"`
	Password string `json:"password" validate:"required,min=8,max=50"`
}

// Signup function to execute on signup route. Checks if user doesn't exist then add the user in the DB
func Signup(ctx *gin.Context) {
	var body SignUpRequest

	v := validator.New()

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.ErrorMessage(http.StatusBadRequest, err.Error()))
		return
	}

	err := v.Struct(body)

	if err != nil {
		var errors []string
		for _, e := range err.(validator.ValidationErrors) {
			log.Println("Errors", e.Field(), e.Tag(), e.Param())
			errors = append(errors, helpers.ValidationMessage(e.Field(), e.Tag(), e.Param()))
		}
		var status []int
		status = append(status, http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, helpers.ErrorMessages(status, errors))
		return
	}

	encrypted, err := helpers.HashPassword(body.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.ErrorMessage(http.StatusInternalServerError, helpers.ErrInternalServerError))
		return
	}
	user := body
	user.Password = string(encrypted)
	ctx.JSON(http.StatusOK, helpers.SuccessMessage(http.StatusOK, helpers.SuccessUserCreated, nil))
}

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

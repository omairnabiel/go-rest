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

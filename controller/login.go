package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/omairnabiel/go-lang-starter/cache"
	"github.com/omairnabiel/go-lang-starter/helpers"
	"github.com/omairnabiel/go-lang-starter/utils"
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
	Email string `json:"email"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

// Login function to execute on signup route. Checks if user doesn't exist then add the user in the DB
func Login(ctx *gin.Context) {
	var creds LoginRequest

	if err := ctx.ShouldBindJSON(&creds); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorMessage(http.StatusBadRequest, err.Error()))
		return
	}

	v := validator.New()

	errValid := v.Struct(creds)

	if errValid != nil {
		var errors []string
		for _, e := range errValid.(validator.ValidationErrors) {
			log.Println("Errors", e)
			errors = append(errors, utils.ValidationMessage(e.Field(), e.Tag(), e.Param()))
		}
		var status []int
		status = append(status, http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, utils.ErrorMessages(status, errors))
		return
	}

	redisErr := cache.Redis.Set(cache.Redis.Context(), "user", "Value", 0).Err()
	if redisErr != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorMessage(http.StatusNotFound, utils.ErrUserDoesntExist))
		return
	}

	cachedUser, found := cache.Redis.Get(cache.Redis.Context(), "user").Result()

	if found != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorMessage(http.StatusNotFound, utils.ErrUserDoesntExist))
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(cachedUser), []byte(creds.Password))

	if err != nil {
		ctx.JSON(http.StatusForbidden, utils.ErrorMessage(http.StatusForbidden, utils.ErrIncorrectPassword))
		return
	}

	token, err := helpers.GenerateToken(creds.Email)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorMessage(http.StatusInternalServerError, utils.ErrInternalServerError))
		return
	}

	var resp interface{}
	resp = &LoginResponse{Email: creds.Email, Token: token, Name: "user.Name"}

	ctx.JSON(http.StatusOK, utils.SuccessMessage(http.StatusOK, utils.SuccessUserLogin, resp))
}

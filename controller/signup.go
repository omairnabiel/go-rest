package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/omairnabiel/go-lang-starter/helpers"
	"github.com/omairnabiel/go-lang-starter/utils"
	"gopkg.in/go-playground/validator.v9"
)

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
		ctx.JSON(http.StatusBadRequest, utils.ErrorMessage(http.StatusBadRequest, err.Error()))
		return
	}

	err := v.Struct(body)

	if err != nil {
		var errors []string
		for _, e := range err.(validator.ValidationErrors) {
			log.Println("Errors", e.Field(), e.Tag(), e.Param())
			errors = append(errors, utils.ValidationMessage(e.Field(), e.Tag(), e.Param()))
		}
		var status []int
		status = append(status, http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, utils.ErrorMessages(status, errors))
		return
	}

	encrypted, err := helpers.HashPassword(body.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorMessage(http.StatusInternalServerError, utils.ErrInternalServerError))
		return
	}
	user := body
	user.Password = string(encrypted)
	ctx.JSON(http.StatusOK, utils.SuccessMessage(http.StatusOK, utils.SuccessUserCreated, nil))
}

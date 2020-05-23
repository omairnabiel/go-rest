package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/omairnabiel/go-lang-starter/helpers"
)

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

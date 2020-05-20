package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/omairnabiel/go-lang-starter/utils"
)

// LogoutRequest maps to logout api call params
type LogoutRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// Logout route
func Logout(ctx *gin.Context) {
	var body LogoutRequest

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorMessage(http.StatusBadRequest, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.SuccessMessage(http.StatusOK, utils.SuccessUserLoggedOut, nil))

}

package helpers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// RespErr response
type RespErr struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// ErrorMessages is used to when you have multiple errors and want to return all of them
func ErrorMessages(status []int, message []string) gin.H {

	var err []RespErr
	for index := range status {
		err = append(err, RespErr{Status: status[index], Message: message[index]})
	}
	resp := map[string]interface{}{
		"error": err,
	}
	return resp
}

// ErrorMessage is used in case of single error, returns error response object
func ErrorMessage(status int, message string) gin.H {
	resp := map[string]interface{}{
		"status":  status,
		"message": message,
	}
	return resp
}

// ValidationMessage recieves a field name, type of error and expectedValue and constructs a validation message
func ValidationMessage(field string, tag string, expectedValue string) string {
	var message string
	switch tag {
	case "min":
		message = fmt.Sprintf("'%s' must be atleast %s characters", field, expectedValue)
	case "max":
		message = fmt.Sprintf("'%s' must be less than %s characters", field, expectedValue)
	case "email":
		message = fmt.Sprintf("'%s' is invalid", field)
	case "required":
		message = fmt.Sprintf("'%s' is a required field", field)
	}
	return message
}

var (
	// ErrInvalidParameters message
	ErrInvalidParameters = "Invalid Parameters"

	// ErrPermissionDenied message
	ErrPermissionDenied = "Permission Denied"

	// ErrUserExists message
	ErrUserExists = "User Already Exists"

	// ErrUserCreate message
	ErrUserCreate = "Failed creating User. Please try again"

	// ErrUserDoesntExist message
	ErrUserDoesntExist = "User doesn't exist. Please signup"

	// ErrInternalServerError message
	ErrInternalServerError = "Internal Server Error"

	// ErrIncorrectPassword message
	ErrIncorrectPassword = "Incorrect Password"

	// ErrLogin message
	ErrLogin = "Failed to login. Please Try Again"

	// ErrAlreadyLoggedOut message
	ErrAlreadyLoggedOut = "You're already logged out"

	// ErrTokenNotValid message
	ErrTokenNotValid = "Token in not valid!"
)

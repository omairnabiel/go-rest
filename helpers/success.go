package helpers

import "github.com/gin-gonic/gin"

// SuccessMessage is used in case of single error, returns error response object
func SuccessMessage(status int, message string, data interface{}) gin.H {
	resp := map[string]interface{}{
		"status":  status,
		"message": message,
	}

	if data == nil {
		return resp
	}
	resp["data"] = data
	return resp
}

var (
	// SuccessUserCreated message
	SuccessUserCreated = "User Created Successfuly"

	// SuccessUserLoggedOut message
	SuccessUserLoggedOut = "User Logged Out Successfuly"

	// SuccessUserLogin message
	SuccessUserLogin = "User Logged In Successfuly"
)

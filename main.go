package main

import (
	"github.com/gin-gonic/gin"
	auth "github.com/omairnabiel/go-lang-starter/auth"
	middlewares "github.com/omairnabiel/go-lang-starter/middlewares"
)

func main() {

	// create new gin router
	router := gin.New()

	// middlewares
	router.Use(gin.Logger())

	// user on-boarding auth routes
	router.POST("/login", auth.Login)
	router.POST("/signup", auth.Signup)

	router.Use(middlewares.VerifyToken)

	router.POST("/logout", auth.Logout)

	// run the server on port 8080
	router.Run(":8080")
}

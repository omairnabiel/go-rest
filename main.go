package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"                            // postgres driver
	_ "github.com/omairnabiel/go-lang-starter/cache" // initialize cache
	"github.com/omairnabiel/go-lang-starter/controller"
	_ "github.com/omairnabiel/go-lang-starter/db" // initialize db
	"github.com/omairnabiel/go-lang-starter/middlewares"
)

func main() {

	// create new gin router
	router := gin.New()

	// middlewares
	router.Use(gin.Logger())

	// user on-boarding auth routes
	router.POST("/login", controller.Login)
	router.POST("/signup", controller.Signup)

	router.Use(middlewares.VerifyToken)

	router.POST("/logout", controller.Logout)

	// run the server on port 8080
	router.Run(":8080")
}

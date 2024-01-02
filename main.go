package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jwt-auth/first-try/controllers"
	"github.com/jwt-auth/first-try/initializers"
	"github.com/jwt-auth/first-try/middlewares"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()

	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.LogIn)
	r.GET("/validate", middlewares.RequireAuth, controllers.Validate)

	r.Run()
	fmt.Println("Hello 2")
}

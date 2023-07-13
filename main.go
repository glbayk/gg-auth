package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/glbayk/gg-auth/controllers"
	"github.com/glbayk/gg-auth/middleware"
	"github.com/glbayk/gg-auth/models"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		os.Exit(1)
	}

	models.Connect()
}

func main() {
	router := gin.Default()

	router.Use(cors.Default())

	// Health check endpoint
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "pong"})
	})

	aC := controllers.AuthController{}
	auth := router.Group("/auth")
	auth.POST("/register", aC.Register)
	auth.POST("/login", aC.Login)
	auth.GET("/refresh-token", aC.RefreshToken)
	auth.POST("/reset-password", aC.ResetPassword)

	uC := controllers.UserController{}
	user := router.Group("/user")
	user.Use(middleware.Authenticated())
	user.GET("/profile", uC.Profile)

	// PORT is set in .env file
	router.Run()
	log.Println("Server running on port " + os.Getenv("PORT"))
}

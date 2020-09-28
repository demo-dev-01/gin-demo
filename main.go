package main

import (
	"github.com/godemo-dev/gin-demo/model"
	"github.com/godemo-dev/gin-demo/pkg"

	"github.com/gin-gonic/gin"
)

func main() {
	model.Setup()

	gin.ForceConsoleColor()
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Static("/upload", "./upload")
	r.Use(pkg.CORS())
	router := r.Group("api/v1")

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// register
	router.POST("users/register", pkg.Register)
	// dang ky
	router.POST("users/login", pkg.Login)
	// JWT
	router.Use(pkg.JWT())
	{
		// File upload
		router.POST("/upload", pkg.UploadFile)
		// Get all User
		router.GET("/users", pkg.GetUsers)
		// Update user
		router.PUT("/users", pkg.UpdateUser)
		// Get user
		router.GET("/users/:id", pkg.GetUserId)
		// Delete user
		router.DELETE("/users/:id", pkg.DeleteUserById)

	}

	r.Run()
}

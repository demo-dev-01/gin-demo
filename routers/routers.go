package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/godemo-dev/gin-demo/jwt"
	v1 "github.com/godemo-dev/gin-demo/routers/v1"
)

func InitRouters() {
	gin.ForceConsoleColor()
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Static("/upload", "./upload")
	r.Use(jwt.CORS())
	router := r.Group("api/v1")

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// register
	router.POST("users/register", v1.Register)
	// dang ky
	router.POST("users/login", v1.Login)
	// JWT
	router.Use(jwt.JWT())
	{
		// File upload
		router.POST("/upload", v1.UploadFile)
		// Get all User
		router.GET("/users", v1.GetUsers)
		// Update user
		router.PUT("/users", v1.UpdateUser)
		// Get user
		router.GET("/users/:id", v1.GetUserId)
		// Delete user
		router.DELETE("/users/:id", v1.DeleteUserById)
	}

	r.Run()
}

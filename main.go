package main

import (
	"fmt"
	"gin-api/model"
	"gin-api/pkg"
	"net/http"
	"path/filepath"
	"strconv"

	"gorm.io/gorm"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

var db *gorm.DB

type UserLogin struct {
	Email    string
	Password string
}

type GetUser struct {
	FirstName string
	LastName  string
	Email     string
}

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data interface{}

		tokenToString := c.Request.Header.Get("Authorization")
		code := http.StatusOK

		if tokenToString == "" {
			code = http.StatusUnauthorized
		}
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenToString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("jdnfksdmfksd"), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  "Unauthorized",
				"data": data,
			})
			c.Abort()
			return
		}
		if token == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  "Unauthorized",
				"data": data,
			})
			c.Abort()
			return
		}

		if code != http.StatusOK {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  "Unauthorized",
				"data": data,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	model.Setup()

	gin.ForceConsoleColor()
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Static("/upload", "./upload")
	r.Use(CORS())
	router := r.Group("api/v1")

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// register
	router.POST("users/register", pkg.Register)

	// dang ky
	router.POST("users/login", func(c *gin.Context) {
		var userLogin UserLogin
		var user model.User
		c.BindJSON(&userLogin)

		if err := db.Where("email = ?", userLogin.Email).First(&user).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Email nout found",
			})
			return
		}

		check := pkg.CheckPasswordHash(userLogin.Password, user.Password)
		if !check {
			c.JSON(http.StatusBadRequest, gin.H{
				"messger": "Email or password does not exist",
			})
			return
		}

		token, err := pkg.CreateToken(user.Email, user.ID)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": "token error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"messager": "Success",
			"token":    token,
		})
	})

	router.Use(JWT())
	{
		// Set a lower memory limit for multipart forms (default is 32 MiB)
		router.POST("/upload", func(c *gin.Context) {
			// Source
			file, err := c.FormFile("file")
			if err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
				return
			}

			filename := "upload/" + filepath.Base(file.Filename)

			if err := c.SaveUploadedFile(file, filename); err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
				return
			}

			c.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully with fields", file.Filename))
		})

		router.GET("/users", func(c *gin.Context) {
			user := []model.User{}
			db.Find(&user)

			c.JSON(http.StatusOK, gin.H{
				"message": "Success",
				"result":  user,
			})
			return
		})

		router.PUT("/users", func(c *gin.Context) {
			var user model.User
			c.BindJSON(&user)
			a := user

			if err := db.Where("email = ?", user.Email).First(&user).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "Email nout found",
				})
				return
			}
			user.FirstName = a.FirstName
			user.LastName = a.LastName

			db.Save(&user)

			c.JSON(http.StatusOK, gin.H{
				"message": "Success",
			})
			return
		})

		router.GET("/users/:id", func(c *gin.Context) {
			var user model.User
			id := c.Param("id")
			if err := db.Where("id = ?", &id).First(&user).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "user nout found",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Success",
				"result": GetUser{
					FirstName: user.FirstName,
					LastName:  user.LastName,
					Email:     user.Email,
				},
			})
			return
		})

		router.DELETE("/users/:id", func(c *gin.Context) {
			id := c.Param("id")
			if id == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "User not found",
				})
				return
			}
			strconv.Atoi(id)
			db.Delete(&model.User{}, id)

			c.JSON(http.StatusOK, gin.H{
				"message": "Success",
			})
			return
		})

	}

	r.Run()
}

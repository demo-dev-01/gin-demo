package jwt

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// Hash password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Compare password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateToken(email string, id uint) (string, error) {
	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["email"] = email
	atClaims["id"] = id
	atClaims["exp"] = time.Now().Add(time.Minute * 1).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
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

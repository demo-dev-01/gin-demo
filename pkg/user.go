package pkg

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var user User
	c.BindJSON(&user)
	password, err := HashPassword(user.Password)
	if err != nil {
		return
	}

	user.Password = password

	CreateUser(user)

	c.JSON(http.StatusOK, gin.H{
		"messager": "Success",
	})

}

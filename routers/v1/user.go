package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/godemo-dev/gin-demo/jwt"
	"github.com/godemo-dev/gin-demo/model"
)

func Register(c *gin.Context) {
	var user model.User
	c.BindJSON(&user)
	password, err := jwt.HashPassword(user.Password)
	if err != nil {
		return
	}

	user.Password = password

	model.CreateUser(user)

	c.JSON(http.StatusOK, gin.H{
		"messager": "Success",
	})

}

func Login(c *gin.Context) {
	var userLogin model.UserLogin
	c.BindJSON(&userLogin)
	user, err := model.FindByEmail(userLogin.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Email nout found",
		})
		return
	}

	check := jwt.CheckPasswordHash(userLogin.Password, user.Password)
	if !check {
		c.JSON(http.StatusBadRequest, gin.H{
			"messger": "Email or password does not exist",
		})
		return
	}

	token, err := jwt.CreateToken(user.Email, user.ID)

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
}

type GetUser struct {
	FirstName string
	LastName  string
	Email     string
}

func GetUserId(c *gin.Context) {
	id := c.Param("id")
	user, err := model.FindById(id)

	if err != nil {
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
}

func DeleteUserById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User not found",
		})
		return
	}

	strconv.Atoi(id)
	err := model.DeleteById(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
	})
	return
}

func GetUsers(c *gin.Context) {
	list, err := model.GetAllUser()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"result":  list,
	})
}

func UpdateUser(c *gin.Context) {
	var user model.User
	c.BindJSON(&user)

	a, err := model.FindByEmail(user.Email)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Email nout found",
		})
		return
	}

	a.FirstName = user.FirstName
	a.LastName = user.LastName

	userErr := model.UpdateByUser(a)

	if userErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Update false",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
	})
	return
}

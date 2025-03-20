package controller

import (
	"funny-login/model"
	"funny-login/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type User struct {
	RG *gin.RouterGroup
}

func (u *User) createUser(c *gin.Context) {
	var payload model.User
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	user, err := usecase.CreateUser(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (u *User) getAllUser(c *gin.Context) {
	users, err := usecase.ListAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	if len(users) > 0 {
		c.JSON(http.StatusOK, users)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "List user empty"})
}

func (u *User) getUserById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := usecase.GetUserById(uint32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (u *User) Route() {

	u.RG.POST("/users", u.createUser)
	u.RG.GET("/users", u.getAllUser)
	u.RG.GET("/users/:id", u.getUserById)

}

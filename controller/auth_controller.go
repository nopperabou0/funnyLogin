package controller

import (
	"fmt"
	"funny-login/model"
	"funny-login/usecase"

	"github.com/gin-gonic/gin"
)

type Auth struct {
	RG *gin.RouterGroup
}

func (a *Auth) loginHandler(c *gin.Context) {
	var payload model.User
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"err": fmt.Errorf("authentication failed : %s", err.Error())})
		return
	}

	token, err := usecase.Login(payload.Name, payload.Password)
	if err != nil {
		c.JSON(500, gin.H{"err": err})
		return
	}
	c.JSON(201, gin.H{"token": token})
}

func (a *Auth) Route() {
	a.RG.POST("/login", a.loginHandler)
}

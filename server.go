package main

import (
	"fmt"
	"funny-login/config"
	"funny-login/controller"
	"funny-login/repository"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	var engine *gin.Engine = gin.Default()
	var conf config.Config
	var err error
	repository.DB, err = conf.DB()
	if err != nil {
		panic(fmt.Errorf("connection error : %v", err.Error()))
	}

	defer repository.CloseDB()

	rg := engine.Group("/api/v1")

	var userController = controller.User{
		RG: rg,
	}
	host := conf.API()
	userController.Route()
	err = engine.Run(host)
	if err != nil {
		panic(fmt.Errorf("server not running on host %s, because error %v", host, err.Error()))
	}
}

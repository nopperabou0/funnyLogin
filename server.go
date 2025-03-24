package main

import (
	"fmt"
	"funny-login/config"
	"funny-login/controller"
	"funny-login/repository"
	jwtservice "funny-login/utils/jwt_service"

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

	jwtservice.Conf = *conf.Token()

	rg := engine.Group("/api/v1")
	var userController = controller.User{
		RG: rg,
	}
	var authController = controller.Auth{
		RG: rg,
	}
	userController.Route()
	authController.Route()

	host := conf.API()
	err = engine.Run(host)

	if err != nil {
		panic(fmt.Errorf("server not running on host %s, because error %v", host, err.Error()))
	}
}

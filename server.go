package main

import (
	"fmt"
	"funny-login/config"
	"funny-login/controller"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	var engine *gin.Engine = gin.Default()
	var conf config.Config

	db := conf.DB()
	rg := engine.Group("/api/v1")

	var userC = controller.User{
		DB: db,
		RG: rg,
	}
	host := fmt.Sprintf(":%s", conf.API())
	userC.Route()
	err := engine.Run(host)
	if err != nil {
		panic(fmt.Errorf("server not running on host %s, because error %v", host, err.Error()))
	}
}

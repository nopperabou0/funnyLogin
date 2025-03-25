package main

import (
	"database/sql"
	"fmt"

	"enigmacamp.com/unit-test-starter-pack/middleware"
	"enigmacamp.com/unit-test-starter-pack/utils/service"

	"enigmacamp.com/unit-test-starter-pack/config"
	"enigmacamp.com/unit-test-starter-pack/controller"
	"enigmacamp.com/unit-test-starter-pack/repository"
	"enigmacamp.com/unit-test-starter-pack/usecase"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Server struct {
	userUC     usecase.UserUseCase
	authUc     usecase.AuthenticationUseCase
	jwtService service.JwtService
	engine     *gin.Engine
	host       string
}

func (s *Server) initRoute() {
	rg := s.engine.Group("/api/v1")
	authMiddleware := middleware.NewAuthMiddleware(s.jwtService)
	controller.NewUserController(s.userUC, rg, authMiddleware).Route()
	controller.NewAuthController(s.authUc, rg).Route()
}

func (s *Server) Run() {
	s.initRoute()
	if err := s.engine.Run(s.host); err != nil {
		panic(fmt.Errorf("server not running on host %s, becauce error %v", s.host, err.Error()))
	}
}

func NewServer() *Server {
	cfg, _ := config.NewConfig()
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database)
	db, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		panic("connection error")
	}
	userRepo := repository.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo)
	jwtService := service.NewJwtService(cfg.TokenConfig)
	authUseCase := usecase.NewAuthenticationUseCase(userUseCase, jwtService)
	engine := gin.Default()
	host := fmt.Sprintf(":%s", cfg.ApiPort)
	return &Server{
		userUC:     userUseCase,
		authUc:     authUseCase,
		jwtService: jwtService,
		engine:     engine,
		host:       host,
	}
}

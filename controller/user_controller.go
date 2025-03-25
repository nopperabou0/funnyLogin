package controller

import (
	"net/http"
	"strconv"

	"enigmacamp.com/unit-test-starter-pack/middleware"
	"enigmacamp.com/unit-test-starter-pack/model"
	"enigmacamp.com/unit-test-starter-pack/usecase"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	useCase        usecase.UserUseCase
	rg             *gin.RouterGroup
	authMiddleware middleware.AuthMiddleware
}

func (b *UserController) createUser(c *gin.Context) {
	var payload model.UserCredential
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	book, err := b.useCase.RegisterNewUser(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, book)
}

func (b *UserController) getAllUser(c *gin.Context) {
	books, err := b.useCase.FindAllUser()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to retrieve data users"})
		return
	}

	if len(books) > 0 {
		c.JSON(http.StatusOK, books)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "List user empty"})
}

func (b *UserController) getUserById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	book, err := b.useCase.FindUserById(uint32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to get user by ID"})
		return
	}

	c.JSON(http.StatusOK, book)
}

func (b *UserController) Route() {
	b.rg.POST("/users", b.authMiddleware.RequireToken("admin"), b.createUser)
	b.rg.GET("/users", b.authMiddleware.RequireToken("admin"), b.getAllUser)
	b.rg.GET("/users/:id", b.authMiddleware.RequireToken("admin"), b.getUserById)
}

func NewUserController(useCase usecase.UserUseCase, rg *gin.RouterGroup, am middleware.AuthMiddleware) *UserController {
	return &UserController{useCase: useCase, rg: rg, authMiddleware: am}
}

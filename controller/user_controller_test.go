package controller

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"enigmacamp.com/unit-test-starter-pack/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserUseCaseMock struct {
	mock.Mock
}
type AuthMiddlewareMock struct {
	mock.Mock
}

func (u *UserUseCaseMock) RegisterNewUser(payload model.UserCredential) (model.UserCredential, error) {
	args := u.Called(payload)
	return args.Get(0).(model.UserCredential), args.Error(1)
}

func (u *UserUseCaseMock) FindAllUser() ([]model.UserCredential, error) {
	args := u.Called()
	return args.Get(0).([]model.UserCredential), args.Error(1)
}

func (u *UserUseCaseMock) FindUserById(id uint32) (model.UserCredential, error) {
	args := u.Called(id)
	return args.Get(0).(model.UserCredential), args.Error(1)
}

func (u *UserUseCaseMock) FindUserByUsernamePassword(username string, password string) (model.UserCredential, error) {
	args := u.Called(username, password)
	return args.Get(0).(model.UserCredential), args.Error(1)
}

func (a *AuthMiddlewareMock) RequireToken(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

type UserControllerTestSuite struct {
	suite.Suite
	rg                 *gin.RouterGroup
	userUseCaseMock    *UserUseCaseMock
	authMiddlewareMock *AuthMiddlewareMock
}

func (suite *UserControllerTestSuite) SetupTest() {
	suite.userUseCaseMock = new(UserUseCaseMock)
	suite.authMiddlewareMock = new(AuthMiddlewareMock)
	r := gin.Default()
	gin.SetMode(gin.TestMode)
	rg := r.Group("/api/v1")
	r.Use(suite.authMiddlewareMock.RequireToken("admin"))
	suite.rg = rg
}

var dummyData = model.UserCredential{
	Username: "Maher",
	Password: "1234",
	Role:     "admin",
}

var expectedUser = model.UserCredential{
	Id:       "1",
	Username: "Maher",
	Role:     "admin",
}

// this is whatever even empty string it will pass
var mockTokenJWT string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30"

func (suite *UserControllerTestSuite) TestCreateUser_Success() {
	suite.userUseCaseMock.On("RegisterNewUser", dummyData).Return(expectedUser, nil)
	userController := NewUserController(suite.userUseCaseMock, suite.rg, suite.authMiddlewareMock)
	userController.Route()

	userJSON := `{
		"username": "Maher",
		"password": "1234",
		"role": "admin"
	}`

	request, err := http.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(userJSON))
	assert.NoError(suite.T(), err)

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+mockTokenJWT)

	record := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = request

	userController.createUser(ctx)

	assert.Equal(suite.T(), http.StatusCreated, record.Code)
}

func (suite *UserControllerTestSuite) TestCreateUser_BadRequest() {
	// arrange
	suite.userUseCaseMock.AssertNotCalled(suite.T(), "RegisterNewUser")
	// inialisasi controller dan route
	userController := NewUserController(suite.userUseCaseMock, suite.rg, suite.authMiddlewareMock)
	userController.Route()
	// payload
	userJSON := `{
	 "username": "Maher",
	 "password": "1234",
 }`
	request, err := http.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(userJSON))
	assert.NoError(suite.T(), err)

	// set header
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+mockTokenJWT)
	// membuat recorder untuk mencatat respon
	record := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = request

	userController.createUser(ctx)

	assert.Equal(suite.T(), http.StatusBadRequest, record.Code)
}

func (suite *UserControllerTestSuite) TestCreateUser_Failed() {
	suite.userUseCaseMock.On("RegisterNewUser", dummyData).Return(model.UserCredential{}, fmt.Errorf("Some error occured"))

	userController := NewUserController(suite.userUseCaseMock, suite.rg, suite.authMiddlewareMock)
	userController.Route()

	userJSON := `{
		"username": "Maher",
		"password": "1234",
		"role": "admin"
	}`

	request, err := http.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(userJSON))
	assert.NoError(suite.T(), err)

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+mockTokenJWT)

	record := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = request

	userController.createUser(ctx)

	assert.Equal(suite.T(), http.StatusInternalServerError, record.Code)
	suite.userUseCaseMock.AssertCalled(suite.T(), "RegisterNewUser", dummyData)
}

func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}

package usecase

import (
	"fmt"
	"strconv"
	"testing"

	"enigmacamp.com/unit-test-starter-pack/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (u *UserRepositoryMock) Create(payload model.UserCredential) (model.UserCredential, error) {
	args := u.Called(payload)
	return args.Get(0).(model.UserCredential), args.Error(1)
}

func (u *UserRepositoryMock) List() ([]model.UserCredential, error) {
	args := u.Called()
	return args.Get(0).([]model.UserCredential), args.Error(1)
}

func (u *UserRepositoryMock) Get(id uint32) (model.UserCredential, error) {
	args := u.Called(id)
	return args.Get(0).(model.UserCredential), args.Error(1)
}

func (u *UserRepositoryMock) GetByUsernamePassword(username string, password string) (model.UserCredential, error) {
	args := u.Called(username, password)
	return args.Get(0).(model.UserCredential), args.Error(1)
}

type UserUseCaseTestSuite struct {
	suite.Suite
	userRepoMock *UserRepositoryMock
	userUseCase  UserUseCase
}

func (suite *UserUseCaseTestSuite) SetupTest() {
	suite.userRepoMock = new(UserRepositoryMock)
	suite.userUseCase = NewUserUseCase(suite.userRepoMock)
}

var expectedUser = model.UserCredential{
	Id:       "1",
	Username: "Maher",
	Password: "1234",
	Role:     "admin",
}

var expectedUsers = []model.UserCredential{
	{
		Id:       "1",
		Username: "Maher",
		Role:     "admin",
	},
	{
		Id:       "2",
		Username: "Fauzi",
		Role:     "user",
	},
}

func (suite *UserUseCaseTestSuite) TestCreate_Success() {
	suite.userRepoMock.On("Create", expectedUser).Return(expectedUser, nil)
	actualData, err := suite.userUseCase.RegisterNewUser(expectedUser)

	assert.Nil(suite.T(), err)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedUser, actualData)
	suite.userRepoMock.AssertCalled(suite.T(), "Create", expectedUser)
}

func (suite *UserUseCaseTestSuite) TestCreate_Failed() {
	suite.userRepoMock.On("Create", model.UserCredential{}).Return(model.UserCredential{}, fmt.Errorf("failed to register user"))

	actualData, err := suite.userUseCase.RegisterNewUser(model.UserCredential{})

	assert.NotNil(suite.T(), err)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), model.UserCredential{}, actualData)
	suite.userRepoMock.AssertCalled(suite.T(), "Create", model.UserCredential{})
}

func (suite *UserUseCaseTestSuite) TestList_Success() {
	suite.userRepoMock.On("List").Return(expectedUsers, nil)
	actualData, err := suite.userUseCase.FindAllUser()

	assert.Nil(suite.T(), err)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedUsers, actualData)
	suite.userRepoMock.AssertCalled(suite.T(), "List")
}

func (suite *UserUseCaseTestSuite) TestList_Failed() {
	suite.userRepoMock.On("List").Return([]model.UserCredential{}, fmt.Errorf("Failed to find all user\n"))
	actualData, err := suite.userUseCase.FindAllUser()

	assert.NotNil(suite.T(), err)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), []model.UserCredential{}, actualData)
	suite.userRepoMock.AssertCalled(suite.T(), "List")
	suite.userRepoMock.AssertExpectations(suite.T())
}

func (suite *UserUseCaseTestSuite) TestGet_Success() {
	id, _ := strconv.Atoi(expectedUser.Id)
	suite.userRepoMock.On("Get", uint32(id)).Return(expectedUser, nil)
	actualData, err := suite.userUseCase.FindUserById(uint32(id))
	assert.Nil(suite.T(), err)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedUser, actualData)
	suite.userRepoMock.AssertCalled(suite.T(), "Get", uint32(id))
}

func (suite *UserUseCaseTestSuite) TestGet_Failed() {
	suite.userRepoMock.On("Get", uint32(2)).Return(model.UserCredential{}, fmt.Errorf("Failed to find user by id"))
	actualData, err := suite.userUseCase.FindUserById(uint32(2))
	assert.NotNil(suite.T(), err)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), model.UserCredential{}, actualData)
	suite.userRepoMock.AssertCalled(suite.T(), "Get", uint32(2))

}

func TestUserUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserUseCaseTestSuite))
}

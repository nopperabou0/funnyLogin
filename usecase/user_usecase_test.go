package usecase

import (
	"database/sql"
	"funny-login/model"
	"funny-login/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (u *UserRepositoryMock) Create(params *repository.Params) (model.User, error) {

	args := u.Called(params)
	return args.Get(0).(model.User), args.Error(1)
}

func (u *UserRepositoryMock) List(params *repository.Params) ([]model.User, error) {
	args := u.Called(params)
	return args.Get(0).([]model.User), args.Error(1)
}

func (u *UserRepositoryMock) Get(params *repository.Params) (model.User, error) {
	args := u.Called(params)
	return args.Get(0).(model.User), args.Error(1)
}

func (u *UserRepositoryMock) GetByNamePassword(params *repository.Params) (model.User, error) {
	args := u.Called(params)
	return args.Get(0).(model.User), args.Error(1)
}

type UserUseCaseTestSuite struct {
	suite.Suite
	userRepoMock *UserRepositoryMock
	mockDB       *sql.DB
	mockSql      sqlmock.Sqlmock
}

func (suite *UserUseCaseTestSuite) SetupTest() {
	suite.userRepoMock = new(UserRepositoryMock)
	db, mock, _ := sqlmock.New()
	suite.mockDB = db
	suite.mockSql = mock

}

var expectedParams = &repository.Params{
	Req: repository.CreateRequest,
	User: model.User{
		Id:       "1",
		Name:     "Maher",
		Password: "1234",
		Role:     "admin",
	},
}

var expectedUser = model.User{
	Id:       "1",
	Name:     "Maher",
	Password: "1234",
	Role:     "admin",
}

func (suite *UserUseCaseTestSuite) TestCreateUser_Success() {

	// suite.userRepoMock.On("Create", expectedParams).Return(expectedUser, nil)

	suite.mockSql.ExpectQuery("INSERT INTO mst_user").WithArgs(
		expectedUser.Name,
		expectedUser.Password,
		expectedUser.Role,
	).WillReturnRows(sqlmock.NewRows([]string{"Id"}).AddRow(expectedUser.Id))

	repository.DB = suite.mockDB

	actualData, err := CreateUser(expectedUser)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), expectedUser, actualData)

}

func TestUserUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserUseCaseTestSuite))
}

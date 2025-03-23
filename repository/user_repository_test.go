package repository

import (
	"database/sql"
	"funny-login/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	mockDB     *sql.DB
	mockSql    sqlmock.Sqlmock
	mockParams *Params
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	db, mock, _ := sqlmock.New()
	suite.mockDB = db
	suite.mockSql = mock
}

var expectedUser = model.User{
	Id:       "1",
	Name:     "Maher",
	Password: "1234",
	Role:     "admin",
}

var expectedCrud = CRUD{
	Create:            expectedUser,
	List:              nil,
	Get:               model.User{},
	GetByNamePassword: model.User{},
}

func (suite *UserRepositoryTestSuite) TestCreate_Success() {
	suite.mockSql.ExpectQuery("INSERT INTO mst_user").WithArgs(
		expectedUser.Name,
		expectedUser.Password,
		expectedUser.Role,
	).WillReturnRows(sqlmock.NewRows([]string{"Id"}).AddRow(expectedUser.Id))
	DB = suite.mockDB
	suite.mockParams = &Params{
		Req:  CreateRequest,
		User: expectedUser,
	}
	actualData, err := User(suite.mockParams)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), &expectedCrud, actualData)
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

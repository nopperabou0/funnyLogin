package repository

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"enigmacamp.com/unit-test-starter-pack/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	mockDB  *sql.DB
	mockSql sqlmock.Sqlmock
	repo    UserRepository
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	db, mock, _ := sqlmock.New()
	suite.mockDB = db
	suite.mockSql = mock
	suite.repo = NewUserRepository(suite.mockDB)
}

var expectedUser = model.UserCredential{
	Id:       "1",
	Username: "Maher",
	Password: "1234",
	Role:     "admin",
}

var expectedUsers = []model.UserCredential{{
	Id:       "1",
	Username: "Maher",
	Role:     "admin",
}, {
	Id:       "2",
	Username: "Fauzi",
	Role:     "user",
},
}

func (suite *UserRepositoryTestSuite) TestCreate_Success() {
	suite.mockSql.ExpectQuery("INSERT INTO mst_user").WithArgs(
		expectedUser.Username,
		expectedUser.Password,
		expectedUser.Role,
	).WillReturnRows(sqlmock.NewRows([]string{"Id"}).AddRow(expectedUser.Id))

	actualData, err := suite.repo.Create(expectedUser)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), expectedUser, actualData)
}

func (suite *UserRepositoryTestSuite) TestCreate_Failed() {
	suite.mockSql.ExpectQuery("INSERT INTO mst_user").WithArgs(
		model.UserCredential{},
	).WillReturnError(fmt.Errorf("failed to create user"))

	actualData, err := suite.repo.Create(expectedUser)

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), model.UserCredential{}, actualData)

}

func (suite *UserRepositoryTestSuite) TestGet_Success() {
	id, _ := strconv.Atoi(expectedUser.Id)
	suite.mockSql.ExpectQuery(regexp.QuoteMeta("SELECT id, username, role FROM mst_user WHERE id = $1")).WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"id", "username", "role"}).AddRow(expectedUser.Id, expectedUser.Username, expectedUser.Role))

	actualData, err := suite.repo.Get(uint32(id))
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), model.UserCredential{
		Id:       expectedUser.Id,
		Username: expectedUser.Username,
		Role:     expectedUser.Role,
	}, actualData)
}

func (suite *UserRepositoryTestSuite) TestGet_Failed() {
	suite.mockSql.ExpectQuery(regexp.QuoteMeta(`SELECT id, username, role FROM mst_user WHERE id = $1`)).WithArgs(2).WillReturnError(fmt.Errorf("no user found"))

	actualData, err := suite.repo.Get(2)
	assert.NotNil(suite.T(), err)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), model.UserCredential{}, actualData)

}

func (suite *UserRepositoryTestSuite) TestList_Success() {
	suite.mockSql.ExpectQuery("SELECT id, username, role FROM mst_user").WithoutArgs().WillReturnRows(sqlmock.NewRows([]string{"id", "username", "role"}).AddRow(expectedUsers[0].Id, expectedUsers[0].Username, expectedUsers[0].Role).AddRow(expectedUsers[1].Id, expectedUsers[1].Username, expectedUsers[1].Role))

	actualData, err := suite.repo.List()
	assert.Nil(suite.T(), err)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedUsers, actualData)
}

func (suite *UserRepositoryTestSuite) TestList_Failed() {
	suite.mockSql.ExpectQuery("SELECT id, username, role FROM mst_user").WithoutArgs().WillReturnError(fmt.Errorf("failed to list user"))

	actualData, err := suite.repo.List()
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), actualData)
	assert.Error(suite.T(), err)
}

func (suite *UserRepositoryTestSuite) TestList_ScanFailed() {
	suite.mockSql.ExpectQuery("SELECT id, username, role FROM mst_user").WithoutArgs().WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(expectedUsers[0].Id, expectedUsers[0].Username).AddRow(expectedUsers[1].Id, expectedUsers[1].Username))

	_, err := suite.repo.List()
	assert.NotNil(suite.T(), err)
	assert.Error(suite.T(), err)

}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

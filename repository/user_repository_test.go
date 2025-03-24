package repository

import (
	"database/sql"
	"fmt"
	"funny-login/model"
	"regexp"
	"strconv"
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

var expectedUsers = []model.User{{
	Id:       "1",
	Name:     "Maher",
	Password: "",
	Role:     "admin",
}, {
	Id:       "2",
	Name:     "Fauzi",
	Password: "",
	Role:     "user",
},
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

func (suite *UserRepositoryTestSuite) TestCreate_Failed() {
	suite.mockSql.ExpectQuery("INSERT INTO mst_user").WithArgs(
		model.User{},
	).WillReturnError(fmt.Errorf("failed to create user"))
	DB = suite.mockDB
	suite.mockParams = &Params{
		Req:  CreateRequest,
		User: model.User{},
	}
	actualData, err := User(suite.mockParams)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), &CRUD{}, actualData)

}
func (suite *UserRepositoryTestSuite) TestGet_Success() {
	id, _ := strconv.Atoi(expectedUser.Id)
	suite.mockSql.ExpectQuery(regexp.QuoteMeta("SELECT id, username, role FROM mst_user WHERE id = $1")).WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"id", "username", "role"}).AddRow(expectedUser.Id, expectedUser.Name, expectedUser.Role))

	DB = suite.mockDB

	suite.mockParams = &Params{
		Req: GetRequest,
		Id:  1,
	}

	expectedCrud = CRUD{
		Create: model.User{},
		List:   nil,
		Get: model.User{
			Id:   expectedUser.Id,
			Name: expectedUser.Name,
			Role: expectedUser.Role,
		},
		GetByNamePassword: model.User{},
	}

	actualData, err := User(suite.mockParams)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), &expectedCrud, actualData)
}

func (suite *UserRepositoryTestSuite) TestGet_Failed() {
	suite.mockSql.ExpectQuery(regexp.QuoteMeta(`SELECT id, username, role FROM mst_user WHERE id = $1`)).WithArgs(2).WillReturnError(fmt.Errorf("no user found"))
	DB = suite.mockDB
	suite.mockParams = &Params{
		Req: GetRequest,
		Id:  2,
	}
	actualData, err := User(suite.mockParams)
	assert.NotNil(suite.T(), err)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), &CRUD{}, actualData)

}

func (suite *UserRepositoryTestSuite) TestList_Success() {
	suite.mockSql.ExpectQuery("SELECT id, username, role FROM mst_user").WithoutArgs().WillReturnRows(sqlmock.NewRows([]string{"id", "username", "role"}).AddRow(expectedUsers[0].Id, expectedUsers[0].Name, expectedUsers[0].Role).AddRow(expectedUsers[1].Id, expectedUsers[1].Name, expectedUsers[1].Role))

	DB = suite.mockDB
	suite.mockParams = &Params{
		Req: ListRequest,
	}

	expectedCrud = CRUD{
		Create:            model.User{},
		Get:               model.User{},
		List:              expectedUsers,
		GetByNamePassword: model.User{},
	}
	actualData, err := User(suite.mockParams)
	assert.Nil(suite.T(), err)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), &expectedCrud, actualData)
}

func (suite *UserRepositoryTestSuite) TestList_Failed() {
	suite.mockSql.ExpectQuery("SELECT id, username, role FROM mst_user").WithoutArgs().WillReturnError(fmt.Errorf("failed to list user"))

	DB = suite.mockDB

	suite.mockParams = &Params{
		Req: ListRequest,
	}

	expectedCrud = CRUD{
		Create:            model.User{},
		Get:               model.User{},
		List:              []model.User{},
		GetByNamePassword: model.User{},
	}

	actualData, err := User(suite.mockParams)
	assert.NotNil(suite.T(), err)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), &expectedCrud, actualData)
}

func (suite *UserRepositoryTestSuite) TestList_ScanFailed() {
	suite.mockSql.ExpectQuery("SELECT id, username, role FROM mst_user").WithoutArgs().WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(expectedUsers[0].Id, expectedUsers[0].Name).AddRow(expectedUsers[1].Id, expectedUsers[1].Name))

	suite.mockParams = &Params{
		Req: ListRequest,
	}

	DB = suite.mockDB

	_, err := User(suite.mockParams)
	assert.NotNil(suite.T(), err)
	assert.Error(suite.T(), err)

}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

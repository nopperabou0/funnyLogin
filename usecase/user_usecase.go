package usecase

import (
	"database/sql"
	"funny-login/model"
	"funny-login/repository"
)

func CreateUser(db *sql.DB, user model.User) model.User {
	params := &repository.Params{
		Req:  repository.CreateRequest,
		DB:   db,
		User: user,
	}
	result := repository.User(params)
	return result.Create
}

func ListAllUsers(db *sql.DB) []model.User {
	params := &repository.Params{
		Req: repository.ListRequest,
		DB:  db,
	}
	result := repository.User(params)
	return result.List
}

func GetUserById(id uint32, db *sql.DB) model.User {
	params := &repository.Params{
		Req: repository.GetRequest,
		DB:  db,
		Id:  id,
	}
	result := repository.User(params)
	return result.Get
}

func GetUserByNamePassword(name string, password string, db *sql.DB) model.User {
	params := &repository.Params{
		Req:      repository.GetByNamePasswordRequest,
		DB:       db,
		Name:     name,
		Password: password,
	}
	result := repository.User(params)
	return result.GetByNamePassword
}

package usecase

import (
	"funny-login/model"
	"funny-login/repository"
)

func CreateUser(user model.User) (model.User, error) {
	params := &repository.Params{
		Req:  repository.CreateRequest,
		User: user,
	}
	result, err := repository.User(params)
	return result.Create, err
}

func ListAllUsers() ([]model.User, error) {
	params := &repository.Params{
		Req: repository.ListRequest,
	}
	result, err := repository.User(params)
	return result.List, err
}

func GetUserById(id uint32) (model.User, error) {
	params := &repository.Params{
		Req: repository.GetRequest,
		Id:  id,
	}
	result, err := repository.User(params)
	return result.Get, err
}

func GetUserByNamePassword(name string, password string) (model.User, error) {
	params := &repository.Params{
		Req:      repository.GetByNamePasswordRequest,
		Name:     name,
		Password: password,
	}
	result, err := repository.User(params)
	return result.GetByNamePassword, err
}

package usecase

import (
	"enigmacamp.com/unit-test-starter-pack/model"
	"enigmacamp.com/unit-test-starter-pack/repository"
)

type UserUseCase interface {
	RegisterNewUser(payload model.UserCredential) (model.UserCredential, error)
	FindAllUser() ([]model.UserCredential, error)
	FindUserById(id uint32) (model.UserCredential, error)
	FindUserByUsernamePassword(username string, password string) (model.UserCredential, error)
}

type userUseCase struct {
	repo repository.UserRepository
}

func (u *userUseCase) RegisterNewUser(payload model.UserCredential) (model.UserCredential, error) {
	return u.repo.Create(payload)
}

func (u *userUseCase) FindAllUser() ([]model.UserCredential, error) {
	return u.repo.List()
}

func (u *userUseCase) FindUserById(id uint32) (model.UserCredential, error) {
	return u.repo.Get(id)
}

func (u *userUseCase) FindUserByUsernamePassword(username string, password string) (model.UserCredential, error) {
	return u.repo.GetByUsernamePassword(username, password)
}

func NewUserUseCase(repo repository.UserRepository) UserUseCase {
	return &userUseCase{repo: repo}
}

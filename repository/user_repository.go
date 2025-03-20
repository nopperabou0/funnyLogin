package repository

import (
	"database/sql"
	"fmt"
	"funny-login/model"
)

type Request string

const (
	CreateRequest            Request = "CreateUser"
	ListRequest              Request = "ListAllUsers"
	GetRequest               Request = "GetUserById"
	GetByNamePasswordRequest Request = "GetUserByNamePassword"
)

type Params struct {
	Req      Request
	User     model.User
	Id       uint32
	Name     string
	Password string
}

type CRUD struct {
	Create            model.User
	List              []model.User
	Get               model.User
	GetByNamePassword model.User
}

func User(withParameter *Params) (*CRUD, error) {
	var crud = CRUD{}
	var err error
	switch withParameter.Req {
	case CreateRequest:
		crud.Create, err = create(withParameter.User)
	case ListRequest:
		crud.List, err = list()
	case GetRequest:
		crud.Get, err = get(withParameter.Id)
	case GetByNamePasswordRequest:
		crud.GetByNamePassword, err = getByNamePassword(withParameter.Name, withParameter.Password)
	}
	return &crud, err
}

var DB *sql.DB

func CloseDB() {
	err := DB.Close()
	if err != nil {
		fmt.Println("Failed to close DB : ", err.Error())
	}
}

func create(user model.User) (model.User, error) {
	err := DB.QueryRow("INSERT INTO mst_user (username, password, role) VALUES  ($1, $2, $3) RETURNING id",
		user.Name, user.Password, user.Role,
	).Scan(&user.Id)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to create user : " + err.Error())
	}
	return user, nil
}

func list() ([]model.User, error) {
	var users []model.User
	rows, err := DB.Query("SELECT id, username, role FROM mst_user")
	if err != nil {
		return []model.User{}, fmt.Errorf("failed to list user from database : " + err.Error())
	}
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.Id, &user.Name, &user.Role)
		if err != nil {
			return []model.User{}, fmt.Errorf("failure occured when scanning data : " + err.Error())
		}
		users = append(users, user)
	}
	return users, nil
}

func get(id uint32) (model.User, error) {
	var user model.User
	err := DB.QueryRow("SELECT id, username, role FROM mst_user WHERE id = $1", id).Scan(&user.Id, &user.Name, &user.Role)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to get user by id : " + err.Error())
	}
	return user, nil
}

func getByNamePassword(name string, password string) (model.User, error) {
	var user model.User
	err := DB.QueryRow("SELECT id, username, role FROM mst_user WHERE username = $1 and password = $2", name, password).Scan(&user.Id, &user.Name, &user.Role)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to gert user by name and password : " + err.Error())
	}
	return user, nil
}

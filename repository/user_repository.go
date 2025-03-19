package repository

import (
	"database/sql"
	"fmt"
	"funny-login/config"
	"funny-login/model"
)

type Params struct {
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

var db *sql.DB

func init() {
	var connect config.Config
	db = connect.DB()
	if err := db.Ping(); err != nil {
		panic(err)
	}
}

func User(withParameter *Params) *CRUD {

	return &CRUD{
		Create:            create(withParameter.User),
		List:              list(),
		Get:               get(withParameter.Id),
		GetByNamePassword: getByNamePassword(withParameter.Name, withParameter.Password),
	}

}

func Close() {
	err := db.Close()
	if err != nil {
		fmt.Println("Failed to close DB : ", err.Error())
	}
}

func create(user model.User) model.User {
	err := db.QueryRow("INSERT INTO mst_user (username, password, role) VALUES  ($1, $2, $3) RETURNING id",
		user.Name, user.Password, user.Role,
	).Scan(&user.Id)
	if err != nil {
		fmt.Println("Failed to create user : ", err.Error())
		return model.User{}
	}
	return user
}

func list() []model.User {
	var users []model.User
	rows, err := db.Query("SELECT id, username, role FROM mst_user")
	if err != nil {
		fmt.Println("Failed to list users : ", err.Error())
		return []model.User{}
	}
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.Id, &user.Name, &user.Role)
		if err != nil {
			fmt.Println("Failure occured when scanning data : ", err.Error())
			return []model.User{}
		}
		users = append(users, user)
	}
	return users
}

func get(id uint32) model.User {
	var user model.User
	err := db.QueryRow("SELECT id, username, role FROM mst_user WHERE id = $1", id).Scan(&user.Id, &user.Name, &user.Role)
	if err != nil {
		fmt.Println("Failed to get user by id : ", err.Error())
		return model.User{}
	}
	return user
}

func getByNamePassword(name string, password string) model.User {
	var user model.User
	err := db.QueryRow("SELECT id, username, role FROM mst_user WHERE username = $1 and password = $2", name, password).Scan(&user.Id, &user.Name, &user.Role)
	if err != nil {
		fmt.Println("Failed to get user by name and password : ", err.Error())
		return model.User{}
	}
	return user
}

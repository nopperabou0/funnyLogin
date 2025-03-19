package repository

import (
	"database/sql"
	"fmt"
	"funny-login/config"
	"funny-login/model"
)

type Params struct {
	user     model.User
	id       uint32
	name     string
	password string
}

type CRUD struct {
	Create            model.User
	List              []model.User
	Get               model.User
	GetByNamePassword model.User
}

func UserRepository(withParameter *Params) *CRUD {
	var connect config.Config

	return &CRUD{
		Create:            create(withParameter.user, connect.DB()),
		List:              list(connect.DB()),
		Get:               get(connect.DB(), withParameter.id),
		GetByNamePassword: getByNamePassword(connect.DB(), withParameter.name, withParameter.password),
	}

}

func create(user model.User, db *sql.DB) model.User {
	defer db.Close()
	err := db.QueryRow("INSERT INTO mst_user (username, password, role) VALUES  ($1, $2, $3) RETURNING id",
		user.Name, user.Password, user.Role,
	).Scan(&user.Id)
	if err != nil {
		fmt.Println("Failed to create user : ", err.Error())
		return model.User{}
	}
	return user
}

func list(db *sql.DB) []model.User {
	defer db.Close()
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

func get(db *sql.DB, id uint32) model.User {
	defer db.Close()
	var user model.User
	err := db.QueryRow("SELECT id, username, role FROM mst_user WHERE id = $1", id).Scan(&user.Id, &user.Name, &user.Role)
	if err != nil {
		fmt.Println("Failed to get user by id : ", err.Error())
		return model.User{}
	}
	return user
}

func getByNamePassword(db *sql.DB, name string, password string) model.User {
	defer db.Close()
	var user model.User
	err := db.QueryRow("SELECT id, username, role FROM mst_user WHERE username = $1 and password = $2", name, password).Scan(&user.Id, &user.Name, &user.Role)
	if err != nil {
		fmt.Println("Failed to get user by name and password : ", err.Error())
		return model.User{}
	}
	return user
}

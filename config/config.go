package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type dbconfig struct {
	host     string
	port     string
	database string
	username string
	password string
	driver   string
}

type apiconfig struct {
	apiport string
}

type Config struct {
	dbconfig
	apiconfig
}

func (c *Config) DB() *sql.DB {
	if c.dbconfig == (dbconfig{}) {
		c.dbconfig = dbconfig{
			host:     "localhost",
			port:     "5432",
			database: "enigma_laundry",
			username: "postgres",
			password: "1234",
			driver:   "postgres",
		}
	}
	var dbConf string = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.host,
		c.port,
		c.username,
		c.password,
		c.database)

	db, err := sql.Open(c.driver, dbConf)

	if err != nil {
		panic(fmt.Sprintln("Connection Error : ", err.Error()))
	}

	return db
}

func (c *Config) API() string {
	if c.apiconfig == (apiconfig{}) {
		c.apiconfig = apiconfig{
			apiport: "8888",
		}
	}
	return c.apiport
}

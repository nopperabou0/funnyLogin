package config

import (
	"fmt"
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

func (c *Config) DB() string {
	if c.dbconfig == (dbconfig{}) {
		c.dbconfig = dbconfig{
			host:     "localhost",
			port:     "5432",
			database: "funnylogin",
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

	return dbConf
}

func (c *Config) API() string {
	if c.apiconfig == (apiconfig{}) {
		c.apiconfig = apiconfig{
			apiport: "8888",
		}
	}
	return c.apiport
}

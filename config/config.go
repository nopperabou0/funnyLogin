package config

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

type TokenConfig struct {
	ApplicationName     string
	JwtSignatureKey     []byte
	JwtSigningMethod    *jwt.SigningMethodHMAC
	AccessTokenLifeTime time.Duration
}

type Config struct {
	dbconfig
	apiconfig
	TokenConfig
}

func (c *Config) DB() (*sql.DB, error) {
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
		return nil, err
	}

	return db, nil
}

func (c *Config) API() string {
	if c.apiconfig == (apiconfig{}) {
		c.apiconfig = apiconfig{
			apiport: ":8888",
		}
	}
	return c.apiport
}

func (c *Config) Token() *TokenConfig {
	accessTokenLifeTime := time.Duration(1) * time.Hour
	return &TokenConfig{
		ApplicationName:     "Hell",
		JwtSignatureKey:     []byte("The humor's not the same, coming from denial"),
		JwtSigningMethod:    jwt.SigningMethodHS256,
		AccessTokenLifeTime: accessTokenLifeTime,
	}
}

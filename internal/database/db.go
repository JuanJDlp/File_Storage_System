package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Database struct {
	config     *databaseConfig
	connection *sql.DB
}

type databaseConfig struct {
	host     string
	port     string
	user     string
	password string
	dbname   string
}

func NewDatabase() *Database {
	config := &databaseConfig{
		host:     os.Getenv("DATABASE_HOST"),
		port:     os.Getenv("DATABASE_PORT"),
		password: os.Getenv("DATABASE_PASSWORD"),
		user:     os.Getenv("DATABASE_USERNAME"),
		dbname:   os.Getenv("DATABASE_NAME"),
	}
	connectionString := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable", config.host, config.port, config.user, config.password, config.dbname)

	connection, err := sql.Open("postgres", connectionString)

	if err != nil {
		panic(err)
	}

	if err := connection.Ping(); err != nil {
		panic(err)
	}

	return &Database{
		config:     config,
		connection: connection,
	}
}

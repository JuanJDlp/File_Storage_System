package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

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

	err = runScript("internal/database/sql/creation_script.sql", connection)
	if err != nil {
		log.Print(err.Error())
	}

	return &Database{
		config:     config,
		connection: connection,
	}
}

func runScript(pathToString string, connection *sql.DB) error {
	data, err := os.ReadFile(pathToString)
	if err != nil {
		return err
	}

	sql := string(data)
	sqls := strings.Split(sql, ";\n")
	for _, sql := range sqls {
		_, err := connection.Exec(sql)
		if err != nil {
			return err
		}
	}
	return nil
}

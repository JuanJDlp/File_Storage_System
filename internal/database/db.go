package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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

// NewDatabase create a new instance and connection with the database, if an error occurs it will panic
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
	var connection *sql.DB
	connection, err := sql.Open("postgres", connectionString)

	if err != nil {
		panic(err)

	}

	if err := connection.Ping(); err != nil {
		time.Sleep(5 * time.Second)

		connection, err = sql.Open("postgres", connectionString)

		if err != nil {
			panic(err)
		}
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

// runScript, will run the script located in the sql folder that created the tables, if the
// tables alredy exist it will throw an error
func runScript(pathToString string, connection *sql.DB) error {
	data, err := os.ReadFile(pathToString)
	if err != nil {
		return err
	}

	sql := string(data)
	sqls := strings.Split(sql, ";")
	for _, sql := range sqls {
		connection.Exec(sql)
	}
	return nil
}

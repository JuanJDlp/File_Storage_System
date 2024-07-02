package database

import (
	"database/sql"
	"fmt"

	"github.com/JuanJDlp/File_Storage_System/internal/model"
)

type UserRepository struct {
	Database  *Database
	TableName string
}

func (ur *UserRepository) Create(user model.UserDatabase) error {
	query := fmt.Sprintf("INSERT INTO %s (email, username, password) VALUES ($1, $2, $3)", ur.TableName)
	_, err := ur.Database.connection.Exec(query, user.Email, user.Username, user.Password)
	return err
}

func (ur *UserRepository) Get(email string) (*model.UserDatabase, error) {
	var user model.UserDatabase
	query := fmt.Sprintf("SELECT * FROM %s  WHERE email = $1", ur.TableName)
	rows, err := ur.Database.connection.Query(query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(
			&user.Email,
			&user.Username,
			&user.Password,
		)

		if err != nil {
			return nil, err
		}
	} else {
		return nil, sql.ErrNoRows
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) Delete(email string) error {
	query := fmt.Sprintf("DELETE %s WHERE email = $1", ur.TableName)
	_, err := ur.Database.connection.Exec(query, email)
	if err != nil {
		return err
	}

	return nil
}

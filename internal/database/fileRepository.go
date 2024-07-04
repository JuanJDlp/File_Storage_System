package database

import (
	"fmt"

	"github.com/JuanJDlp/File_Storage_System/internal/model"
)

type FileRepository struct {
	Database  *Database
	TableName string
}

func (fr *FileRepository) Clear() error {
	_, err := fr.Database.connection.Exec("DELETE FROM files")
	return err
}

func (fr *FileRepository) Save(file model.FileDatabase) error {
	query := fmt.Sprintf("INSERT INTO %s VALUES ($1, $2, $3, $4, $5)", fr.TableName)
	_, err := fr.Database.connection.Exec(query,
		file.Hash,
		file.FileName,
		file.Size,
		file.Date_of_upload,
		file.Owner,
	)
	return err

}

func (fr *FileRepository) Delete(email, hash string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE hash = $1 AND owner = $2", fr.TableName)
	_, err := fr.Database.connection.Exec(query,
		hash,
		email,
	)
	return err
}
func (fr *FileRepository) IsUserOwner(email, hash string) bool {
	query := fmt.Sprintf("SELECT FROM %s WHERE hash = $1 AND owner = $2", fr.TableName)
	rows, err := fr.Database.connection.Query(query,
		hash,
		email,
	)

	if err != nil {
		return false
	}

	defer rows.Close()

	if rows.Next() {
		return true
	} else {
		return false
	}
}

func (fr *FileRepository) GetAllFilesFromUser(email string) (*[]model.FilesDTO, error) {
	var files []model.FilesDTO
	query := fmt.Sprintf("SELECT filename, size, date_of_upload FROM %s WHERE owner = $1", fr.TableName)
	rows, err := fr.Database.connection.Query(query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var file model.FilesDTO
		if err := rows.Scan(&file.FileName, &file.Size, &file.Date_of_upload); err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &files, nil
}

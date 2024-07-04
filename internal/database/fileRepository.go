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

func (fr *FileRepository) Delete(file model.FileDatabase) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE hash = $1 AND owner = $2", fr.TableName)
	_, err := fr.Database.connection.Exec(query,
		file.Hash,
		file.Owner,
	)
	return err
}

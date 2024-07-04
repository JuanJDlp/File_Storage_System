package database

type FileRepository struct {
	Database  *Database
	TableName string
}

func (fr *FileRepository) Clear() error {
	_, err := fr.Database.connection.Exec("DELETE FROM $1", fr.TableName)
	return err
}

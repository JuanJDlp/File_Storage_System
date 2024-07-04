package internal

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/JuanJDlp/File_Storage_System/internal/database"
	"github.com/JuanJDlp/File_Storage_System/internal/model"
)

type Storage struct {
	BlockSize     int
	DefaultFolder string
	fileRepo      *database.FileRepository
	userRepo      *database.UserRepository
}

type file struct {
	Path string
	Name string
}

// FullPath return the path with the folder that is in the path attribute plus the file name at the end
func (f *file) FullPath() string {
	return fmt.Sprintf("%s/%s", f.Path, f.Name)
}

const defautlFolderName = "storage"
const defaultBlockSize = 8

func NewStorage(blockSize int) *Storage {
	if blockSize == 0 {
		blockSize = defaultBlockSize
	}
	db := database.NewDatabase()
	return &Storage{
		BlockSize:     blockSize,
		DefaultFolder: defautlFolderName,
		fileRepo: &database.FileRepository{
			TableName: "files",
			Database:  db,
		},
		userRepo: &database.UserRepository{
			TableName: "users",
			Database:  db,
		},
	}
}

// CreatePathForFile will create a path for a given file name
// The path contains the original file name at the end and the hashed one as the folders where it is stored
func (st *Storage) CreatePathForFile(fileName string) file {
	//Get the hash of the file name
	hashedName := HashString(fileName)
	//Find how many folders you will create
	sliceLen := len(hashedName) / st.BlockSize
	//Empty array to put the diffrent folder names
	paths := make([]string, sliceLen)
	for i := 0; i < sliceLen; i++ {
		from, to := i*st.BlockSize, (i*st.BlockSize)+st.BlockSize
		paths[i] = hashedName[from:to]
	}
	return file{
		Path: strings.Join(paths, "/"),
		Name: fileName,
	}
}

// SaveFile calls writeToFile
// The io.Reader representes the content of the file
func (st *Storage) SaveFile(fileName string, r io.Reader, owner string) (int64, error) {
	if st.Exists(fileName) {
		return 0, errors.New("the file alredy exists")
	}
	size, err := st.writeToFile(fileName, r)
	if err != nil {
		return 0, err
	}
	fileStruct := model.FileDatabase{
		Hash:           HashString(fileName),
		FileName:       fileName,
		Size:           size,
		Date_of_upload: time.Now(),
		Owner:          owner,
	}
	err = st.fileRepo.Save(fileStruct)

	return size, err
}

// ReadFile generates the complete path of the file and calls readFileStram
// it returns the size of the file, the reader of the file and an error
func (st *Storage) ReadFile(fileName, email string) (int64, io.Reader, error) {
	path := st.CreatePathForFile(fileName)
	fullPath := fmt.Sprintf("%s/%s", st.DefaultFolder, path.FullPath())
	if st.fileRepo.IsUserOwner(email, HashString(fileName)) {
		return st.readFileStream(fullPath)

	} else {
		return 0, nil, errors.New("this user does not own this file")
	}
}

// readFileStream gets the file
func (st *Storage) readFileStream(path string) (int64, io.ReadCloser, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, nil, err
	}
	fileStats, err := file.Stat()
	if err != nil {
		return 0, nil, err
	}
	return fileStats.Size(), file, nil
}

// Delete will erase the file and any directories it was in that do not contain other files in it.
func (st *Storage) Delete(fileName, email string) error {
	if !st.Exists(fileName) {
		return errors.New("the file does not exist or it was not found" + fileName)
	}

	path := st.CreatePathForFile(fileName)
	filePath := filepath.Join(st.DefaultFolder, path.FullPath())

	err := st.fileRepo.Delete(email, HashString(fileName))
	if err != nil {
		return fmt.Errorf("failed to remove the file from the database: %w", err)
	}

	err = os.RemoveAll(filePath)
	if err != nil {
		return fmt.Errorf("failed to remove file or directory: %w", err)
	}
	paths := strings.Split(path.Path, "/")

	return st.deleteEmptyFolders(fmt.Sprintf("%s/%s", st.DefaultFolder, paths[0]))
}

// deleteEmptyFolders will recursively delete the nested folder for a file.
// it start from the root and since the root is not empty it keeps goig down
// once it reacher the end, since the dic is empty because we removed the file
// the recursion starts
// I had to implement it this way since we dont want to delete another files as colateral damage
func (st *Storage) deleteEmptyFolders(root string) error {
	isEmpty, err := isEmptyDir(root)
	if err != nil {
		return err
	}

	if isEmpty {
		return os.Remove(root)
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			err = st.deleteEmptyFolders(filepath.Join(root, entry.Name()))
			if err != nil {
				return err
			}
		}
	}

	// Re-check if the directory is empty after attempting to delete subdirectories.
	isEmpty, err = isEmptyDir(root)
	if err != nil {
		return err
	}

	if isEmpty && root != st.DefaultFolder {
		return os.Remove(root)
	}

	return nil
}

// isEmptyDir returns if a folder is empty
func isEmptyDir(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	list, err := f.Readdirnames(-1)
	if err != nil {
		return false, err
	}

	return len(list) == 0, nil
}

// Exists return true if the file is present or false if it isnt present
func (st *Storage) Exists(fileName string) bool {
	path := st.CreatePathForFile(fileName)
	fullPathWithRoot := fmt.Sprintf("%s/%s", st.DefaultFolder, path.FullPath())
	_, err := os.Stat(fullPathWithRoot)
	return !errors.Is(err, os.ErrNotExist)
}

// Clear deletes all the files that the storage is managing
func (st *Storage) Clear() {
	os.RemoveAll(st.DefaultFolder)
	st.fileRepo.Clear()
}

// openFileForWriting creates the necesary directories, once that is done, it creates the fiel
// It return an error if something happened while creating the directories and the file that it just created
func (st *Storage) openFileForWriting(path string) (*os.File, error) {
	file := st.CreatePathForFile(path)
	pathWithRoot := fmt.Sprintf("%s/%s", st.DefaultFolder, file.Path)
	if err := os.MkdirAll(pathWithRoot, os.ModePerm); err != nil {
		return nil, err
	}
	fullPath := fmt.Sprintf("%s/%s", st.DefaultFolder, file.FullPath())
	return os.Create(fullPath)
}

// writeToFile create the file calling openFileForWriting and return the amount of bytes coppied and an error if necesary
func (st *Storage) writeToFile(path string, r io.Reader) (int64, error) {
	file, err := st.openFileForWriting(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	return io.Copy(file, r)
}

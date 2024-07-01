package internal

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Storage struct {
	BlockSize     int
	DefaultFolder string
}

type File struct {
	Path string
	Name string
}

// FullPath return the path with the folder that is in the path attribute plus the file name at the end
func (f *File) FullPath() string {
	return fmt.Sprintf("%s/%s", f.Path, f.Name)
}

const defautlFolderName = "storage"
const defaultBlockSize = 8

func NewStorage(blockSize int) *Storage {
	if blockSize == 0 {
		blockSize = defaultBlockSize
	}
	return &Storage{
		BlockSize:     blockSize,
		DefaultFolder: defautlFolderName,
	}
}

// CreatePathForFile will create a path for a given file name
// The path contains the original file name at the end and the hashed one as the folders where it is stored
func (st *Storage) CreatePathForFile(fileName string) File {
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
	return File{
		Path: strings.Join(paths, "/"),
		Name: fileName,
	}
}

// SaveFile calls writeToFile
func (st *Storage) SaveFile(path string, r io.Reader) {
	st.writeToFile(path, r)
}

// ReadFile generates the complete path of the file and calls readFileStram
// it returns the size of the file, the reader of the file and an error
func (st *Storage) ReadFile(fileName string) (int64, io.Reader, error) {
	path := st.CreatePathForFile(fileName)
	fullPath := fmt.Sprintf("%s/%s", st.DefaultFolder, path.FullPath())
	return st.readFileStream(fullPath)
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

// Exists return true if the file is present or false if it isnt present
func (st *Storage) Exists(fileName string) bool {
	path := st.CreatePathForFile(fileName)
	fullPathWithRoot := fmt.Sprintf("%s/%s", st.DefaultFolder, path.Path)
	_, err := os.Stat(fullPathWithRoot)
	return !errors.Is(err, os.ErrNotExist)
}

// Clear deletes all the files that the storage is managing
func (st *Storage) Clear() {
	os.RemoveAll(st.DefaultFolder)
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

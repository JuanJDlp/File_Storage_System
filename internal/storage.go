package internal

import "strings"

type Storage struct {
	BlockSize     int
	DefaultFolder string
}

type File struct {
	Path string
	Name string
}

func NewStorage(blockSize int) *Storage {
	if blockSize == 0 {
		blockSize = 8
	}
	return &Storage{
		BlockSize:     blockSize,
		DefaultFolder: "storage",
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
	paths := make([]string, sliceLen+1)
	for i := 0; i < sliceLen; i++ {
		from, to := i*st.BlockSize, (i*st.BlockSize)*st.BlockSize
		paths[i] = hashedName[from:to]
	}
	paths[len(paths)-1] = fileName
	return File{
		Path: strings.Join(paths, "/"),
		Name: fileName,
	}
}

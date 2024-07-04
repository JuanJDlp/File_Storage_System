package test

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/JuanJDlp/File_Storage_System/internal"
	"github.com/JuanJDlp/File_Storage_System/internal/database"
	"github.com/joho/godotenv"
)

func TestGeneratePathFromFileName(t *testing.T) {
	err := godotenv.Load("../.env")
	db := database.NewDatabase()
	if err != nil {
		panic(err)
	}
	stg := internal.NewStorage(0, db)
	fileName := "Final Exam.pdf"
	file := stg.CreatePathForFile(fileName)
	expectedPath := "2fdbc3b3/fc1373f6/90d09213/4c6e764d/16a7c4d6/6b1216da/bd4681df/eb3ee96a/" + fileName
	if file.FullPath() != expectedPath {
		t.Errorf("have %s want %s", file.FullPath(), expectedPath)
	}
}

func TestSaveFile(t *testing.T) {
	err := godotenv.Load("../.env")
	db := database.NewDatabase()
	if err != nil {
		panic(err)
	}
	stg := internal.NewStorage(0, db)
	contentString := "Hello, Reader!"
	fileName := "Final Exam.txt"
	path := stg.CreatePathForFile(fileName)
	content := strings.NewReader(contentString)
	stg.SaveFile(fileName, content, "j@gmail.com")
	file, err := os.Open(fmt.Sprintf("%s/%s", stg.DefaultFolder, path.FullPath()))
	if err != nil {
		t.Errorf(err.Error())
	}
	contentInTheFile, err := io.ReadAll(file)
	if err != nil {
		t.Errorf(err.Error())
	}
	if string(contentInTheFile) != contentString {
		t.Errorf("Got %s , expected %s", string(contentInTheFile), contentString)
	}
	stg.Clear()
}

func TestReadFile(t *testing.T) {
	err := godotenv.Load("../.env")
	db := database.NewDatabase()
	if err != nil {
		panic(err)
	}
	stg := internal.NewStorage(0, db)
	contentString := "Hello, Reader!"
	fileName := "Final Exam.txt"
	content := strings.NewReader(contentString)
	stg.SaveFile(fileName, content, "j@gmail.com")
	size, file, err := stg.ReadFile(fileName, "j@gmail.com")
	if err != nil {
		t.Error(err.Error())
	}
	if size == 0 {
		t.Error("The file is empty")
	}
	words, err := io.ReadAll(file)
	if string(words) != contentString {
		t.Error("The content does not match")
	}
	if err != nil {
		t.Error(err.Error())
	}
	stg.Clear()
}

func TestDelete(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}
	db := database.NewDatabase()
	stg := internal.NewStorage(0, db)
	contentString := "Hello, Reader!"
	fileName := "Final Exam.txt"
	content := strings.NewReader(contentString)
	stg.SaveFile(fileName, content, "j@gmail.com")
	err = stg.Delete(fileName, "j@gmail.com")
	if err != nil {
		t.Error(err.Error())
	}
	if stg.Exists(fileName) {
		t.Error("the file was not removed")
	}
	stg.Clear()
}

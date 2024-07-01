package test

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/JuanJDlp/File_Storage_System/internal"
)

func TestGeneratePathFromFileName(t *testing.T) {
	stg := internal.NewStorage(0)
	fileName := "Final Exam.pdf"
	file := stg.CreatePathForFile(fileName)
	expectedPath := "2fdbc3b3/fc1373f6/90d09213/4c6e764d/16a7c4d6/6b1216da/bd4681df/eb3ee96a/" + fileName
	if file.FullPath() != expectedPath {
		t.Errorf("have %s want %s", file.FullPath(), expectedPath)
	}
}

func TestSaveFile(t *testing.T) {
	stg := internal.NewStorage(0)
	contentString := "Hello, Reader!"
	fileName := "Final Exam.txt"
	path := stg.CreatePathForFile(fileName)
	content := strings.NewReader(contentString)
	stg.SaveFile(fileName, content)
	file, err := os.Open(fmt.Sprintf("%s/%s", stg.DefaultFolder, path.FullPath()))
	if err != nil {
		t.Errorf(err.Error())
	}
	fileStat, err := file.Stat()
	if err != nil {
		t.Errorf(err.Error())
	}
	contentInTheFile := make([]byte, fileStat.Size())
	for {
		n, err := file.Read(contentInTheFile)
		if err != nil && err != io.EOF {
			t.Error(err.Error())
		}
		if n == 0 {
			break //End of the file reached
		}
	}
	if string(contentInTheFile) != contentString {
		t.Errorf("Got %s , expected %s", string(contentInTheFile), contentString)
	}
	stg.Clear()
}

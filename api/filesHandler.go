package api

import (
	"net/http"

	"github.com/JuanJDlp/File_Storage_System/internal"
	"github.com/labstack/echo/v4"
)

type FilesHanlder struct {
	e       *echo.Echo
	storage *internal.Storage
}

func (fh *FilesHanlder) Start() {
	fh.e.POST("/api/v1/files", fh.saveFile)
}

func (fh *FilesHanlder) saveFile(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	src, err := file.Open()

	fh.storage.SaveFile(file.Filename, src)

	fh.e.Logger.Print(file)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, "ok")
}

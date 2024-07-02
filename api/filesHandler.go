package api

import (
	"net/http"
	"strconv"

	"github.com/JuanJDlp/File_Storage_System/internal"
	"github.com/labstack/echo/v4"
)

type FilesHanlder struct {
	e       *echo.Echo
	storage *internal.Storage
}

func (fh *FilesHanlder) Start() {
	fh.e.POST("/api/v1/files", fh.saveFile)
	fh.e.DELETE("/api/v1/files", fh.deleteFile)
	fh.e.DELETE("/api/v1/files/:name", fh.deleteOneFile)
	fh.e.GET("/api/v1/files/:name", fh.dowloadFile)
}

func (fh *FilesHanlder) saveFile(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]
	for _, file := range files {

		src, err := file.Open()

		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		defer src.Close()

		fh.storage.SaveFile(file.Filename, src)
	}

	return c.JSON(http.StatusOK, "ok")
}

func (fh *FilesHanlder) deleteFile(c echo.Context) error {
	params := struct {
		Files []string `json:"files"`
	}{}
	c.Bind(&params)
	for _, fileName := range params.Files {
		err := fh.storage.Delete(fileName)
		if err != nil {
			fh.e.Logger.Print(err.Error())
			return c.JSON(http.StatusBadRequest, err.Error())
		}

	}

	return c.JSON(http.StatusOK, "ok")
}

func (fh *FilesHanlder) deleteOneFile(c echo.Context) error {
	fileName := c.Param("name")
	err := fh.storage.Delete(fileName)
	if err != nil {
		fh.e.Logger.Print(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, "ok")
}
func (fh *FilesHanlder) dowloadFile(c echo.Context) error {
	fileName := c.Param("name")
	size, file, err := fh.storage.ReadFile(fileName)
	sizeString := strconv.Itoa(int(size))
	if err != nil {
		fh.e.Logger.Print(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	// Set the appropriate headers
	c.Response().Header().Set(echo.HeaderContentDisposition, "attachment; filename="+fileName)
	c.Response().Header().Set(echo.HeaderContentType, "application/octet-stream")
	c.Response().Header().Set(echo.HeaderContentLength, sizeString)

	// Stream the file content to the response
	return c.Stream(http.StatusOK, "application/octet-stream", file)
}

func (fh *FilesHanlder) Clear() {
	fh.storage.Clear()
}

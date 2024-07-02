package api

import (
	"mime/multipart"
	"net/http"
	"strconv"
	"sync"

	"github.com/JuanJDlp/File_Storage_System/internal"
	"github.com/labstack/echo/v4"
)

type FilesHandler struct {
	e       *echo.Echo
	storage *internal.Storage
}

func (fh *FilesHandler) Start() {
	fh.e.POST("/api/v1/files", fh.saveFile)
	fh.e.DELETE("/api/v1/files", fh.deleteFile)
	fh.e.DELETE("/api/v1/files/:name", fh.deleteOneFile)
	fh.e.GET("/api/v1/files/:name", fh.dowloadFile)
}

func (fh *FilesHandler) saveFile(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]
	var wg sync.WaitGroup
	errChan := make(chan error, len(files))
	for _, file := range files {
		wg.Add(1)

		go func(file *multipart.FileHeader) {

			defer wg.Done()

			src, err := file.Open()

			if err != nil {
				errChan <- err
				return
			}
			defer src.Close()

			_, err = fh.storage.SaveFile(file.Filename, src)
			if err != nil {
				errChan <- err
				return
			}
		}(file)

	}

	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
	}

	return c.JSON(http.StatusOK, "ok")
}

func (fh *FilesHandler) deleteFile(c echo.Context) error {
	var wg sync.WaitGroup
	params := struct {
		Files []string `json:"files"`
	}{}
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	errChan := make(chan error, len(params.Files))
	for _, fileName := range params.Files {
		wg.Add(1)
		go func(fileName string) {
			defer wg.Done()
			if err := fh.storage.Delete(fileName); err != nil {
				errChan <- err
			}
		}(fileName)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			fh.e.Logger.Print(err.Error())
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, "ok")
}

func (fh *FilesHandler) deleteOneFile(c echo.Context) error {
	fileName := c.Param("name")
	err := fh.storage.Delete(fileName)
	if err != nil {
		fh.e.Logger.Print(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, "ok")
}
func (fh *FilesHandler) dowloadFile(c echo.Context) error {
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

func (fh *FilesHandler) Clear() {
	fh.storage.Clear()
}

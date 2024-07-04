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
	e       *echo.Group
	storage *internal.Storage
}

// NewFileHandler create an instance of the fileHandler
func NewFileHandler(e *echo.Group, storage *internal.Storage) *FilesHandler {
	return &FilesHandler{
		e:       e,
		storage: storage,
	}
}

// Start, start all the routing
func (fh *FilesHandler) Start() {
	fh.e.POST("", fh.saveFile)
	fh.e.DELETE("", fh.deleteFile)
	fh.e.DELETE("/:name", fh.deleteOneFile)
	fh.e.GET("/:name", fh.dowloadFile)
}

// saveFile will take the files uploaded by the user in the multipartForm and sabe them
func (fh *FilesHandler) saveFile(c echo.Context) error {
	id := c.Request().Context().Value(internal.ContextUserKey).(string)
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

			_, err = fh.storage.SaveFile(file.Filename, src, id)
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

// deleteFile takes the name of the files to be delete and deletes them all
// it should be use when you want to delete more that 1 file at the time
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
			c.Logger().Print(err.Error())
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, "ok")
}

// deleteOneFIle takes the path parameter and deletes the file with that name
func (fh *FilesHandler) deleteOneFile(c echo.Context) error {
	fileName := c.Param("name")
	err := fh.storage.Delete(fileName)
	if err != nil {
		c.Logger().Print(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, "ok")
}

// dowloadFile takes the name of the file and dowloadsit to the user
func (fh *FilesHandler) dowloadFile(c echo.Context) error {
	fileName := c.Param("name")
	size, file, err := fh.storage.ReadFile(fileName)
	sizeString := strconv.Itoa(int(size))
	if err != nil {
		c.Logger().Print(err.Error())
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

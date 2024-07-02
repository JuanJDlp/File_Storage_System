package api

import (
	"github.com/JuanJDlp/File_Storage_System/internal"
	"github.com/JuanJDlp/File_Storage_System/internal/database"
	"github.com/labstack/echo/v4"
)

type router struct {
	e            *echo.Echo
	filesHandler *FilesHandler
}

func NewRouter() *router {
	e := echo.New()
	storage := internal.NewStorage(0)
	db := database.NewDatabase()
	filesHanlder := NewFileHandler(e, storage, db)

	return &router{
		e,
		filesHanlder,
	}

}

func (r *router) Start(port string) error {
	r.filesHandler.Start()
	return r.e.Start(":" + port)
}

func (r *router) Clear() {
	r.filesHandler.Clear()
}

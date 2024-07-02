package api

import (
	"github.com/JuanJDlp/File_Storage_System/internal"
	"github.com/labstack/echo/v4"
)

type router struct {
	e            *echo.Echo
	filesHandler *FilesHanlder
}

func NewRouter() *router {
	e := echo.New()
	storage := internal.NewStorage(0)
	filesHadler := &FilesHanlder{e, storage}

	return &router{
		e,
		filesHadler,
	}

}

func (r *router) Start(port string) error {
	r.filesHandler.Start()
	return r.e.Start(":" + port)
}

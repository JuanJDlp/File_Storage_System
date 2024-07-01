package api

import "github.com/labstack/echo/v4"

type router struct {
	e            *echo.Echo
	filesHandler *FilesHanlder
}

func NewRouter() *router {
	e := echo.New()

	return &router{
		e,
		&FilesHanlder{e},
	}
}

func (r *router) Start(port string) error {
	return r.e.Start(":" + port)
}

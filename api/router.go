package api

import (
	"github.com/JuanJDlp/File_Storage_System/internal"
	"github.com/JuanJDlp/File_Storage_System/internal/database"
	"github.com/labstack/echo/v4"
)

type router struct {
	e            *echo.Echo
	filesHandler *FilesHandler
	userHandler  *UserHandler
}

// NewRouter creates an instance of the echo router
func NewRouter() *router {
	e := echo.New()
	g := e.Group("/api/v1")

	filesEcho := g.Group("/files")
	usersEcho := g.Group("/users")
	storage := internal.NewStorage(0)
	db := database.NewDatabase()
	jwt := internal.NewJwtService()
	//Add auth to anything related with files
	filesEcho.Use(jwt.ValidateJWT)

	filesHanlder := NewFileHandler(filesEcho, storage)
	userHanlder := NewUserHandler(usersEcho,db)

	return &router{
		e:            e,
		filesHandler: filesHanlder,
		userHandler:  userHanlder,
	}

}

// Start, starts the echo server
func (r *router) Start(port string) error {
	r.filesHandler.Start()
	r.userHandler.Start()

	return r.e.Start(":" + port)
}

// Clear will delete all the files saved in the storage
func (r *router) Clear() {
	r.filesHandler.Clear()
}

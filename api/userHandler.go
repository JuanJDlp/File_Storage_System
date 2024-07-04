package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/JuanJDlp/File_Storage_System/internal"
	"github.com/JuanJDlp/File_Storage_System/internal/database"
	"github.com/JuanJDlp/File_Storage_System/internal/model"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	e              *echo.Group
	jwtService     *internal.JwtService
	userRepository *database.UserRepository
}

// NewUserHandler creater a new user hanlder with a jwtService and a userRepository
func NewUserHandler(e *echo.Group, db *database.Database) *UserHandler {
	return &UserHandler{
		e: e,
		jwtService: &internal.JwtService{
			JwtSecret: os.Getenv("JWT_SECRET"),
		},
		userRepository: &database.UserRepository{
			Database:  db,
			TableName: "users",
		},
	}
}

// Start will activate the endpoints
func (uh *UserHandler) Start() {
	uh.e.POST("/login", uh.logIn)
	uh.e.POST("/register", uh.register)
	uh.e.POST("/update", internal.NewJwtService().ValidateJWT(uh.updateUser))
}

// logIn is the functin that will be called when trying to log in
func (uh *UserHandler) logIn(c echo.Context) error {
	var user model.UserDatabase
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "there was an error binding")
	}

	userDb, err := uh.userRepository.Get(user.Email)
	if err != nil {
		log.Print(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if valid := unhashPassword(user.Password, userDb.Password); valid {
		token, err := uh.jwtService.CreateJwtTokenForUser(user.Email)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		response := struct {
			Email string
			Token string
		}{
			user.Email,
			token,
		}
		return c.JSON(http.StatusOK, response)
	} else {
		return echo.NewHTTPError(http.StatusBadRequest, "the password is incorrect")
	}
}

// register is the function that is called when the user is creating a new account
func (uh *UserHandler) register(c echo.Context) error {
	var user model.UserDatabase
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "it wasn't possible to get the user, please check your request")
	}

	password, err := hashPassword(user.Password)
	if err != nil {
		c.Logger().Print(err)
		return err
	}

	_, err = uh.userRepository.Get(user.Email)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return echo.NewHTTPError(http.StatusBadRequest, "the user alredy exists")
	}

	user.Password = password

	err = uh.userRepository.Create(user)

	if err != nil {
		c.Logger().Print(err)
		return echo.NewHTTPError(http.StatusBadRequest, "there was an error creating the user")
	}
	response := struct {
		Email   string
		created bool
	}{
		Email:   user.Email,
		created: true,
	}
	return c.JSON(http.StatusOK, response)
}

// hashPassword will create a 32 bytes hash for the password
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// unhashPassword will tell you if two passwords are the same
func unhashPassword(passwordInserted, passwordHashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHashed), []byte(passwordInserted))
	return err == nil
}

// updateUser will receive a new username and password and change it on the database
func (uh *UserHandler) updateUser(c echo.Context) error {
	var user model.UserDatabase

	id := c.Request().Context().Value(internal.ContextUserKey).(string)
	c.Bind(&user)

	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "it wasn't possible to get the user, please check your request")
	}

	password, err := hashPassword(user.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "There was an error hashing the password", err)
	}
	err = uh.userRepository.Update(id, user.Username, password)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "The user does not exist", err)
	}

	response := struct {
		NewUsername    string
		paswordChanged bool
	}{
		user.Username,
		true,
	}
	return c.JSON(http.StatusOK, response)
}

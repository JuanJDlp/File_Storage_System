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
	e              *echo.Echo
	jwtService     *internal.JwtService
	userRepository *database.UserRepository
}

func NewUserHandler(e *echo.Echo, db *database.Database) *UserHandler {
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

func (uh *UserHandler) Start() {
	uh.e.POST("/api/v1/login", uh.logIn)
	uh.e.POST("/api/v1/register", uh.register)
}

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

func (uh *UserHandler) register(c echo.Context) error {
	var user model.UserDatabase
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "it wasn't possible to get the user, please check your request")
	}

	password, err := hashPasword(user.Password)
	if err != nil {
		uh.e.Logger.Print(err)
		return err
	}

	_, err = uh.userRepository.Get(user.Email)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return echo.NewHTTPError(http.StatusBadRequest, "the user alredy exists")
	}

	user.Password = password

	err = uh.userRepository.Create(user)

	if err != nil {
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

func hashPasword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func unhashPassword(passwordInserted, passwordHashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHashed), []byte(passwordInserted))
	return err == nil
}

func (uh *UserHandler) updateUser(c echo.Context) error {
	var user model.UserDatabase

	id := c.Request().Context().Value(internal.ContextUserKey).(string)
	c.Bind(&user)

	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "it wasn't possible to get the user, please check your request")
	}

	password, err := hashPasword(user.Password)
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

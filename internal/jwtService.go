package internal

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type IdContext string

const ContextUserKey IdContext = "user"

type JwtService struct {
	JwtSecret string
}

func NewJwtService() *JwtService {
	return &JwtService{
		JwtSecret: os.Getenv("JWT_SECRET")}
}

// ValidateJWT is a middleware that will check if a token is valid or not
func (jw *JwtService) ValidateJWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := jw.GetJwtToken(c.Request().Header.Get("Authorization"))
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid JWT")

		}
		if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
			id := claims.Subject
			if id == "" {
				return echo.NewHTTPError(http.StatusBadRequest, "an id was not provided")

			}
			ctx := context.WithValue(c.Request().Context(), ContextUserKey, id)
			req := c.Request().WithContext(ctx)
			c.SetRequest(req)
		} else {
			return echo.NewHTTPError(http.StatusUnauthorized, "the token is not valid")
		}
		return next(c)
	}

}

// ExtractTokenString will get the token with an specific work, it is implemented this way so
// you can use it to stract apikeys or jwt token
func ExtractTokenString(word, header string) (string, error) {

	tokenString := strings.TrimSpace(strings.TrimPrefix(header, word))
	if tokenString == "" {
		return "", errors.New("the token is empty")
	}
	return tokenString, nil
}

// GetJWTToken will create a new token struct from a rawToken String
func (jw *JwtService) GetJwtToken(rawToken string) (*jwt.Token, error) {
	tokenString, err := ExtractTokenString("Bearer", rawToken)
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(jw.JwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

// CreateJWTTokenForUser will create a new JWT token for a given email
func (jw *JwtService) CreateJwtTokenForUser(email string) (string, error) {
	claimsMap := jwt.RegisteredClaims{
		Subject:  email,
		Issuer:   "Chirp",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsMap)
	s, err := t.SignedString([]byte(jw.JwtSecret))
	if err != nil {
		return "", err
	}
	return s, nil
}

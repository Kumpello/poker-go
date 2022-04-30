package webapi

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"pokergo/pkg/jwt"
)

type Router interface {
	Route(e *echo.Echo, prefix string) error
}

func GetJWTToken(c echo.Context) (jwt.SignedToken, error) {
	token := c.Get("user")
	jwtToken, ok := token.(jwt.SignedToken)
	if !ok {
		return jwt.SignedToken{}, errors.New("invalid jwt token")
	}

	return jwtToken, nil
}

func NewEcho(
	jwtInstance *jwt.JWT,
	authMux Router,
	orgRouter Router,
	gameRouter Router) *echo.Echo {
	e := echo.New()

	_ = authMux.Route(e, "auth")
	_ = orgRouter.Route(e, "org")
	_ = gameRouter.Route(e, "game")

	e.GET("health", func(c echo.Context) error {
		return c.JSON(200, "ok")
	})

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if uri := c.Request().RequestURI; uri == "/auth/signup" || uri == "/auth/login" {
				// auth method, don't check JWT token
				return next(c)
			}

			jwtToken := c.Request().Header.Get("token")
			if jwtToken == "" {
				return c.String(403, "missing jwt token")
			}

			v, err := jwtInstance.ValidateToken(jwtToken)
			if err != nil {
				return c.String(403, fmt.Sprintf("invalid token: %s", err.Error()))
			}
			c.Set("user", v)

			return next(c)
		}
	})

	return e
}

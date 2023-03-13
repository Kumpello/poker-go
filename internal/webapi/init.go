package webapi

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"pokergo/pkg/jwt"
	"pokergo/pkg/logger"
)

type Router interface {
	Route(g *echo.Group)
}

type echoValidator struct {
	validator *validator.Validate
}

func (v *echoValidator) Validate(i any) error {
	return v.validator.Struct(i) // nolint:wrapcheck  // that's ok (echo framework)
}

func GetJWTToken(c echo.Context) (jwt.SignedToken, error) {
	token := c.Get("user")
	jwtToken, ok := token.(jwt.SignedToken)
	if !ok {
		return jwt.SignedToken{}, errors.New("invalid jwt token")
	}

	return jwtToken, nil
}

type EchoRouters struct {
	AuthRouter Router
	OrgRouter  Router
	GameRouter Router
	NewsRouter Router
}

func NewEcho(
	validate *validator.Validate,
	jwtInstance *jwt.JWT,
	routers EchoRouters,
	log logger.Logger,
	debug bool,
) *echo.Echo {
	e := echo.New()
	e.Debug = debug
	e.Validator = &echoValidator{validator: validate}

	auth := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			jwtToken := c.Request().Header.Get("Authorization")
			if jwtToken == "" {
				return c.String(403, "missing jwt token")
			}

			if !strings.HasPrefix(jwtToken, "Bearer: ") {
				return c.String(400, "token must start with bearer:")
			}

			v, err := jwtInstance.ValidateToken(strings.TrimPrefix(jwtToken, "Bearer: "))
			if err != nil {
				return c.String(403, fmt.Sprintf("invalid token: %s", err.Error()))
			}
			c.Set("user", v)

			return next(c)
		}
	}

	authRouter := e.Group("/auth")
	orgRouter := e.Group("/org", auth)
	gameRouter := e.Group("/game", auth)
	newsRouter := e.Group("/news")

	routers.AuthRouter.Route(authRouter)
	routers.OrgRouter.Route(orgRouter)
	routers.GameRouter.Route(gameRouter)
	routers.NewsRouter.Route(newsRouter)

	e.GET("health", func(c echo.Context) error {
		return c.JSON(200, "ok")
	})

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger.MakeEchoLogEntry(log, c).Info("incoming request")
			return next(c)
		}
	})

	e.Use()

	// Required for the vue-js framework
	// better security rules
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// for local development
			c.Response().Header().Add("Access-Control-Allow-Origin", "http://localhost:3000")
			c.Response().Header().Add("Access-Control-Allow-Methods", "PUT, POST, PATCH, DELETE, GET")
			c.Response().Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

			return next(c)
		}
	})

	return e
}

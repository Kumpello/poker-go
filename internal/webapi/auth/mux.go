package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"pokergo/internal/users"
	"pokergo/pkg/crypto"
	"pokergo/pkg/id"
	"pokergo/pkg/jwt"
	"pokergo/pkg/timer"
	"time"
)

type mux struct {
	userAdapter users.Adapter
	timer       timer.Timer
	jwt         *jwt.JWT
	validator   *validator.Validate
}

func NewMux(
	userAdapter users.Adapter,
	timer timer.Timer,
	jwt *jwt.JWT,
	validator *validator.Validate,
) *mux {
	return &mux{userAdapter: userAdapter, timer: timer, jwt: jwt, validator: validator}
}

func (m *mux) Route(e *echo.Echo, prefix string) error {
	e.POST(prefix+"/signup", m.SignUp)
	e.POST(prefix+"/login", m.LogIn)
	return nil
}

func (m *mux) SignUp(c echo.Context) error {
	reqCtx, cancel := context.WithTimeout(c.Request().Context(), time.Duration(60)*time.Second)
	defer cancel()

	var request signUpRequest
	if err := c.Bind(&request); err != nil {
		return c.String(400, fmt.Sprintf("invalid request: %s", err.Error()))
	}
	if err := m.validator.Struct(request); err != nil {
		return c.String(400, fmt.Sprintf("invalid request: %s", err.Error()))
	}

	encPass, err := crypto.HashPassword(request.Password)
	if err != nil {
		return c.String(400, fmt.Sprintf("cannot encrypt password: %s", err.Error()))
	}

	u := users.User{
		ID:           id.ID{}, // stub
		Username:     request.Name,
		Password:     encPass,
		Email:        request.Email,
		Token:        "",
		RefreshToken: "",
		CreatedAt:    m.timer.Now(),
		UpdatedAt:    m.timer.Now(),
	}

	if u, err = m.userAdapter.NewUser(reqCtx, u); err != nil { // overwrite user for ID and generated data
		return c.String(500, fmt.Sprintf("cannot create user: %s", err.Error()))
	}

	token, refresh, err := m.jwt.GenerateTokens(u.Email, u.Username, u.ID)
	if err != nil {
		return c.String(500, fmt.Sprintf("cannot generate user token, "+
			"but the user was created: %s", err.Error()))
	}

	if err := m.userAdapter.UpdateTokens(reqCtx, u.ID, &token, &refresh); err != nil {
		return c.String(500, fmt.Sprintf("cannot update user token, "+
			"but user was created: %s", err.Error()))
	}

	return c.JSON(200, authResponse{
		ID:           u.ID.Hex(),
		Token:        token,
		RefreshToken: refresh,
	})
}

func (m *mux) LogIn(c echo.Context) error {
	reqCtx, cancel := context.WithTimeout(c.Request().Context(), time.Duration(60)*time.Second)
	defer cancel()

	var request logInRequest
	if err := c.Bind(&request); err != nil {
		return c.String(400, fmt.Sprintf("cannot bind input data: %s", err.Error()))
	}
	if err := m.validator.Struct(request); err != nil {
		return c.String(400, fmt.Sprintf("invalid request: %s", err.Error()))
	}

	u, err := m.userAdapter.GetUserByName(reqCtx, request.Name)
	if err != nil {
		if errors.Is(err, users.ErrUserNotExists) {
			return c.String(404, "not exists")
		}
		return c.String(500, fmt.Sprintf("cannot find user (internal error): %s", err.Error()))
	}

	if err := crypto.VerifyPassword(u.Password, request.Password); err != nil {
		return c.String(403, fmt.Sprintf("invalid password: %s", err.Error()))
	}

	token, refresh, err := m.jwt.GenerateTokens(u.Email, u.Username, u.ID)
	if err != nil {
		return c.String(500, fmt.Sprintf("cannot generate user token, "+
			"but the user was created: %s", err.Error()))
	}

	if err := m.userAdapter.UpdateTokens(reqCtx, u.ID, &token, &refresh); err != nil {
		return c.String(500, fmt.Sprintf("cannot update user token, "+
			"but user was created: %s", err.Error()))
	}

	return c.JSON(200, authResponse{
		ID:           u.ID.Hex(),
		Token:        token,
		RefreshToken: refresh,
	})
}

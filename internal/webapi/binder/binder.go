package binder

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"pokergo/internal/webapi"
	"pokergo/pkg/id"
	"pokergo/pkg/jwt"
	"time"
)

type Context struct {
	Ctx    context.Context
	Cancel context.CancelFunc
	Echo   echo.Context

	UserID    id.ID
	TokenData jwt.SignedToken
}

type StructValidator interface {
	Struct(any interface{}) error
}

type BindError struct {
	Code    int
	Message string
}

func (b BindError) Error() string {
	return fmt.Sprintf("%d: %s", b.Code, b.Message)
}

var _ error = (*BindError)(nil)

func BindRequest[T any](c echo.Context, v StructValidator) (Context, T, *BindError) {
	result := Context{
		Echo: c,
	}
	var t T

	// Obtain context and cancel
	reqCtx, cancel := context.WithTimeout(c.Request().Context(), time.Duration(60)*time.Second)
	result.Ctx = reqCtx
	result.Cancel = cancel

	// Obtain jwt token
	jwtToken, err := webapi.GetJWTToken(c)
	if err != nil {
		return result, t, &BindError{403, "jwt token invalid"}
	}
	requesterID, err := id.FromString(jwtToken.ID)
	if err != nil {
		return result, t, &BindError{400, fmt.Sprintf("invalid user id: %s", err)}
	}
	result.UserID = requesterID
	result.TokenData = jwtToken

	// Obtain request
	var request T
	if err := c.Bind(&request); err != nil {
		return result, t, &BindError{400, fmt.Sprintf("invalid request: %s", err.Error())}
	}

	if err := v.Struct(request); err != nil {
		return result, t, &BindError{400, fmt.Sprintf("invalid request: %s", err.Error())}
	}

	return result, request, nil
}

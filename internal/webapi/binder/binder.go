package binder

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/labstack/echo/v4"
	"pokergo/internal/webapi"
	"pokergo/pkg/id"
	"pokergo/pkg/jwt"
)

// BaseContext is interface over Context without generic type
// Allows to use Context without generic type
type BaseContext interface {
	Context() context.Context
	Cancel() context.CancelFunc
	Echo() echo.Context
	UserID() id.ID
	TokenData() jwt.SignedToken
}

type Context[T any, Q any] struct {
	ctx    context.Context
	cancel context.CancelFunc
	echo   echo.Context

	userID    id.ID
	tokenData jwt.SignedToken

	Request T
	Query   Q
}

func (c Context[T, Q]) Context() context.Context {
	return c.ctx
}

func (c Context[T, Q]) Cancel() context.CancelFunc {
	return c.cancel
}

func (c Context[T, Q]) Echo() echo.Context { // nolint:ireturn // nolintlint
	return c.echo
}

func (c Context[T, Q]) UserID() id.ID {
	return c.userID
}

func (c Context[T, Q]) TokenData() jwt.SignedToken {
	return c.tokenData
}

type StructValidator interface {
	Struct(str any) error
}

// BindRequest bind requests returning Context, user data (if requireAuth) and an error.
// T must be a simple type to be validated (pointers are not validated).
func BindRequest[T any, Q any](
	c echo.Context,
	requireAuth bool,
) (*Context[T, Q], *BindError) {
	result := &Context[T, Q]{
		echo: c,
	}
	var t T
	var q Q

	// Obtain context and cancel
	reqCtx, cancel := context.WithTimeout(c.Request().Context(), time.Duration(60)*time.Second)
	result.ctx = reqCtx
	result.cancel = cancel

	if requireAuth {
		jwtToken, err := webapi.GetJWTToken(c)
		if err != nil {
			return result, &BindError{403, "jwt token invalid"}
		}
		requesterID, err := id.FromString(jwtToken.ID)
		if err != nil {
			return result, &BindError{400, fmt.Sprintf("invalid user id: %s", err)}
		}
		result.userID = requesterID
		result.tokenData = jwtToken
	}

	// Obtain request
	if err := c.Bind(&t); err != nil {
		return result, &BindError{400, fmt.Sprintf("invalid request: %s", err.Error())}
	}

	if val := reflect.ValueOf(t); val.Kind() == reflect.Struct { // don't validate interface{} type
		if err := c.Validate(t); err != nil {
			return result, &BindError{400, fmt.Sprintf("invalid request: %s", err.Error())}
		}
	}

	result.Request = t
	result.Query = q
	return result, nil
}

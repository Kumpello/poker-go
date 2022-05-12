package binder

import "fmt"

type BindError struct {
	Code    int
	Message string
}

func (b BindError) Error() string {
	return fmt.Sprintf("%d: %s", b.Code, b.Message)
}

var _ error = (*BindError)(nil)

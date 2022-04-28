package game

import "errors"

var (
	ErrUserNotFound      = errors.New("user not exists")
	ErrGameNotFinished   = errors.New("some players have not their final stack set")
	ErrStackInconsistent = errors.New("the sum of final stacks is differ than the sum of buy ins")
)

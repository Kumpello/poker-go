package game

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrUserNotFound = errors.New("user not exists")
	ErrOrgNotFound  = errors.New("org not exists")

	ErrGameNotFinished   = errors.New("some players have not their final stack set")
	ErrStackInconsistent = errors.New("the sum of final stacks is differ than the sum of buy ins")

	ErrInsufficientPermissions = errors.New("insufficient permissions to manage game")

	ErrGameNotExists = mongo.ErrNoDocuments
)

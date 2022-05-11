package id

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ID = primitive.ObjectID

var ZeroID = [12]byte{} // nolint:gochecknoglobals // cannot be const

func NewID() ID {
	return primitive.NewObjectID()
}

func FromString(s string) (ID, error) {
	if s == "" {
		return ZeroID, nil
	}

	id, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		return ZeroID, fmt.Errorf("cannot parse id: %w", err)
	}
	return id, nil
}

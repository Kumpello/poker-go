package id

import "go.mongodb.org/mongo-driver/bson/primitive"

type ID = primitive.ObjectID

var ZeroID = [12]byte{}

func NewID() ID {
	return primitive.NewObjectID()
}

func FromString(s string) ID {
	var emptyID primitive.ObjectID
	id, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		return emptyID
	}

	return id
}

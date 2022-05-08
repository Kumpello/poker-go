package game

import (
	"time"

	"pokergo/pkg/id"
)

type Data struct {
	ID           id.ID     `bson:"_id"` // nolint:tagliatelle // mongo-id
	Organizer    id.ID     `bson:"organizer"`
	Organization id.ID     `bson:"organization"`
	Start        time.Time `bson:"start"`
	Players      []Player  `bson:"players"`
}

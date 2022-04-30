package game

import (
	"pokergo/pkg/id"
	"time"
)

type Data struct {
	ID           id.ID     `bson:"_id"`
	Organizer    id.ID     `bson:"organizer"`
	Organization id.ID     `bson:"organization"`
	Start        time.Time `bson:"start"`
	Players      []Player  `bson:"players"`
}

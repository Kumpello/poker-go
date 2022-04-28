package game

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"pokergo/internal/users"
	"pokergo/pkg/id"
	"pokergo/pkg/timer"
)

// Adapter allows creating and updating games (persistent part)
type Adapter interface {
	// NewGame create a new game (organizer must exist) and returns the created object
	NewGame(ctx context.Context, organizer string) (Data, error)
	// Update updates game model in database (this is quite inefficient because the whole document is replaced)
	Update(ctx context.Context, new Data) error
}

type mongoAdapter struct {
	coll mongo.Collection

	users users.Adapter
	timer timer.Timer
}

func (m *mongoAdapter) NewGame(ctx context.Context, organizer string) (Data, error) {
	u, err := m.users.GetUserByName(ctx, organizer)
	if err != nil {
		return Data{}, fmt.Errorf("cannot find organizer: %w", err)
	}

	gameData := Data{
		ID:        id.NewID(),
		Organizer: u.ID,
		Start:     m.timer.Now(),
		Players:   nil,
	}

	_, err = m.coll.InsertOne(ctx, gameData)
	if err != nil {
		return Data{}, fmt.Errorf("cannot create a new game in mongo: %w", err)
	}

	return gameData, nil
}

func (m *mongoAdapter) Update(ctx context.Context, new Data) error {
	filter := bson.M{
		"id": new.ID,
	}

	_, err := m.coll.ReplaceOne(ctx, filter, new)
	if err != nil {
		return fmt.Errorf("cannot update game: %w", err)
	}

	return nil
}

var _ Adapter = (*mongoAdapter)(nil)

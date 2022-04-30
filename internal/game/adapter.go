package game

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"pokergo/pkg/id"
	"pokergo/pkg/timer"
)

// Adapter allows creating and updating games (persistent part)
type Adapter interface {
	// NewGame create a new game (organizer must exist) and returns the created object
	NewGame(ctx context.Context, orgID, uID id.ID) (Data, error)
	// Update updates game model in database (this is quite inefficient because the whole document is replaced)
	Update(ctx context.Context, new Data) error
	// FindGameByID looks for a game by id
	FindGameByID(ctx context.Context, uID id.ID) (Data, error)
}

type mongoAdapter struct {
	coll  *mongo.Collection
	timer timer.Timer
}

func NewMongoAdapter(
	coll *mongo.Collection,
	timer timer.Timer,
) *mongoAdapter {
	return &mongoAdapter{coll: coll, timer: timer}
}

func (m *mongoAdapter) EnsureIndexes(ctx context.Context) error {
	return nil // no indexes required (except the default one)
}

func (m *mongoAdapter) NewGame(ctx context.Context, orgID, uID id.ID) (Data, error) {
	gameData := Data{
		ID:           id.NewID(),
		Organizer:    orgID,
		Organization: uID,
		Start:        m.timer.Now(),
		Players:      nil,
	}

	_, err := m.coll.InsertOne(ctx, gameData)
	if err != nil {
		return Data{}, fmt.Errorf("cannot create a new game in mongo: %w", err)
	}

	return gameData, nil
}

func (m *mongoAdapter) Update(ctx context.Context, new Data) error {
	filter := bson.M{
		"_id": new.ID,
	}

	_, err := m.coll.ReplaceOne(ctx, filter, new)
	if err != nil {
		return fmt.Errorf("cannot update game: %w", err)
	}

	return nil
}

func (m *mongoAdapter) FindGameByID(ctx context.Context, uID id.ID) (Data, error) {
	filter := bson.M{
		"_id": uID,
	}

	res := m.coll.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Data{}, ErrGameNotExists
		}
		return Data{}, fmt.Errorf("cannot perform query: %w", err)
	}

	var data Data
	if err := res.Decode(&data); err != nil {
		return Data{}, fmt.Errorf("cannot decode result data: %w", err)
	}

	return data, nil
}

var _ Adapter = (*mongoAdapter)(nil)

package game

import (
	"context"
	"fmt"
	"pokergo/internal/users"
	"pokergo/pkg/id"
	"pokergo/pkg/logger"
	"sync"
)

type Manager interface {
	// CreateGame creates an empty game
	CreateGame(ctx context.Context, organizer string) (*Game, error)
}

type manager struct {
	gameAdapter  Adapter
	usersAdapter users.Adapter

	gamesMux sync.Mutex
	games    map[id.ID]*Game
}

func (m *manager) CreateGame(ctx context.Context, organizer string) (*Game, error) {
	gameData, err := m.gameAdapter.NewGame(ctx, organizer)
	if err != nil {
		return nil, fmt.Errorf("cannot create a new game: %w", err)
	}

	game := Game{
		Data:         gameData,
		playerMux:    sync.Mutex{},
		gameLogger:   logger.NewLogger(),
		usersAdapter: m.usersAdapter,
	}

	m.gamesMux.Lock()
	defer m.gamesMux.Unlock()
	if _, exists := m.games[game.ID]; exists {
		return nil, fmt.Errorf("game exists in the cache")
	}
	m.games[game.ID] = &game

	return &game, nil
}

package game

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"pokergo/internal/org"
	"pokergo/internal/users"
	"pokergo/pkg/id"
	"pokergo/pkg/logger"
)

type Manager interface {
	// CreateGame creates an empty game
	CreateGame(ctx context.Context, uID id.ID, orgName string) (*Game, error)
	// GetGame returns a game from cache or creates a new instance (if game exists).
	// Changes made to obtained Game are not saved, must be committed via Commit
	GetGame(ctx context.Context, callerID, id id.ID) (*Game, error)
	// Commit saves changes made to game in the persistent storage
	Commit(ctx context.Context, callerID, gID id.ID) error
}

type manager struct {
	gameAdapter  Adapter
	usersAdapter users.Adapter
	orgAdapter   org.Adapter

	gamesMux sync.Mutex
	games    map[id.ID]*Game // TODO(pmaterna): This should be TTL or something
}

func NewManager(
	gameAdapter Adapter,
	usersAdapter users.Adapter,
	orgAdapter org.Adapter,
) *manager {
	return &manager{
		gameAdapter:  gameAdapter,
		usersAdapter: usersAdapter,
		orgAdapter:   orgAdapter,
		gamesMux:     sync.Mutex{},
		games:        make(map[id.ID]*Game),
	}
}

func (m *manager) CreateGame(ctx context.Context, uID id.ID, orgName string) (*Game, error) {
	o, err := m.orgAdapter.GetOrgByName(ctx, orgName)
	if err != nil {
		if errors.Is(err, org.ErrOrgNotExists) {
			return nil, ErrOrgNotFound
		}
		return nil, fmt.Errorf("cannot find org: %w", err)
	}

	gameData, err := m.gameAdapter.NewGame(ctx, uID, o.ID)
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

func (m *manager) GetGame(ctx context.Context, callerID, id id.ID) (*Game, error) {
	m.gamesMux.Lock()
	defer m.gamesMux.Unlock()

	g, ok := m.games[id]
	if ok {
		return g, nil
	}

	d, err := m.gameAdapter.FindGameByID(ctx, id)
	if err != nil {
		return nil, ErrGameNotExists
	}

	o, err := m.orgAdapter.GetOrgByID(ctx, d.Organization)
	if err != nil {
		return nil, ErrOrgNotFound
	}
	if !o.IsMember(callerID) {
		return nil, ErrInsufficientPermissions
	}

	g = &Game{
		Data:         d,
		playerMux:    sync.Mutex{},
		gameLogger:   logger.NewLogger(),
		usersAdapter: m.usersAdapter,
	}

	m.games[g.ID] = g

	return g, nil
}

func (m *manager) Commit(ctx context.Context, callerID, gID id.ID) error {
	g, err := m.GetGame(ctx, callerID, gID)
	if err != nil {
		return fmt.Errorf("cannot obtain the game: %w", err)
	}

	if err := m.gameAdapter.Update(ctx, g.Data); err != nil {
		return fmt.Errorf("cannot update the game: %w", err)
	}

	return nil
}

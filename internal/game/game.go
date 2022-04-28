package game

import (
	"context"
	"fmt"
	"pokergo/internal/users"
	"pokergo/pkg/id"
	"pokergo/pkg/logger"
	"sync"
	"time"
)

type Data struct {
	ID        id.ID     `bson:"_id" json:"id"`
	Organizer id.ID     `bson:"organizer" json:"organizer"`
	Start     time.Time `bson:"start" json:"start"`
	Players   []Player  `bson:"players" json:"players"`
}

type Game struct {
	Data

	playerMux    sync.Mutex
	gameLogger   logger.Logger
	usersAdapter users.Adapter
}

type inGameTransaction struct {
	amount int64
	reason string
	from   id.ID
}

// AppendPlayer appends a new player to the game
func (g *Game) AppendPlayer(ctx context.Context, user id.ID, startStack int64) error {
	g.playerMux.Lock()
	defer g.playerMux.Unlock()

	u, err := g.usersAdapter.GetUserByID(ctx, user)
	if err != nil {
		return fmt.Errorf("cannot find user: %w", err)
	}

	g.Players = append(g.Players, Player{
		UserID:            u.ID,
		BuyIn:             startStack,
		BuyOut:            nil,
		AdditionalIncomes: []inGameTransaction{},
	})

	g.gameLogger.Infof("adding new player to the game: %s with starting stack %d", u.Username, startStack)
	return nil
}

func (g *Game) SetFinishStack(id id.ID, stack int64) error {
	g.playerMux.Lock()
	defer g.playerMux.Unlock()

	if u, err := g.findPlayer(id); err == nil {
		u.BuyOut = &stack
		g.gameLogger.Infof("user %s finishes the game with %d stack", id.Hex(), stack)
		return nil
	}

	g.gameLogger.Errorf("user %s not exists", id.Hex())
	return ErrUserNotFound
}

// ReBuyIn is a re-buy-in from the bank (just increase starting stack)
func (g *Game) ReBuyIn(player id.ID, amount int64) error {
	g.playerMux.Lock()
	defer g.playerMux.Unlock()

	p, err := g.findPlayer(player)
	if err != nil {
		g.gameLogger.Errorf("user %s not exists", player.Hex())
		return err
	}

	g.gameLogger.Errorf("user %s re-bought for %d", player.Hex(), amount)

	p.BuyIn += amount
	return nil
}

// ReBuyInFromPlayer is a re-buy-in but paid by another player
func (g *Game) ReBuyInFromPlayer(buyer, seller id.ID, amount int64) error {
	g.playerMux.Lock()
	defer g.playerMux.Unlock()

	buyerPlayer, err := g.findPlayer(buyer)
	if err != nil {
		g.gameLogger.Errorf("user (buyer) %s not exists", buyer.Hex())
		return err
	}

	sellerPlayer, err := g.findPlayer(seller)
	if err != nil {
		g.gameLogger.Errorf("user (seller) %s not exists", seller.Hex())
		return err
	}

	// add cash to seller
	sellerPlayer.AdditionalIncomes = append(sellerPlayer.AdditionalIncomes, inGameTransaction{
		amount: amount,
		reason: "re-buy-in from another player",
		from:   buyer,
	})

	g.gameLogger.Infof("transaction done: %s --%d--> %s",
		sellerPlayer.UserID, buyerPlayer.UserID, amount)

	// and to the buyer as a start cash
	buyerPlayer.BuyIn += amount
	return nil
}

func (g *Game) Verify() error {
	g.playerMux.Lock()
	defer g.playerMux.Unlock()

	// 1. All players have finishStack set
	// 2. sum(buy-ins) == sum(take-outs)
	var buyIns int64
	var buyOuts int64
	for _, p := range g.Players {
		if p.BuyOut == nil {
			return ErrGameNotFinished
		}

		buyIns += p.BuyIn
		buyOuts += *p.BuyOut

		// add inGameIncomes from transactions too
		// buyIns were modified on transaction
		for _, a := range p.AdditionalIncomes {
			buyOuts += a.amount
		}
	}

	if buyOuts != buyIns {
		return ErrStackInconsistent
	}

	return nil
}

// findPlayer looks for a player in current game and returns a pointer to it.
func (g *Game) findPlayer(id id.ID) (*Player, error) {
	for i := range g.Players {
		if id == g.Players[i].UserID {
			return &g.Players[i], nil
		}
	}

	return nil, ErrUserNotFound
}

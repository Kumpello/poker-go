package game

import (
	"context"
	"fmt"
	"sync"

	"pokergo/internal/users"
	"pokergo/pkg/id"
	"pokergo/pkg/logger"
)

// Game is internal struct for keeping game state
type Game struct {
	Data

	playerMux    sync.Mutex
	gameLogger   logger.Logger
	usersAdapter users.Adapter
}

// AppendPlayer appends a new player to the game
func (g *Game) AppendPlayer(ctx context.Context, uID *id.ID, name string, startStack int64) error {
	g.playerMux.Lock()
	defer g.playerMux.Unlock()

	for _, p := range g.Players {
		if p.UserName == name {
			return fmt.Errorf("user already exists")
		}
	}

	player := Player{
		BuyIn:             startStack,
		AdditionalIncomes: []inGameTransaction{},
	}
	if uID != nil {
		u, err := g.usersAdapter.GetUserByID(ctx, *uID)
		if err != nil {
			return fmt.Errorf("cannot find user: %w", err)
		}
		player.UserID = &u.ID
		player.UserName = u.Username
	} else {
		player.UserName = name
	}

	g.Players = append(g.Players, player)

	g.gameLogger.Infof("adding new player to the game(anonymous: %b, uID: %s, name: %s startStack: %d)",
		uID == nil, fmt.Sprint(uID), player.UserName, startStack)

	return nil
}

func (g *Game) SetFinishStack(name string, stack int64) error {
	g.playerMux.Lock()
	defer g.playerMux.Unlock()

	if u, err := g.findPlayer(name); err == nil {
		u.BuyOut = &stack
		g.gameLogger.Infof("user %s finishes the game with %d stack", name, stack)
		return nil
	}

	g.gameLogger.Errorf("user %s not exists", name)
	return ErrUserNotFound
}

// ReBuyIn is a re-buy-in from the bank (just increase starting stack)
func (g *Game) ReBuyIn(player string, amount int64) error {
	g.playerMux.Lock()
	defer g.playerMux.Unlock()

	p, err := g.findPlayer(player)
	if err != nil {
		g.gameLogger.Errorf("user %s not exists", p.UserName)
		return ErrUserNotFound
	}

	g.gameLogger.Infof("user %s re-bought for %d", p.UserName, amount)

	p.BuyIn += amount
	return nil
}

// ReBuyInFromPlayer is a re-buy-in but paid by another player
func (g *Game) ReBuyInFromPlayer(buyer, seller string, amount int64) error {
	g.playerMux.Lock()
	defer g.playerMux.Unlock()

	buyerPlayer, err := g.findPlayer(buyer)
	if err != nil {
		g.gameLogger.Errorf("user (buyer) %s not exists", buyerPlayer.UserName)
		return err
	}

	sellerPlayer, err := g.findPlayer(seller)
	if err != nil {
		g.gameLogger.Errorf("user (seller) %s not exists", sellerPlayer.UserName)
		return err
	}

	// add cash to seller
	sellerPlayer.AdditionalIncomes = append(sellerPlayer.AdditionalIncomes, inGameTransaction{
		Amount:   amount,
		Reason:   "re-buy-in from another player",
		From:     buyerPlayer.UserID,
		FromName: buyerPlayer.UserName,
	})

	g.gameLogger.Infof("transaction done: %s --%d--> %s",
		sellerPlayer.UserName, buyerPlayer.UserName, amount)

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
			buyOuts += a.Amount
		}
	}

	if buyOuts != buyIns {
		return ErrStackInconsistent
	}

	return nil
}

func (g *Game) Report() map[string]int64 {
	g.playerMux.Lock()
	defer g.playerMux.Unlock()

	res := make(map[string]int64, len(g.Players))
	for _, p := range g.Players {
		var incomes int64
		for _, a := range p.AdditionalIncomes {
			incomes += a.Amount
		}
		incomes += *p.BuyOut

		res[p.UserName] = p.BuyIn - incomes
	}

	return res
}

// findPlayer looks for a player in current game and returns a pointer to it.
func (g *Game) findPlayer(name string) (*Player, error) {
	for i := range g.Players {
		if name == g.Players[i].UserName {
			return &g.Players[i], nil
		}
	}

	return nil, ErrUserNotFound
}

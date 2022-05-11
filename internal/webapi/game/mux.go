package game

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"pokergo/internal/game"
	"pokergo/internal/webapi/binder"
	"pokergo/pkg/id"
)

type mux struct {
	binder.StructValidator
	gameManager game.Manager
}

func NewMux(validator binder.StructValidator, gameManager game.Manager) *mux {
	return &mux{validator, gameManager}
}

func (m *mux) Route(e *echo.Echo, prefix string) error {
	e.POST(prefix+"/createGame", m.CreateGame)
	e.POST(prefix+"/appendPlayer", m.AppendPlayer)
	e.POST(prefix+"/setFinishStack", m.SetFinishStack)
	e.POST(prefix+"/reBuyIn", m.ReBuyIn)
	e.POST(prefix+"/reBuyInFromPlayer", m.ReBuyInFromPlayer)
	return nil
}

// CreateGame just creates a game for a specific user.
// The game is empty and has no players attached (except the organizer).
func (m *mux) CreateGame(c echo.Context) error {
	data, req, bindErr := binder.BindRequest[createGameRequest](c, true, m)
	if bindErr != nil {
		return c.String(bindErr.Code, bindErr.Message)
	}
	defer data.Cancel()

	g, err := m.gameManager.CreateGame(data.Ctx, data.UserID, req.Org)
	if err != nil {
		return c.String(500, fmt.Sprintf("cannot create a new g: %s", err.Error()))
	}

	return c.JSON(200, createGameResponse{g.ID.Hex()})
}

func (m *mux) AppendPlayer(c echo.Context) error {
	data, req, bindErr := binder.BindRequest[appendPlayerRequest](c, true, m)
	if bindErr != nil {
		return c.String(bindErr.Code, bindErr.Message)
	}
	defer data.Cancel()

	return m.performOnGame(data, req.GameID, func(g *game.Game) (bool, int, string) {
		// No need to verify if requester has the right to the organization - manager do the job.
		isAnonymous := req.UserID == nil
		var i *id.ID
		if !isAnonymous {
			ii, fErr := id.FromString(*req.UserID)
			i = &ii
			if fErr != nil {
				return false, 400, fmt.Sprintf("invalid user id: %s", fErr.Error())
			}
		}

		fErr := g.AppendPlayer(data.Ctx, i, req.UserName, *req.StartStack)
		if fErr != nil {
			return false, 500, fmt.Sprintf("cannot add the player: %s", fErr.Error())
		}

		return true, 200, "ok"
	})
}

func (m *mux) SetFinishStack(c echo.Context) error {
	data, req, bindErr := binder.BindRequest[setFinishStack](c, true, m)
	if bindErr != nil {
		return c.String(bindErr.Code, bindErr.Message)
	}
	defer data.Cancel()

	return m.performOnGame(data, req.GameID, func(g *game.Game) (bool, int, string) {
		if fErr := g.SetFinishStack(req.UserName, *req.FinishStack); fErr != nil {
			return false, 500, fmt.Sprintf("cannot set finish stack: %s", fErr.Error())
		}

		return true, 200, "ok"
	})
}

func (m *mux) ReBuyIn(c echo.Context) error {
	data, req, bindErr := binder.BindRequest[reBuyIn](c, true, m)
	if bindErr != nil {
		return c.String(bindErr.Code, bindErr.Message)
	}
	defer data.Cancel()

	return m.performOnGame(data, req.GameID, func(g *game.Game) (bool, int, string) {
		if fErr := g.ReBuyIn(req.UserName, req.BuyIn); fErr != nil {
			return false, 500, fmt.Sprintf("error on rebuy-in: %s", fErr.Error())
		}

		return true, 200, "ok"
	})
}

func (m *mux) ReBuyInFromPlayer(c echo.Context) error {
	data, req, bindErr := binder.BindRequest[reBuyInFromPlayer](c, true, m)
	if bindErr != nil {
		return c.String(bindErr.Code, bindErr.Message)
	}
	defer data.Cancel()

	return m.performOnGame(data, req.GameID, func(g *game.Game) (bool, int, string) {
		if fErr := g.ReBuyInFromPlayer(req.UserName, req.FromName, req.BuyIn); fErr != nil {
			return false, 500, fmt.Sprintf("error on rebuy-in: %s", fErr.Error())
		}

		return true, 200, "ok"
	})
}

// performOnGame makes direct write on data.Echo
// this should be last call of the function
//
// f should return indicator if the call was ok (if so, the commit is done)
// return code and then string
func (m *mux) performOnGame(
	data binder.Context,
	game string,
	f func(*game.Game) (bool, int, string),
) error {

	gameID, err := id.FromString(game)
	if err != nil {
		return data.Echo.String(400, fmt.Sprintf("invalid game id: %s", err))
	}

	g, err := m.gameManager.GetGame(data.Ctx, data.UserID, gameID)
	if err != nil {
		return data.Echo.String(500, fmt.Sprintf("cannot get the game: %s", err.Error()))
	}

	ok, code, msg := f(g)
	if !ok {
		return data.Echo.String(code, msg)
	}

	err = m.gameManager.Commit(data.Ctx, data.UserID, gameID)
	if err != nil {
		return data.Echo.String(500, fmt.Sprintf("cannot commit the state: %s", err.Error()))
	}

	return data.Echo.String(code, msg)
}

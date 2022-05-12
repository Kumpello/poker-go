package game

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"pokergo/internal/game"
	"pokergo/internal/webapi/binder"
	"pokergo/pkg/id"
)

type mux struct {
	gameManager game.Manager
}

func NewMux(gameManager game.Manager) *mux {
	return &mux{gameManager}
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
	data, bindErr := binder.BindRequest[createGameRequest, any](c, true)
	if bindErr != nil {
		return c.String(bindErr.Code, bindErr.Message)
	}
	defer data.Cancel()

	g, err := m.gameManager.CreateGame(data.Context(), data.UserID(), data.Request.Org)
	if err != nil {
		return c.String(500, fmt.Sprintf("cannot create a new g: %s", err.Error()))
	}

	return c.JSON(200, createGameResponse{g.ID.Hex()})
}

func (m *mux) AppendPlayer(c echo.Context) error {
	data, bindErr := binder.BindRequest[appendPlayerRequest, any](c, true)
	if bindErr != nil {
		return c.String(bindErr.Code, bindErr.Message)
	}
	defer data.Cancel()

	return m.performOnGame(data, data.Request.GameID, func(g *game.Game) (bool, int, string) {
		// No need to verify if requester has the right to the organization - manager do the job.
		isAnonymous := data.Request.UserID == nil
		var i *id.ID
		if !isAnonymous {
			ii, fErr := id.FromString(*data.Request.UserID)
			i = &ii
			if fErr != nil {
				return false, 400, fmt.Sprintf("invalid user id: %s", fErr.Error())
			}
		}

		fErr := g.AppendPlayer(data.Context(), i, data.Request.UserName, *data.Request.StartStack)
		if fErr != nil {
			return false, 500, fmt.Sprintf("cannot add the player: %s", fErr.Error())
		}

		return true, 200, "ok"
	})
}

func (m *mux) SetFinishStack(c echo.Context) error {
	data, bindErr := binder.BindRequest[setFinishStack, any](c, true)
	if bindErr != nil {
		return c.String(bindErr.Code, bindErr.Message)
	}
	defer data.Cancel()

	return m.performOnGame(data, data.Request.GameID, func(g *game.Game) (bool, int, string) {
		if fErr := g.SetFinishStack(data.Request.UserName, *data.Request.FinishStack); fErr != nil {
			return false, 500, fmt.Sprintf("cannot set finish stack: %s", fErr.Error())
		}

		return true, 200, "ok"
	})
}

func (m *mux) ReBuyIn(c echo.Context) error {
	data, bindErr := binder.BindRequest[reBuyIn, any](c, true)
	if bindErr != nil {
		return c.String(bindErr.Code, bindErr.Message)
	}
	defer data.Cancel()

	return m.performOnGame(data, data.Request.GameID, func(g *game.Game) (bool, int, string) {
		if fErr := g.ReBuyIn(data.Request.UserName, data.Request.BuyIn); fErr != nil {
			return false, 500, fmt.Sprintf("error on rebuy-in: %s", fErr.Error())
		}

		return true, 200, "ok"
	})
}

func (m *mux) ReBuyInFromPlayer(c echo.Context) error {
	data, bindErr := binder.BindRequest[reBuyInFromPlayer, any](c, true)
	if bindErr != nil {
		return c.String(bindErr.Code, bindErr.Message)
	}
	defer data.Cancel()

	return m.performOnGame(data, data.Request.GameID, func(g *game.Game) (bool, int, string) {
		if fErr := g.ReBuyInFromPlayer(
			data.Request.UserName,
			data.Request.FromName,
			data.Request.BuyIn,
		); fErr != nil {
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
	binder binder.BaseContext,
	game string,
	f func(*game.Game) (bool, int, string),
) error {

	gameID, err := id.FromString(game)
	if err != nil {
		return binder.Echo().String(400, fmt.Sprintf("invalid game id: %s", err))
	}

	g, err := m.gameManager.GetGame(binder.Context(), binder.UserID(), gameID)
	if err != nil {
		return binder.Echo().String(500, fmt.Sprintf("cannot get the game: %s", err.Error()))
	}

	ok, code, msg := f(g)
	if !ok {
		return binder.Echo().String(code, msg)
	}

	err = m.gameManager.Commit(binder.Context(), binder.UserID(), gameID)
	if err != nil {
		return binder.Echo().String(500, fmt.Sprintf("cannot commit the state: %s", err.Error()))
	}

	return binder.Echo().String(code, msg)
}

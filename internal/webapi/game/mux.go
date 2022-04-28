package game

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"pokergo/internal/game"
)

type mux struct {
	gameManager game.Manager
}

type CreateGameRequest struct {
	// TODO: Remove this and replace with logging and context requests.
	Organizer    string `json:"organizer"`
	Organization string `json:"organization"`
}

type CreateGameResponse struct {
	GameID string `json:"game_id"`
}

// CreateGame just creates a game for a specific user.
// The game is empty and has no players attached (except the organizer).
func (m *mux) CreateGame(c echo.Context) error {
	var request CreateGameRequest
	if err := c.Bind(&request); err != nil {
		return c.String(400, fmt.Sprintf("invalid request: %s", err.Error()))
	}

	reqCtx := c.Request().Context()
	g, err := m.gameManager.CreateGame(reqCtx, request.Organizer)
	if err != nil {
		return c.String(500, fmt.Sprintf("cannot create a new g: %s", err.Error()))
	}

	return c.JSON(200, CreateGameResponse{GameID: g.ID.Hex()})
}

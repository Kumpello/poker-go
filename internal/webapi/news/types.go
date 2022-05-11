package news

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"pokergo/internal/articles"
)

// TODO: Create an echo-binder based on tags with validate

type getURLOpts struct {
	lastDocID string
	no        int
}

func (g *getURLOpts) BindQuery(ctx echo.Context) error {
	err := echo.QueryParamsBinder(ctx).
		String("lastDocID", &g.lastDocID).
		Int("no", &g.no).
		BindError()

	if err != nil {
		return fmt.Errorf("cannot bind the query: %w", err)
	}

	if g.no == 0 {
		g.no = 20
	}

	return nil
}

type newsResponseItem = articles.Article

type getNewsResponse struct {
	News []newsResponseItem `json:"news"`
}

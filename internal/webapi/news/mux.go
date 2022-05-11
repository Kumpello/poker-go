package news

import (
	"github.com/labstack/echo/v4"
	"pokergo/internal/articles"
	"pokergo/internal/webapi/binder"
)

type mux struct {
	binder.StructValidator
	artsAdapter articles.Adapter
}

func NewMux(structValidator binder.StructValidator, artsAdapter articles.Adapter) *mux {
	return &mux{StructValidator: structValidator, artsAdapter: artsAdapter}
}

func (m *mux) Route(e *echo.Echo, prefix string) error {
	e.GET(prefix, m.GetNews)
	return nil
}

// GetNews returns poker-news
func (m *mux) GetNews(c echo.Context) error {
	data, _, bindErr := binder.BindRequest[any](c, false, m)
	if bindErr != nil {
		return c.String(bindErr.Code, bindErr.Message)
	}
	defer data.Cancel()

	var res []newsResponseItem
	arts, err := m.artsAdapter.GetAll(data.Ctx)
	if err != nil {
		return c.String(500, "fetch arts err")
	}
	for _, a := range arts {
		res = append(res, a.Article)
	}

	return c.JSON(200, getNewsResponse{res})
}

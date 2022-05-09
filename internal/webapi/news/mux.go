package news

import (
	"github.com/labstack/echo/v4"
	"pokergo/internal/articles"
	"pokergo/internal/webapi/binder"
)

type mux struct {
	binder.StructValidator
	getters []articles.Getter
}

func NewMux(structValidator binder.StructValidator, getters []articles.Getter) *mux {
	return &mux{StructValidator: structValidator, getters: getters}
}

func (m *mux) Route(e *echo.Echo, prefix string) error {
	e.GET(prefix, m.GetNews)
	return nil
}

// GetNews returns poker-news
// TODO: Cache them instead of scrapping every time
func (m *mux) GetNews(c echo.Context) error {
	var res []newsResponseItem

	for _, v := range m.getters {
		a, err := v.Get()
		if err != nil {
			return c.String(500, "cannot get articles")
		}
		res = append(res, a...)
	}

	return c.JSON(200, getNewsResponse{res})
}

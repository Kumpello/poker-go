package news

import (
	"github.com/labstack/echo/v4"
	"pokergo/internal/articles"
	"pokergo/internal/webapi/binder"
	"pokergo/pkg/id"
)

type mux struct {
	artsAdapter articles.Adapter
}

func NewMux(artsAdapter articles.Adapter) *mux {
	return &mux{artsAdapter: artsAdapter}
}

func (m *mux) Route(e *echo.Echo, prefix string) error {
	e.GET(prefix, m.GetNews)
	return nil
}

// GetNews returns poker-news
// QueryParams:
//	lastDocID = string, default empty (returns from the begging)
//	no = int, default 20, min 5, max 40
func (m *mux) GetNews(c echo.Context) error {
	data, _, bindErr := binder.BindRequest[any](c, false)
	if bindErr != nil {
		return c.String(bindErr.Code, bindErr.Message)
	}
	defer data.Cancel()

	var queryParams getURLOpts
	if err := queryParams.BindQuery(c); err != nil {
		return c.String(400, "invalid query params")
	}
	lastItemID, err := id.FromString(queryParams.lastDocID)
	if err != nil {
		return c.String(400, "unparseable last item id")
	}

	var res []newsResponseItem
	arts, err := m.artsAdapter.GetNext(data.Ctx, lastItemID, queryParams.no)
	if err != nil {
		return c.String(500, "fetch arts err")
	}
	for _, a := range arts {
		res = append(res, a.Article)
	}

	return c.JSON(200, getNewsResponse{res})
}

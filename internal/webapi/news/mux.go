package news

import (
	"github.com/labstack/echo/v4"
	"pokergo/internal/webapi/binder"
)

type mux struct {
	binder.StructValidator
}

func NewMux(structValidator binder.StructValidator) *mux {
	return &mux{StructValidator: structValidator}
}

func (m *mux) Route(e *echo.Echo, prefix string) error {
	e.GET(prefix, m.GetNews)
	return nil
}

// GetNews returns poker-news
// TODO: this is mock :)
func (m *mux) GetNews(c echo.Context) error {
	return c.JSON(200, getNewsResponse{
		News: []newsResponseItem{
			{
				Title: "The Muck: Did \"Poker Brat\" Phil Hellmuth Angle-Shoot vs. Slime on Hustler Casino Live?",
				URL:   "https://www.pokernews.com/news/2022/04/phil-hellmuth-and-tom-dwan-to-play-800k-high-stakes-duel-41098.htm",
				Img:   "https://pnimg.net/w/articles/1/626/f71bdbae23.jpg",
			},
			{
				Title: "Phil Hellmuth & Tom Dwan to Play $800K 'High Stakes Duel' Match on May 12",
				URL:   "https://www.pokernews.com/news/2022/04/phil-hellmuth-and-tom-dwan-to-play-800k-high-stakes-duel-41098.htm",
				Img:   "https://pnimg.net/w/articles/1/626/ec0359b5f4.png",
			},
			{
				Title: "Alec Torelli Asks: \"Did Solvers Ruin Poker?\"",
				URL:   "https://www.pokernews.com/strategy/did-solvers-ruin-poker-41113.htm",
				Img:   "https://pnimg.net/w/articles/1/626/ec0359b5f4.png",
			},
		},
	})
}

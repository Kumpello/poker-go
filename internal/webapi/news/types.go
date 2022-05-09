package news

import (
	"pokergo/internal/articles"
)

type newsResponseItem = articles.Article

type getNewsResponse struct {
	News []newsResponseItem `json:"news"`
}

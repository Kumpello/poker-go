package news

import (
	"pokergo/internal/articles"
)

type getNewsRequest struct {
	LastDocID *string `query:"lastDocID" validate:"hexadecimal,len=24"`
	NO        int     `query:"no" validate:"omitempty,gte=5,lte=40"`
}

type newsResponseItem = articles.Article

type getNewsResponse struct {
	News []newsResponseItem `json:"news"`
}

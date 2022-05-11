package articles

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"pokergo/pkg/id"
)

type Source string

const (
	SourcePokerNews Source = "PokerNews"
)

type Article struct {
	ID       id.ID             `json:"id" bson:"_id"` // nolint:tagliatelle // mongo-id
	Source   Source            `json:"source" bson:"source"`
	IMGSrc   string            `json:"img_src" bson:"img_src"`
	IMGAlt   string            `json:"img_alt" bson:"img_alt"`
	IMGTitle string            `json:"img_title" bson:"img_title"`
	Href     string            `json:"href" bson:"href"`
	Title    string            `json:"title" bson:"title"`
	Metadata map[string]string `json:"metadata" bson:"metadata"`
	Date     time.Time         `json:"date" bson:"date"`
}

type shaArticle struct {
	Article `bson:",inline"`

	ShaHash string `bson:"hash_code"`
}

func convertToSHA(article Article) (shaArticle, error) {
	// remove date from hash-calc
	oldDate := article.Date
	article.Date = time.Time{}

	marshall, err := bson.Marshal(article)
	hash := sha256.New()
	hash.Write(marshall)
	if err != nil {
		return shaArticle{}, fmt.Errorf("cannot conver article to shaArticle: %w", err)
	}

	article.Date = oldDate
	return shaArticle{article, hex.EncodeToString(hash.Sum(nil))}, nil
}

package commands

import (
	"fmt"

	"pokergo/internal/articles"
)

func (c *commandApp) updateArticles() {
	fetchers := []articles.Getter{
		articles.NewPokerNewsExtractor(c.timer),
	}

	var arts []articles.Article
	for _, fetcher := range fetchers {
		a, err := fetcher.Get()
		if err != nil {
			c.logger.Fatalf("fetch articles error: %s", err.Error())
		}
		arts = append(arts, a...)
	}

	ids, err := c.artsAdapter.Save(c.ctx, arts)
	if err != nil {
		c.logger.Fatalf("save article error: %s", err.Error())
	}

	c.logger.Infof("inserted ids: %s", fmt.Sprint(ids))
}

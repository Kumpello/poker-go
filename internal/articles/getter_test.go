package articles

import (
	"path/filepath"
	"testing"
)

func Test_PokerNewsGetter(t *testing.T) {
	getter := NewPokerNewsExtractor()
	// replace provider to not use real-url
	getter.provider = fileProvider{}
	abs, _ := filepath.Abs("./poker_news.html.test")
	getter.fetchURL = abs

	articles, err := getter.Get()
	if err != nil {
		t.Fatalf("cannot get articles: %s", err.Error())
		return
	}

	if len(articles) != 20 {
		return
	}
}

package articles

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

type htmlGetter func(document *goquery.Document) ([]Article, error)

func extractFromHTML(provider provider, url string, getter htmlGetter) ([]Article, error) {
	res, err := provider.provide(url)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch input")
	}

	defer res.Close()

	doc, err := goquery.NewDocumentFromReader(res)
	if err != nil {
		return nil, fmt.Errorf("cannot load doc: %w", err)
	}

	docs, err := getter(doc)
	if err != nil {
		return nil, fmt.Errorf("getter err: %w", err)
	}
	return docs, nil
}

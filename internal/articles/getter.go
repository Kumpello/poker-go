package articles

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Getter interface {
	Get() ([]Article, error)
}

type PokerNewsGetter struct {
	fetchURL string
	provider provider
	getter   func(provider provider, url string, getter htmlGetter) ([]Article, error)
	newsURL  string
}

func NewPokerNewsExtractor() *PokerNewsGetter {
	return &PokerNewsGetter{
		fetchURL: "https://www.pokernews.com/news/",
		newsURL:  "https://www.pokernews.com/news/",
		provider: &htmlProvider{},
		getter:   extractFromHTML,
	}
}

func (p PokerNewsGetter) Get() ([]Article, error) {
	return p.getter(p.provider, p.fetchURL, p.extract)
}

func (p PokerNewsGetter) extract(document *goquery.Document) ([]Article, error) {
	var res []Article
	document.Find(".articles.clrfix li").Each(func(i int, selection *goquery.Selection) {
		a := selection.Find("a")
		img := a.Find("img")
		div := selection.Find("div")

		imgSrc, _ := img.Attr("data-src")
		imgAlt, _ := img.Attr("alt")
		imgTitle, _ := img.Attr("title")
		href, _ := a.Attr("href")
		category := div.Find("p").Text()
		title := div.Text()

		res = append(res, Article{
			IMGSrc:   imgSrc,
			IMGAlt:   imgAlt,
			IMGTitle: imgTitle,
			Href:     href,
			Title:    title,
			Metadata: map[string]string{
				"category": strings.Trim(category, " "),
			},
		})
	})

	return res, nil
}

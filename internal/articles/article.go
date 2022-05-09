package articles

type Article struct {
	IMGSrc   string            `json:"img_src"`
	IMGAlt   string            `json:"img_alt"`
	IMGTitle string            `json:"img_title"`
	Href     string            `json:"href"`
	Title    string            `json:"title"`
	Metadata map[string]string `json:"metadata"`
}

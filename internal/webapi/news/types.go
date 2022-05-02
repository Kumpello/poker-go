package news

type newsResponseItem struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Img   string `json:"img"`
}

type getNewsResponse struct {
	News []newsResponseItem `json:"news"`
}

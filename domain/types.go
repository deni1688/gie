package domain

type Repo struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Owner string `json:"owner"`
}

type Issue struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Url           string `json:"url"`
	ExtractedLine string `json:"-"`
}

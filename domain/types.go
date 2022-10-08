package domain

type Repo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Issue struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

func (r *Issue) Reset() {
	r.Title = ""
	r.Desc = ""
}

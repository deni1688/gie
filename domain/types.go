package domain

type Repo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Issue struct {
	Title     string `json:"title"`
	Desc      string `json:"desc"`
	Weight    string `json:"weight"`
	Milestone string `json:"milestone"`
}

func (r *Issue) Reset() {
	r.Title = ""
	r.Desc = ""
	r.Weight = ""
	r.Milestone = ""
}

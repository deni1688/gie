package main

type Repo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Issue struct {
	Title     string `json:"title"`
	Desc      string `json:"desc"`
	Weight    int    `json:"weight"`
	Milestone string `json:"milestone"`
}

package main

type GitProvider interface {
	GetRepos() (*[]Repo, error)
	CreateIssue(repo Repo, issue Issue) error
}

type Service interface {
	ListRepos() (*[]Repo, error)
	DefaultIssue() Issue
	SubmitIssues(Repo, *[]Issue) error
}

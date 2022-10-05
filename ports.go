package main

type GitProvider interface {
	GetRepos() (*[]Repo, error)
	CreateIssue(repo Repo, issue Issue) error
	CreateIssues(repo Repo, issues []Issue) error
}

type Service interface {
	ParseArgs(args []string) ([]Issue, error)
	Execute([]Issue) error
}

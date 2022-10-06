package domain

type Provider interface {
	GetRepos() (*[]Repo, error)
	CreateIssue(repo Repo, issue Issue) error
}

type Service interface {
	ListRepos() (*[]Repo, error)
	DefaultIssue() Issue
	SubmitIssues(Repo, *[]Issue) error
}

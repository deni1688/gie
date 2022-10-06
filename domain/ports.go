package domain

type Provider interface {
	GetRepos() (*[]Repo, error)
	CreateIssue(repo Repo, issue Issue) error
}

type Service interface {
	ListRepos() (*[]Repo, error)
	SubmitIssues(Repo, *[]Issue) error
}

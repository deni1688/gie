package domain

type GitHost interface {
	GetRepos() (*[]Repo, error)
	CreateIssue(repo Repo, issue Issue) error
}

type Notifier interface {
	Notify(issues *[]Issue) error
}

type Service interface {
	Notify(issues *[]Issue) error
	ListRepos() (*[]Repo, error)
	ExtractIssues(content, source string) (*[]Issue, error)
	SubmitIssues(Repo, *[]Issue) error
}

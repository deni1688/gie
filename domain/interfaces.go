package domain

type GitProvider interface {
	GetRepos() (*[]Repo, error)
	CreateIssue(repo *Repo, issue *Issue) error
}

type Notifier interface {
	Notify(issues *[]Issue) error
}

type Service interface {
	Notify(issues *[]Issue) error
	ListRepos() (*[]Repo, error)
	ExtractIssues(content, source string) (*[]Issue, error)
	SubmitIssue(repo *Repo, issue *Issue) error
	FindRepoByName(name string) (*Repo, error)
}

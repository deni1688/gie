package issues

type GitProvider interface {
	GetRepos() (*[]Repo, error)
	CreateIssue(repo *Repo, issue *Issue) error
}

type Notifier interface {
	Notify(issues *[]Issue) error
}

type Service interface {
	ExtractIssues(content, source *string) (*[]Issue, error)
	FindRepoByName(name string) (*Repo, error)
	SubmitIssue(repo *Repo, issue *Issue) error
	GetUpdatedLine(issue Issue) string
	Notify(issues *[]Issue) error
	listRepos() (*[]Repo, error)
}

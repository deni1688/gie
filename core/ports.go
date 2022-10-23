package core

// GitProvider is an interface for interacting with a git provider, e.g. GitHub or GitLab
type GitProvider interface {
	GetRepos() (*[]Repo, error)
	CreateIssue(repo *Repo, issue *Issue) error
}

// Notifier is an interface for notifying about issues to be invoked after issue extraction
type Notifier interface {
	Notify(issues *[]Issue) error
}

// Service is the interface for the core service
type Service interface {
	// ExtractIssues extracts issues from the provided string content reference and returns a list of issues
	ExtractIssues(content, source *string) (*[]Issue, error)
	// FindRepoByName returns a repo if found
	FindRepoByName(name string) (*Repo, error)
	// SubmitIssue submits an issue to the provided repo using gitProvider
	SubmitIssue(repo *Repo, issue *Issue) error
	// GetUpdatedLine returns the updated line with the issue ID and URL
	GetUpdatedLine(issue Issue) string
	// Notify uses the notifier to publish the issues to other interested parties
	Notify(issues *[]Issue) error
}

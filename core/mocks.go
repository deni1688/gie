package core

type mockGitProvider struct {
	repos []Repo
	issue Issue
	err   error
}

func (m *mockGitProvider) GetRepos() (*[]Repo, error) {
	return &m.repos, m.err
}

func (m *mockGitProvider) CreateIssue(repo *Repo, issue *Issue) error {
	issue.ID = m.issue.ID
	issue.Url = m.issue.Url
	return m.err
}

type mockNotifier struct {
	err error
}

func (m mockNotifier) Notify(issues *[]Issue) error {
	return m.err
}
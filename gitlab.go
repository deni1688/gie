package main

type gitlabProvider struct {
	Token string
}

func NewGitlabProvider(token string) Provider {
	return &gitlabProvider{
		Token: token,
	}
}

func (r *gitlabProvider) GetRepos() (*[]Repo, error) {
	return &[]Repo{
		{Name: "Repo 1", ID: 1},
		{Name: "Repo 2", ID: 2},
		{Name: "Repo 3", ID: 3},
		{Name: "Repo 4", ID: 4},
	}, nil
}

func (r *gitlabProvider) CreateIssues(repo Repo, issues []Issue) error {
	var err error

	for _, issue := range issues {
		err = r.CreateIssue(repo, issue)
	}

	return err
}

func (r *gitlabProvider) CreateIssue(repo Repo, issue Issue) error {
	return nil
}

package infra

import (
	"deni1688/gitissue/domain"
	"net/http"
)

type githubProvider struct {
	token  string
	host   string
	query  string
	client *http.Client
}

func NewGithubProvider(token string, host string, query string) domain.Provider {
	return &githubProvider{token, host, query, http.DefaultClient}
}

func (g githubProvider) GetRepos() (*[]domain.Repo, error) {
	//TODO implement me
	panic("implement me")
}

func (g githubProvider) CreateIssue(repo domain.Repo, issue domain.Issue) error {
	//TODO implement me
	panic("implement me")
}

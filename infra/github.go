package infra

import (
	"deni1688/gitissue/domain"
	"encoding/json"
	"net/http"
)

type github struct {
	token  string
	host   string
	query  string
	client *http.Client
}

func NewGithub(token string, host string, query string) domain.GitHost {
	return &github{token, host, query, http.DefaultClient}
}

func (r github) GetRepos() (*[]domain.Repo, error) {
	req, err := r.request("GET", "/user/repos")
	if err != nil {
		return nil, err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}

	var repos []domain.Repo
	if err = json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, err
	}

	return &repos, nil
}

func (r github) CreateIssue(repo domain.Repo, issue domain.Issue) error {
	return nil
}

func (r github) request(method, resource string) (*http.Request, error) {
	req, err := http.NewRequest(method, r.host+resource, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+r.token)
	req.URL.RawQuery = r.query

	return req, err
}

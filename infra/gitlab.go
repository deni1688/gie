package infra

import (
	"deni1688/gitissue/domain"
	"encoding/json"
	"net/http"
)

type gitlab struct {
	token  string
	host   string
	query  string
	client *http.Client
}

func NewGitlab(token, host, query string) domain.GitHost {
	return &gitlab{token, host, query, http.DefaultClient}
}

func (r gitlab) GetRepos() (*[]domain.Repo, error) {
	req, err := r.request("GET", "projects")
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

func (r gitlab) CreateIssue(repo domain.Repo, issue domain.Issue) error {
	return nil
}

func (r gitlab) request(method, resource string) (*http.Request, error) {
	req, err := http.NewRequest(method, r.endpoint(resource), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("PRIVATE-TOKEN", r.token)
	req.URL.RawQuery = r.query

	return req, err
}

func (r gitlab) endpoint(resource string) string {
	return r.host + "/api/v4/" + resource
}

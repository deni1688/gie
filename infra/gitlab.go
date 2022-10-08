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
	client HttpClient
}

func NewGitlab(token, host, query string, client HttpClient) domain.GitProvider {
	return &gitlab{token, host, query, client}
}

func (r gitlab) GetRepos() (*[]domain.Repo, error) {
	req, err := r.request("GET", "projects")
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = r.query
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

func (r gitlab) CreateIssue(repo *domain.Repo, issue *domain.Issue) error {
	return nil
}

func (r gitlab) request(method, resource string) (*http.Request, error) {
	req, err := http.NewRequest(method, r.endpoint(resource), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("PRIVATE-TOKEN", r.token)

	return req, err
}

func (r gitlab) endpoint(resource string) string {
	return r.host + "/api/v4/" + resource
}

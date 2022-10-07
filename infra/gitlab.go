package infra

import (
	"deni1688/gitissue/domain"
	"encoding/json"
	"net/http"
)

type gitlabProvider struct {
	token  string
	host   string
	query  string
	client *http.Client
}

func NewGitlabProvider(token, host, query string) domain.Provider {
	return &gitlabProvider{token, host, query, http.DefaultClient}
}

func (r gitlabProvider) GetRepos() (*[]domain.Repo, error) {
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

func (r gitlabProvider) CreateIssue(repo domain.Repo, issue domain.Issue) error {
	return nil
}

func (r gitlabProvider) request(method, resource string) (*http.Request, error) {
	req, err := http.NewRequest(method, r.endpoint(resource), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("PRIVATE-TOKEN", r.token)
	req.URL.RawQuery = r.query
	return req, err

}

func (r gitlabProvider) endpoint(resource string) string {
	return r.host + "/api/v4/" + resource
}

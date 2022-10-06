package infra

import (
	"deni1688/gitissue/domain"
	"encoding/json"
	"net/http"
)

type gitlabProvider struct {
	token  string
	host   string
	client *http.Client
}

func NewGitlabProvider(token string, host string) domain.Provider {
	return &gitlabProvider{token, host, http.DefaultClient}
}

func (r gitlabProvider) GetRepos() (*[]domain.Repo, error) {
	req, err := http.NewRequest(http.MethodGet, r.host+"/api/v4/projects", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("PRIVATE-TOKEN", r.token)

	q := req.URL.Query()
	q.Add("per_page", "100")
	q.Add("order_by", "name")
	q.Add("archived", "false")
	q.Add("sort", "asc")
	q.Add("visibility", "private")
	req.URL.RawQuery = q.Encode()

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

func (r gitlabProvider) GetTokenUserId() (string, error) {
	return "123", nil
}

package gitlab

import (
	"bytes"
	"deni1688/gie/adapters/shared"
	"deni1688/gie/core"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type gitlab struct {
	token  string
	host   string
	query  string
	client shared.HttpClient
	repos  *[]core.Repo
}

type gitlabIssue struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Desc   string `json:"description"`
	WebUrl string `json:"web_url"`
}

func New(token, host, query string, client shared.HttpClient) core.GitProvider {
	return &gitlab{token, host, query, client, nil}
}

func (r gitlab) GetRepos() (*[]core.Repo, error) {
	if r.repos != nil {
		return r.repos, nil
	}

	req, err := r.request("GET", "projects")
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = r.query
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var repos []core.Repo
	if err = json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, err
	}

	return &repos, nil
}

func (r gitlab) CreateIssue(repo *core.Repo, issue *core.Issue) error {
	req, err := r.request("POST", fmt.Sprintf("projects/%d/issues", repo.ID))
	if err != nil {
		return err
	}

	body, err := json.Marshal(gitlabIssue{Title: issue.Title, Desc: issue.Desc})
	if err != nil {
		return err
	}

	req.Body = io.NopCloser(bytes.NewReader(body))

	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var createdIssue gitlabIssue
	if err = json.NewDecoder(resp.Body).Decode(&createdIssue); err != nil {
		return err
	}

	issue.ID = createdIssue.ID
	issue.Url = createdIssue.WebUrl

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

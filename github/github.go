package github

import (
	"bytes"
	"deni1688/gie/common"
	"deni1688/gie/internal/issues"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type github struct {
	token  string
	host   string
	query  string
	client common.HttpClient
	repos  *[]issues.Repo
}

func New(token string, host string, query string, client common.HttpClient) issues.GitProvider {
	return &github{token, host, query, client, nil}
}

type githubRepo struct {
	*issues.Repo
	Owner struct {
		Login string `json:"login"`
	} `json:"owner"`
}

type githubIssue struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	HtmlUrl string `json:"html_url"`
}

func (r github) GetRepos() (*[]issues.Repo, error) {
	if r.repos != nil {
		return r.repos, nil
	}

	req, err := r.request("GET", "/user/repos")
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = r.query
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var githubRepos []githubRepo
	if err = json.NewDecoder(resp.Body).Decode(&githubRepos); err != nil {
		return nil, err
	}

	repos := make([]issues.Repo, 0)
	for _, repo := range githubRepos {
		repos = append(repos, issues.Repo{ID: repo.ID, Name: repo.Name, Owner: repo.Owner.Login})
	}

	return &repos, nil
}

func (r github) CreateIssue(repo *issues.Repo, issue *issues.Issue) error {
	req, err := r.request("POST", "/repos/"+repo.Owner+"/"+repo.Name+"/issues")
	if err != nil {
		return err
	}

	body, err := json.Marshal(githubIssue{Title: issue.Title, Body: issue.Desc})
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

	var createdIssue githubIssue
	if err = json.NewDecoder(resp.Body).Decode(&createdIssue); err != nil {
		return err
	}

	issue.ID = createdIssue.ID
	issue.Url = createdIssue.HtmlUrl

	return nil
}

func (r github) request(method, resource string) (*http.Request, error) {
	req, err := http.NewRequest(method, r.host+resource, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	req.Header.Add("Authorization", "Bearer "+r.token)

	return req, nil
}

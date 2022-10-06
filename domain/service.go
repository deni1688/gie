package domain

import (
	"fmt"
	"time"
)

type service struct {
	provider Provider
}

func NewService(provider Provider) Service {
	return &service{provider: provider}
}

func (r service) ListRepos() (*[]Repo, error) {
	return r.provider.GetRepos()
}

func (r service) SubmitIssues(repo Repo, issues *[]Issue) error {
	for _, issue := range *issues {
		fmt.Println(fmt.Sprintf("Creating issue %s in %s", issue.Title, repo.Name))
		time.Sleep(1 * time.Second)
		err := r.provider.CreateIssue(repo, issue)
		if err != nil {
			return err
		}
	}

	return nil
}

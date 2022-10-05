package main

import (
	"fmt"
	"time"
)

type service struct {
	provider GitProvider
}

func NewService(provider GitProvider) Service {
	return &service{
		provider: provider,
	}
}

func (r *service) DefaultIssue() Issue {
	return Issue{
		Title:     "Default title",
		Desc:      "Default description",
		Weight:    15,
		Milestone: "",
	}
}

func (r *service) ListRepos() (*[]Repo, error) {
	return r.provider.GetRepos()
}

func (r *service) SubmitIssues(repo Repo, issues *[]Issue) error {
	fmt.Printf("Submitting %d issues to repo: %s\n\n", len(*issues), repo.Name)
	for _, issue := range *issues {
		fmt.Println("Creating issue: ", issue.Title)
		time.Sleep(1 * time.Second)
		err := r.provider.CreateIssue(repo, issue)
		if err != nil {
			return err
		}
	}

	return nil
}

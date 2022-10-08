package domain

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type service struct {
	gitHost  GitHost
	notifier Notifier
	prefix   string
}

func NewService(gitHost GitHost, notifier Notifier, prefix string) Service {
	return &service{gitHost, notifier, prefix}
}

func (r service) ListRepos() (*[]Repo, error) {
	return r.gitHost.GetRepos()
}

func (r service) SubmitIssues(repo Repo, issues *[]Issue) error {
	for _, issue := range *issues {
		fmt.Println(fmt.Sprintf("Creating issue %s in %s", issue.Title, repo.Name))
		time.Sleep(1 * time.Second)
		err := r.gitHost.CreateIssue(repo, issue)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r service) ExtractIssues(content, source string) (*[]Issue, error) {
	regx, err := regexp.Compile(r.prefix + "(.*)\n")

	var issues []Issue
	if err != nil {
		return nil, err
	}

	if strings.Contains(content, r.prefix) {
		foundIssues := regx.FindAllString(content, -1)
		for _, title := range foundIssues {
			issue := Issue{}
			issue.Title = strings.TrimPrefix(title, r.prefix)
			issue.Desc = "Extracted from " + source
			issues = append(issues, issue)
		}
	}

	return &issues, nil
}

func (r service) Notify(issues *[]Issue) error {
	return r.notifier.Notify(issues)
}

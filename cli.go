package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Cli struct {
	service Service
}

func NewCli(service Service) *Cli {
	return &Cli{service}
}

const helpError = `please provide one or more issues to be created, format: t:"New issue title" d:"This should be fixed" w:30 m:"sprint/45" -- ..."`

func (r *Cli) Run() error {
	args := os.Args[1:]

	if len(args) == 0 {
		return errors.New(helpError)
	}

	issues, err := r.issuesFromArgs(args)
	if err != nil {
		return err
	}

	repos, err := r.service.ListRepos()
	fmt.Printf("Found %d repos\n", len(*repos))
	for i, repo := range *repos {
		fmt.Printf("%d: %s\n", i, repo.Name)
	}

	var repoIndex string
	fmt.Printf("Select a repo to create issues in (0-%d): ", len(*repos)-1)
	if _, err := fmt.Scanln(&repoIndex); err != nil {
		return err
	}

	index, err := strconv.Atoi(repoIndex)
	if err != nil || index < 0 || index >= len(*repos) {
		return fmt.Errorf("invalid repo index: %s", repoIndex)
	}

	repo := (*repos)[index]
	return r.service.SubmitIssues(repo, &issues)
}

func (r *Cli) issuesFromArgs(args []string) ([]Issue, error) {
	var issues []Issue
	issue := r.service.DefaultIssue()

	for _, arg := range args {
		switch {
		case setFromArgWithPrefix(arg, "t:", &issue.Title):
		case setFromArgWithPrefix(arg, "d:", &issue.Desc):
		case setFromArgWithPrefix(arg, "w:", &issue.Weight):
		case setFromArgWithPrefix(arg, "m:", &issue.Milestone):
		case arg == "--":
			issues = append(issues, issue)
			issue = r.service.DefaultIssue()
		default:
			return nil, errors.New("invalid argument: " + arg)
		}
	}

	return issues, nil
}

func setFromArgWithPrefix(arg, prefix string, value *string) bool {
	if strings.HasPrefix(arg, prefix) {
		*value = strings.TrimPrefix(arg, prefix)
		return true
	}

	return false
}

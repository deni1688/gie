package main

import (
	"fmt"
	"strconv"
	"strings"
)

type service struct {
	provider GitProvider
}

func NewService(provider GitProvider) Service {
	return &service{
		provider: provider,
	}
}

func (r *service) ParseArgs(args []string) ([]Issue, error) {
	issues := []Issue{}
	issue := getDefaultIssue()

	for _, arg := range args {
		switch {
		case strings.HasPrefix(arg, "-t="):
			issue.Title = strings.TrimPrefix(arg, "-t=")
		case strings.HasPrefix(arg, "-d="):
			issue.Desc = strings.TrimPrefix(arg, "-d=")
		case strings.HasPrefix(arg, "-w="):
			weight, err := strconv.Atoi(strings.TrimPrefix(arg, "-w="))
			if err != nil {
				fmt.Println("Weight must be an integer")
				return nil, err
			}

			issue.Weight = weight
		case strings.HasPrefix(arg, "-m="):
			issue.Milestone = strings.TrimPrefix(arg, "-m=")
		case arg == "--":
			issues = append(issues, issue)
			issue = getDefaultIssue()
		default:
			fmt.Println("Invalid argument: ", arg)
		}
	}

	return issues, nil
}

func (r *service) Execute(issues []Issue) error {
	repos, err := r.provider.GetRepos()
	for i, repo := range *repos {
		fmt.Printf("%d: %s\n", i, repo.Name)
	}

	var repoIndex string
	fmt.Print("Select a repo to create issues in: ")
	fmt.Scanln(&repoIndex)

	index, err := strconv.Atoi(repoIndex)
	if err != nil || index < 0 || index >= len(*repos) {
		return fmt.Errorf("Invalid repo index: %s", repoIndex)
	}

	repo := (*repos)[index]
	fmt.Println("Creating issues in repo: ", repo.Name)
	fmt.Println(issues)

	return r.provider.CreateIssues(repo, issues)
}

func getDefaultIssue() Issue {
	return Issue{
		Title:     "Default title",
		Desc:      "Default desc",
		Weight:    10,
		Milestone: "Default milestone",
	}
}

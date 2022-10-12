package issues

import (
	"fmt"
	"path"
	"regexp"
	"strings"
)

type service struct {
	gitProvider GitProvider
	notifier    Notifier
	prefix      string
}

func NewService(gitProvider GitProvider, notifier Notifier, prefix string) Service {
	return &service{gitProvider, notifier, prefix}
}

func (r service) listRepos() (*[]Repo, error) {
	return r.gitProvider.GetRepos()
}

func (r service) SubmitIssue(repo *Repo, issue *Issue) error {
	fmt.Printf("Submitting issue=[%s] to repo=[%s]\n", issue.Title, repo.Name)
	if err := r.gitProvider.CreateIssue(repo, issue); err != nil {
		return err
	}

	fmt.Printf("Issue created at url=[%s]\n", issue.Url)

	return nil
}

func (r service) ExtractIssues(content, source *string) (*[]Issue, error) {
	var issues []Issue
	regx, err := regexp.Compile(r.prefix + "(.*)\n")
	if err != nil {
		return nil, err
	}

	issuesMap := make(map[string]Issue)
	if strings.Contains(*content, r.prefix) {
		foundIssues := regx.FindAllString(*content, -1)
		for _, title := range foundIssues {
			if strings.Contains(title, " -> ") || issuesMap[title] != (Issue{}) {
				continue
			}

			issue := Issue{}
			trimmedTitle := strings.Trim(strings.TrimPrefix(title, r.prefix), " \n")
			issue.Title = strings.ToUpper(trimmedTitle[:1]) + trimmedTitle[1:]
			issue.Desc = "Extracted from " + *source
			issue.ExtractedLine = title
			issues = append(issues, issue)
			issuesMap[title] = issue
		}
	}

	return &issues, nil
}

func (r service) Notify(issues *[]Issue) error {
	return r.notifier.Notify(issues)
}

func (r service) FindRepoByName(name string) (*Repo, error) {
	repos, err := r.listRepos()
	if err != nil {
		return &Repo{}, err
	}

	var current Repo
	base := path.Base(name)
	for _, repo := range *repos {
		if strings.Contains(base, repo.Name) {
			current = repo
			return &current, nil
		}
	}

	return &Repo{}, fmt.Errorf("repo=[%s] not found", name)
}

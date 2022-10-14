package issues

import (
	"fmt"
	"path"
	"regexp"
	"strings"
)

const LABEL = "-> closes"

type service struct {
	gitProvider GitProvider
	notifier    Notifier
	prefix      string
}

func New(gitProvider GitProvider, notifier Notifier, prefix string) Service {
	return &service{gitProvider, notifier, prefix}
}

func (r service) listRepos() (*[]Repo, error) {
	return r.gitProvider.GetRepos()
}

func (r service) SubmitIssue(repo *Repo, issue *Issue) error {
	fmt.Printf("Submitting issue=[%s] to repo=[%s]\n", issue.Title, repo.Name)
	if err := r.gitProvider.CreateIssue(repo, issue); err != nil {
		return fmt.Errorf("failed to create issue for repo=[%s] with error=[%s]", repo.Name, err)
	}
	fmt.Printf("Issue created at url=[%s]\n", issue.Url)

	return nil
}

func (r service) ExtractIssues(content, source *string) (*[]Issue, error) {
	regx, err := regexp.Compile(r.prefix + "(.*)\n")
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex with provided prefix=[%s] with error=[%s]", r.prefix, err)
	}

	var issues []Issue
	issuesMap := make(map[string]Issue)
	if strings.Contains(*content, r.prefix) {
		foundIssues := regx.FindAllString(*content, -1)
		for _, title := range foundIssues {
			if strings.Contains(title, LABEL) || issuesMap[title] != (Issue{}) {
				continue
			}

			issue := Issue{}
			trimmedTitle := strings.Trim(strings.TrimPrefix(title, r.prefix), " \n")

			if trimmedTitle == "" {
				continue
			}

			issue.Title = strings.ToUpper(trimmedTitle[:1]) + trimmedTitle[1:]
			issue.Desc = "Extracted from " + *source
			issue.ExtractedLine = title
			issues = append(issues, issue)
			issuesMap[title] = issue
		}
	}

	return &issues, nil
}

func (r service) GetUpdatedLine(issue Issue) string {
	return fmt.Sprintf("%s %s %s\n",
		strings.Trim(issue.ExtractedLine, "\n"),
		LABEL,
		issue.Url)
}

func (r service) FindRepoByName(name string) (*Repo, error) {
	repos, err := r.listRepos()
	if err != nil {
		return &Repo{}, fmt.Errorf("failed to list repos with error=[%s]", err)
	}

	if len(*repos) < 1 {
		return &Repo{}, fmt.Errorf("no repos found")
	}

	base := path.Base(name)
	for _, repo := range *repos {
		if strings.Contains(base, repo.Name) {
			return &repo, nil
		}
	}

	return &Repo{}, fmt.Errorf("repo=[%s] not found", name)
}

func (r service) Notify(issues *[]Issue) error {
	if err := r.notifier.Notify(issues); err != nil {
		return fmt.Errorf("failed to notify with error=[%s]", err)
	}

	return nil
}

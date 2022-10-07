package infra

import (
	"deni1688/gitissue/domain"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

type Cli struct {
	prefix  string
	path    string
	service domain.Service
}

func NewCli(prefix, path string, service domain.Service) *Cli {
	return &Cli{prefix, path, service}
}

// Issue: #2
func (r Cli) Run() error {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	origin, err := cmd.Output()
	if err != nil {
		return err
	}

	b, err := os.ReadFile(r.path)
	if err != nil {
		return err
	}

	issues, err := r.parseIssues(string(b))
	if err != nil {
		return err
	}

	repos, err := r.service.ListRepos()
	if err != nil {
		return err
	}

	var currentRepo domain.Repo
	for _, repo := range *repos {
		if strings.Contains(path.Base(string(origin)), repo.Name) {
			fmt.Println("Found repo: ", repo)
			currentRepo = repo
			break
		}
	}
	fmt.Println(currentRepo)

	err = r.service.SubmitIssues(domain.Repo{Name: string(origin)}, &issues)
	if err != nil {
		return err
	}

	return nil
}

// Issue: #1
func (r Cli) parseIssues(content string) ([]domain.Issue, error) {
	var issues []domain.Issue
	issue := domain.Issue{}
	regx, err := regexp.Compile(r.prefix + "(.*)\n")
	if err != nil {
		return nil, err
	}

	if strings.Contains(content, r.prefix) {
		issueLines := regx.FindAllString(content, -1)

		for _, title := range issueLines {
			issue.Title = strings.TrimPrefix(title, r.prefix)
			issue.Desc = "Extracted from " + r.path
			issues = append(issues, issue)
			issue.Reset()
		}
	}

	return issues, nil
}

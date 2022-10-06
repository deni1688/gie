package infra

import (
	"deni1688/gitissue/domain"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Cli struct {
	prefix  string
	service domain.Service
}

func NewCli(prefix string, service domain.Service) *Cli {
	return &Cli{prefix, service}
}

//ISSUE: #2
func (r Cli) Run() error {
	p := flag.String("path", "./issues.txt", "please provide file path to parse issues from")
	flag.Parse()

	fmt.Println(*p)

	cmd := exec.Command("git", "remote", "get-url", "origin")
	origin, err := cmd.Output()
	if err != nil {
		return err
	}

	b, err := os.ReadFile(*p)
	if err != nil {
		return err
	}

	issues, err := r.parseIssues(strings.Split(string(b), "\n"), p)
	if err != nil {
		return err
	}

	err = r.service.SubmitIssues(domain.Repo{Name: string(origin)}, &issues)
	if err != nil {
		return err
	}

	return nil
}

//ISSUE: #1
func (r Cli) parseIssues(lines []string, path *string) ([]domain.Issue, error) {
	var issues []domain.Issue
	issue := domain.Issue{}

	for i, line := range lines {
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, r.prefix) {
			issue.Title = strings.TrimPrefix(line, r.prefix)
			issue.Desc = fmt.Sprintf("Extracted from line %d in %s", i, *path)
			issues = append(issues, issue)
			issue.Reset()
		}
	}

	return issues, nil
}

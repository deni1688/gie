package infra

import (
	"deni1688/gitissue/domain"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
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

	issues, err := r.parseIssues(string(b), p)
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
func (r Cli) parseIssues(content string, path *string) ([]domain.Issue, error) {
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
			issue.Desc = "Extracted from " + *path
			issues = append(issues, issue)
			issue.Reset()
		}
	}

	return issues, nil
}

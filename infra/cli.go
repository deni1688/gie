package infra

import (
	"deni1688/gogie/domain"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Cli struct {
	path    string
	service domain.Service
}

func NewCli(path string, service domain.Service) *Cli {
	return &Cli{path, service}
}

func (r Cli) Execute() error {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	origin, err := cmd.Output()
	if err != nil {
		return err
	}

	b, err := os.ReadFile(r.path)
	if err != nil {
		return err
	}

	content := string(b)
	name := string(origin)

	issues, err := r.service.ExtractIssues(content, r.path)

	if len(*issues) < 1 {
		fmt.Println("No issues found")
		return nil
	}

	repo, err := r.service.FindRepoByName(name)

	for _, issue := range *issues {
		fmt.Printf("\n")
		if err = r.service.SubmitIssue(repo, &issue); err != nil {
			return err
		}

		updatedLine := fmt.Sprintf("%s -> %s\n", strings.Trim(issue.ExtractedLine, "\n"), issue.Url)
		content = strings.Replace(content, issue.ExtractedLine, updatedLine, 1)
	}

	fmt.Printf("\n")

	if err = r.service.Notify(issues); err != nil {
		return err
	}

	return os.WriteFile(r.path, []byte(content), 0644)
}

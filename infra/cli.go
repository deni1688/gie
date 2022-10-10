package infra

import (
	"deni1688/gogie/internal/issues"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Cli struct {
	service issues.Service
}

func NewCli(service issues.Service) *Cli {
	return &Cli{service}
}

func (r Cli) Execute(path string) error {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	origin, err := cmd.Output()
	if err != nil {
		return err
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	content := string(b)
	name := string(origin)

	issues, err := r.service.ExtractIssues(content, path)

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

	return os.WriteFile(path, []byte(content), 0644)
}

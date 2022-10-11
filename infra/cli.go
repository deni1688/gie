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
	dry     bool
}

func NewCli(service issues.Service, dry bool) *Cli {
	return &Cli{service, dry}
}

func (r Cli) Execute(path string) error {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	origin, err := cmd.Output()
	if err != nil {
		return err
	}

	fi, err := os.Stat(path)
	if err != nil {
		return err
	}

	if fi.IsDir() {
		var files []os.DirEntry
		files, err = os.ReadDir(path)
		if err != nil {
			return err
		}

		for _, file := range files {
			err = r.Execute(path + "/" + file.Name())
			if err != nil {
				return err
			}
		}

		return nil
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	content := string(b)
	name := string(origin)

	foundIssues, err := r.service.ExtractIssues(&content, &path)
	if len(*foundIssues) < 1 {
		return nil
	}

	if r.dry {
		for _, issue := range *foundIssues {
			fmt.Printf("Found issue=[%s] in file=[%s]\n", issue.Title, path)
		}

		return nil
	}

	repo, err := r.service.FindRepoByName(name)
	for _, issue := range *foundIssues {
		fmt.Printf("\n")
		if err = r.service.SubmitIssue(repo, &issue); err != nil {
			return err
		}

		updatedLine := fmt.Sprintf("%s -> %s\n", strings.Trim(issue.ExtractedLine, "\n"), issue.Url)
		content = strings.Replace(content, issue.ExtractedLine, updatedLine, 1)
	}

	fmt.Printf("\n")
	if err = r.service.Notify(foundIssues); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(content), 0644)
}

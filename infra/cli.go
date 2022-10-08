package infra

import (
	"deni1688/gitissue/domain"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

type Cli struct {
	path    string
	service domain.Service
}

func NewCli(path string, service domain.Service) *Cli {
	return &Cli{path, service}
}

// Issue: #2
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

	issues, err := r.service.ExtractIssues(string(b), r.path)
	repos, err := r.service.ListRepos()
	if err != nil {
		return err
	}

	var current domain.Repo
	base := path.Base(string(origin))
	for _, repo := range *repos {
		if strings.Contains(base, repo.Name) {
			fmt.Println("Found current: ", repo)
			current = repo
			break
		}
	}

	err = r.service.SubmitIssues(current, issues)
	if err != nil {
		return err
	}

	return r.service.Notify(issues)
}

package infra

import (
	"deni1688/gitissue/domain"
	"fmt"
	"os"
	"os/exec"
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

	issues, err := r.service.ExtractIssues(string(b), r.path)
	repo, err := r.service.FindRepoByName(string(origin))

	for _, issue := range *issues {
		fmt.Printf("\n")
		if err = r.service.SubmitIssue(repo, &issue); err != nil {
			return err
		}
	}

	// Issue: Not sure if this is the best way to do it
	if err = r.service.Notify(issues); err != nil {
		return err
	}

	fmt.Printf("\n")
	return nil
}

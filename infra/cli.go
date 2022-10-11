package infra

import (
	"context"
	"deni1688/gogie/internal/issues"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Cli struct {
	service issues.Service
	dry     bool
	ctx     context.Context
}

func NewCli(service issues.Service, dry bool) *Cli {
	return &Cli{service, dry, context.Background()}
}

func (r Cli) Execute(path string) error {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	origin, err := cmd.Output()
	if err != nil {
		return err
	}

	inf, err := os.Stat(path)
	if err != nil {
		return err
	}

	if inf.IsDir() {
		if err = r.ExecuteConcurrently(path); err != nil {
			return err
		}

		return nil
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	b, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	f.Close()

	name := string(origin)
	content := string(b)
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

	// Todo: Optimally the service.Notify should be called after all issues are submitted and files are updated -> https://github.com/deni1688/gogie/issues/29
	fmt.Printf("\n")
	if err = r.service.Notify(foundIssues); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(content), 0644)
}

func (r Cli) ExecuteConcurrently(path string) error {
	var err error
	var files []os.DirEntry

	files, err = os.ReadDir(path)
	if err != nil {
		return err
	}

	g, ctx := errgroup.WithContext(r.ctx)
	g.SetLimit(15)
	for _, dirEntry := range files {
		de := dirEntry
		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return r.Execute(path + "/" + de.Name())
			}
		})
	}

	if err = g.Wait(); err != nil {
		return err
	}

	return nil
}

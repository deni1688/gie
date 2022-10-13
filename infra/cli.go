package infra

import (
	"context"
	"deni1688/gie/internal/issues"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"os"
	"strings"
)

type Cli struct {
	service  issues.Service
	dry      bool
	repoName string
	ctx      context.Context
}

func NewCli(service issues.Service, dry bool, repoName string) *Cli {
	return &Cli{service, dry, repoName, context.Background()}
}

func (r Cli) Execute(path string) error {
	allIssues := make([]issues.Issue, 0)
	issueChan := make(chan issues.Issue)

	if err := r.handlePath(path, &issueChan); err != nil {
		return err
	}

	allIssues = append(allIssues, <-issueChan)

	return r.service.Notify(&allIssues)
}

func (r Cli) handlePath(path string, issueChan *chan issues.Issue) error {
	inf, err := os.Stat(path)
	if err != nil {
		return err
	}

	if inf.IsDir() {
		if err = r.handleDirPath(path, r.repoName, issueChan); err != nil {
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

	if err = f.Close(); err != nil {
		fmt.Println("Error closing file")
	}

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

	for _, issue := range *foundIssues {
		*issueChan <- issue
	}

	<-r.ctx.Done()
	close(*issueChan)

	repo, err := r.service.FindRepoByName(r.repoName)
	for _, issue := range *foundIssues {
		fmt.Printf("\n")
		if err = r.service.SubmitIssue(repo, issue); err != nil {
			return err
		}

		content = strings.Replace(
			content,
			issue.ExtractedLine,
			r.service.GetUpdatedLine(issue), 1)
	}

	return os.WriteFile(path, []byte(content), 0600)
}

func (r Cli) handleDirPath(path, repoName string, issueChan *chan issues.Issue) error {
	var err error
	var files []os.DirEntry

	files, err = os.ReadDir(path)
	if err != nil {
		return err
	}

	g, ctx := errgroup.WithContext(r.ctx)
	g.SetLimit(15)

	for _, de := range files {
		dirEntry := de
		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return r.handlePath(path+"/"+dirEntry.Name(), issueChan)
			}
		})
	}

	if err = g.Wait(); err != nil {
		return err
	}

	return nil
}

package cli

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
	dry      bool
	repoName string
	service  issues.Service
	ctx      context.Context
}

func New(service issues.Service, dry bool, repoName string) *Cli {
	return &Cli{dry, repoName, service, context.Background()}
}

// Issue: Make recursive search through dir thread safe -> closes https://github.com/deni1688/gie/issues/41
func (r Cli) Execute(path string) error {
	fmt.Println("Searching...")

	allIssues := make([]issues.Issue, 0)
	issueCh := make(chan issues.Issue)

	doneCh := make(chan struct{})
	defer close(doneCh)

	var err error
	go func() {
		defer close(issueCh)
		if err = r.handlePath(path, &issueCh); err != nil {
			return
		}
	}()

	go func() {
		for issue := range issueCh {
			allIssues = append(allIssues, issue)
		}

		doneCh <- struct{}{}
	}()
	<-doneCh

	fmt.Printf("Found %d issue(s)\n", len(allIssues))

	return r.service.Notify(&allIssues)
}

// Issue: Ignore private files and dirs (e.g. .git, .idea, etc) -> closes https://github.com/deni1688/gie/issues/39
func (r Cli) handlePath(path string, issueCh *chan issues.Issue) error {
	inf, err := os.Stat(path)
	if err != nil {
		return err
	}

	if inf.IsDir() {
		if err = r.handleDirPath(path, r.repoName, issueCh); err != nil {
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
	found, err := r.service.ExtractIssues(&content, &path)
	if err != nil {
		return err
	}

	if len(*found) < 1 {
		return nil
	}

	repo, err := r.service.FindRepoByName(r.repoName)
	if err != nil {
		return err
	}

	for _, issue := range *found {
		*issueCh <- issue
		fmt.Printf("Found issue=[%s] in path=[%s]\n", issue.Title, path)
		if r.dry {
			continue
		}

		if err = r.service.SubmitIssue(repo, &issue); err != nil {
			return err
		}

		content = strings.Replace(
			content,
			issue.ExtractedLine,
			r.service.GetUpdatedLine(issue), 1)
	}

	if r.dry {
		return nil
	}

	return os.WriteFile(path, []byte(content), 0600)
}

func (r Cli) handleDirPath(path, repoName string, issueCh *chan issues.Issue) error {
	var err error
	var files []os.DirEntry

	files, err = os.ReadDir(path)
	if err != nil {
		return err
	}

	// Issue: Find the best way to do recursive concurrency with errgroup -> closes https://github.com/deni1688/gie/issues/40
	g, ctx := errgroup.WithContext(r.ctx)
	g.SetLimit(15)

	for _, de := range files {
		dirEntry := de

		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return r.handlePath(path+"/"+dirEntry.Name(), issueCh)
			}
		})
	}

	return g.Wait()
}

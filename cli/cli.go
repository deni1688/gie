package cli

import (
	"context"
	"deni1688/gie/core"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"os"
	pth "path"
	"strings"
)

type Cli struct {
	dry      bool
	repoName string
	service  core.Service
	ctx      context.Context
}

func New(service core.Service, dry bool, repoName string) *Cli {
	return &Cli{dry, repoName, service, context.Background()}
}

func (r Cli) Execute(path string) error {
	fmt.Println("Searching...")

	allIssues := make([]core.Issue, 0)
	issueCh := make(chan core.Issue)

	var g errgroup.Group

	g.Go(func() error {
		defer close(issueCh)
		return r.handlePath(path, &issueCh)
	})

	g.Go(func() error {
		for issue := range issueCh {
			allIssues = append(allIssues, issue)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	fmt.Printf("Found %d issue(s)\n", len(allIssues))

	if r.dry {
		return nil
	}

	return r.service.Notify(&allIssues)
}

func (r Cli) handlePath(path string, issueCh *chan core.Issue) error {
	inf, err := os.Stat(path)
	if err != nil {
		return err
	}

	if strings.HasPrefix(pth.Base(path), ".") && path != "." {
		return nil
	}

	if inf.IsDir() {
		if err = r.handleDirPath(path, issueCh); err != nil {
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

func (r Cli) handleDirPath(path string, issueCh *chan core.Issue) error {
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
				return r.handlePath(path+"/"+dirEntry.Name(), issueCh)
			}
		})
	}

	return g.Wait()
}

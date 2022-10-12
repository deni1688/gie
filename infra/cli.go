package infra

import (
	"context"
	"deni1688/gogie/internal/issues"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"os"
	"os/exec"
	"regexp"
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
	repoName, err := getCurrentRepoName()
	if err != nil {
		return err
	}

	foundIssues := make([]issues.Issue, 0)
	if err = r.execute(path, repoName, &foundIssues); err != nil {
		return err
	}

	// Todo: Optimally the service.Notify should be called after all issues are submitted and files are updated -> https://github.com/deni1688/gogie/issues/29
	fmt.Printf("\n")
	return r.service.Notify(&foundIssues)
}

func (r Cli) execute(path, repoName string, issues *[]issues.Issue) error {
	inf, err := os.Stat(path)
	if err != nil {
		return err
	}

	if inf.IsDir() {
		if err = r.executeConcurrently(path, repoName, issues); err != nil {
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

	content := string(b)
	foundIssues, err := r.service.ExtractIssues(&content, &path)
	if len(*foundIssues) < 1 {
		return nil
	}
	*issues = append(*issues, *foundIssues...)

	if r.dry {
		for _, issue := range *foundIssues {
			fmt.Printf("Found issue=[%s] in file=[%s]\n", issue.Title, path)
		}

		return nil
	}

	repo, err := r.service.FindRepoByName(repoName)
	for _, issue := range *foundIssues {
		fmt.Printf("\n")
		if err = r.service.SubmitIssue(repo, &issue); err != nil {
			return err
		}

		updatedLine := fmt.Sprintf("%s -> %s\n", strings.Trim(issue.ExtractedLine, "\n"), issue.Url)
		content = strings.Replace(content, issue.ExtractedLine, updatedLine, 1)
	}

	if err = os.WriteFile(path, []byte(content), 0644); err != nil {
		return err
	}

	return nil
}

func getCurrentRepoName() (string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	res, err := cmd.Output()
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`/(.*)\.git`)
	matches := re.FindStringSubmatch(string(res))
	if matches == nil {
		return "", errors.New("could not find current repo name")
	}

	return matches[1], nil
}

func (r Cli) executeConcurrently(path, repoName string, issues *[]issues.Issue) error {
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
				return r.execute(path+"/"+de.Name(), repoName, issues)
			}
		})
	}

	if err = g.Wait(); err != nil {
		return err
	}

	return nil
}

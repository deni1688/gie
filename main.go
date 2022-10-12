package main

import (
	"deni1688/gogie/infra"
	"deni1688/gogie/internal/issues"
	"flag"
	"fmt"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	setup := flag.Bool("setup", false, "creates a config file")
	configPath := flag.String("config", "", "custom config file path")
	path := flag.String("path", "./issues.txt", "file path to parse issues from")
	prefix := flag.String("prefix", "", "prefix to override config")
	dry := flag.Bool("dry", false, "dry run")
	flag.Parse()

	repoName, err := getCurrentRepoName()
	if err != nil {
		fmt.Println(err)
		return
	}

	if *path == "" {
		fmt.Println("Please provide file path to parse issues from")
		return
	}

	c := new(config)
	if *setup {
		if err = c.Setup(); err != nil {
			fmt.Println(err)
		}

		return
	}

	if err = c.Load(*configPath); err != nil {
		fmt.Println("Error reading config:", err)
		return
	}

	if *prefix != "" {
		c.Prefix = *prefix
	}

	provider, err := newGitProvider(c)
	if err != nil {
		fmt.Println("Error getting provider:", err)
		return
	}

	notifier := infra.NewWebhookNotifier(c.WebHooks, http.DefaultClient)
	service := issues.NewService(provider, notifier, c.Prefix)
	cli := infra.NewCli(service, *dry, repoName)

	if err = cli.Execute(*path); err != nil {
		fmt.Println("Error running cli:", err)
	}
}

func newGitProvider(c *config) (issues.GitProvider, error) {
	switch {
	case strings.Contains(c.Host, "gitlab"):
		return infra.NewGitlab(c.Token, c.Host, c.Query, http.DefaultClient), nil
	case strings.Contains(c.Host, "github"):
		return infra.NewGithub(c.Token, c.Host, c.Query, http.DefaultClient), nil
	default:
		return nil, fmt.Errorf("invalid provider %s", c.Host)
	}
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
		return "", fmt.Errorf("could not find current repo name. Make sure you are in a git repo with a remote origin")
	}

	return matches[1], nil
}

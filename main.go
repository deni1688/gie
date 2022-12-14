package main

import (
	"deni1688/gie/adapters/github"
	"deni1688/gie/adapters/gitlab"
	"deni1688/gie/adapters/webhook"
	"deni1688/gie/cli"
	"deni1688/gie/config"
	"deni1688/gie/core"
	"flag"
	"fmt"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
)

// Supported hosts
const (
	GITHUB = "github"
	GITLAB = "gitlab"
)

var defaultPath = "./issues.txt"

func init() {
	cmd, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}
	defaultPath = strings.Trim(string(cmd), "\n ")
}

func main() {
	setup := flag.Bool("setup", false, "creates a config file")
	configPath := flag.String("config", "", "custom config file path")
	path := flag.String("path", defaultPath, "file path to parse issues from")
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

	c := new(config.Config)
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

	notifier := webhook.New(c.WebHooks, http.DefaultClient)
	service := core.New(provider, notifier, c.Prefix)
	cliApp := cli.New(service, *dry, repoName, c.Exclude)

	if err = cliApp.Execute(*path); err != nil {
		fmt.Println("Error running cli:", err)
	}
}

func newGitProvider(c *config.Config) (core.GitProvider, error) {
	switch {
	case strings.Contains(c.Host, GITLAB):
		return gitlab.New(c.Token, c.Host, c.Query, http.DefaultClient), nil
	case strings.Contains(c.Host, GITHUB):
		return github.New(c.Token, c.Host, c.Query, http.DefaultClient), nil
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

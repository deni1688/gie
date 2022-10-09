package main

import (
	"deni1688/gogie/domain"
	"deni1688/gogie/infra"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"strings"
)

func main() {
	setup := flag.Bool("setup", false, "creates a config file")
	configPath := flag.String("config", "", "custom config file path")
	path := flag.String("path", "./issues.txt", "file path to parse issues from")
	prefix := flag.String("prefix", "", "prefix to override config")
	flag.Parse()

	if *path == "" {
		fmt.Println("Please provide file path to parse issues from")
		return
	}

	c := new(config)
	if *setup {
		if err := c.Setup(); err != nil {
			fmt.Println(err)
		}
		return
	}

	if err := c.Load(*configPath); err != nil {
		fmt.Println("Error reading config:", err)
		return
	}

	if *prefix != "" {
		c.Prefix = *prefix
	}

	p, err := newGitProvider(c)
	if err != nil {
		fmt.Println("Error getting provider:", err)
		return
	}

	n := infra.NewWebhookNotifier(c.WebHooks, http.DefaultClient)
	s := domain.NewService(p, n, c.Prefix)
	cli := infra.NewCli(s)

	// Todo: Make it possible to run the cli.Execute() with dir path -> https://github.com/deni1688/gogie/issues/26
	if err = cli.Execute(*path); err != nil {
		fmt.Println("Error running cli:", err)
	}

	fmt.Println("Done!")
}

func newGitProvider(c *config) (domain.GitProvider, error) {
	switch {
	case strings.Contains(c.Host, "gitlab"):
		return infra.NewGitlab(c.Token, c.Host, c.Query, http.DefaultClient), nil
	case strings.Contains(c.Host, "github"):
		return infra.NewGithub(c.Token, c.Host, c.Query, http.DefaultClient), nil
	default:
		return nil, errors.New("invalid provider " + c.Host)
	}
}

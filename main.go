package main

import (
	"deni1688/gitissue/domain"
	"deni1688/gitissue/infra"
	"errors"
	"flag"
	"fmt"
)

func main() {
	setup := flag.Bool("setup", false, "setup config file")
	path := flag.String("path", "./issues.txt", "please provide file path to parse issues from")
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

	if err := c.Load(); err != nil {
		fmt.Println("Error reading config at $HOME/.config/gitissue.json")
		return
	}

	p, err := getProvider(c)
	if err != nil {
		fmt.Println("Error getting provider: ", err)
		return
	}

	cli := infra.NewCli(c.Prefix, *path, domain.NewService(p))
	if err := cli.Run(); err != nil {
		fmt.Println("Error running cli:", err)
		return
	}

	fmt.Println("Done!")
}

func getProvider(c *config) (domain.Provider, error) {
	switch c.Provider {
	case "gitlab":
		return infra.NewGitlabProvider(c.Token, c.Host, c.Query), nil
	case "github":
		return infra.NewGithubProvider(c.Token, c.Host, c.Query), nil
	default:
		return nil, errors.New("invalid provider " + c.Provider)
	}
}

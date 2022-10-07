package main

import (
	"deni1688/gitissue/domain"
	"deni1688/gitissue/infra"
	"errors"
	"fmt"
)

func main() {
	c := new(config)
	if err := c.Load(); err != nil {
		fmt.Println("Error reading config at $HOME/.config/gitissue.json")
		return
	}

	p, err := getProvider(c)
	if err != nil {
		fmt.Println("Error getting provider: ", err)
		return
	}

	s := domain.NewService(p)
	cli := infra.NewCli(c.Prefix, s)
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

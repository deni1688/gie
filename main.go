package main

import (
	"deni1688/gitissue/domain"
	"deni1688/gitissue/infra"
	"errors"
	"fmt"
	"os"
)

func main() {
	c, err := loadConfig()
	if err != nil {
		fmt.Println("Error reading config at $HOME/.config/gitissue.json")
		os.Exit(1)
	}

	p, err := getProvider(c.Get("provider"), c.Get("token"), c.Get("host"), c.Get("query"))
	if err != nil {
		fmt.Println("Error getting provider: ", err)
		os.Exit(1)
	}

	s := domain.NewService(p)
	cli := infra.NewCli(c.Get("prefix"), s)
	if err := cli.Run(); err != nil {
		fmt.Println("Error running cli:", err)
		os.Exit(1)
	}

	fmt.Println("Done!")
}

func getProvider(provider, token, host, prefix string) (domain.Provider, error) {
	switch provider {
	case "gitlab":
		return infra.NewGitlabProvider(token, host, prefix), nil
	case "github":
		return infra.NewGithubProvider(token, host, prefix), nil
	default:
		return nil, errors.New("invalid provider " + provider)
	}
}

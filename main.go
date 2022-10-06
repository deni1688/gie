package main

import (
	"deni1688/gitissue/domain"
	"deni1688/gitissue/infra"
	"fmt"
	"os"
)

func main() {
	c, err := loadConfig()
	if err != nil {
		fmt.Println("Unable to read the config file at $HOME/.config/gitissue.json")
		os.Exit(1)
	}

	var p domain.Provider
	switch c.Get("provider") {
	case "gitlab":
		p = infra.NewGitlabProvider(
			c.Get("token"),
			c.Get("host"),
			c.Get("query"),
		)
	case "github":
		p = infra.NewGithubProvider(
			c.Get("token"),
			c.Get("host"),
			c.Get("query"),
		)
	default:
		fmt.Println("Invalid provider", c.Get("provider"))
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

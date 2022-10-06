package main

import (
	"deni1688/gitissue/domain"
	"deni1688/gitissue/infra"
	"fmt"
)

func main() {
	provider := infra.NewGitlabProvider("sometoken")
	s := domain.NewService(provider)

	cli := infra.NewCli(s)
	if err := cli.Run(); err != nil {
		fmt.Println("Error running cli:", err)
		return
	}

	fmt.Println("Done")
}

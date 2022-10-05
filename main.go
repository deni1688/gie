package main

import (
	"fmt"
)

func main() {
	provider := NewGitlabProvider("sometoken")
	s := NewService(provider)

	cli := NewCli(s)
	if err := cli.Run(); err != nil {
		fmt.Println("Error running cli:", err)
		return
	}

	fmt.Println("Done")
}

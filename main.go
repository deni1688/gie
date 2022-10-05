package main

import (
	"fmt"
	"os"
)

func main() {
	provider := NewGitlabProvider("sometoken")
	s := NewService(provider)

	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println(`Please provide one or more issues to be created, format: -t="New issue title" -d="This should be fixed" -w=30 -m="sprint/45" -- ..."`)
		return
	}

	issues, err := s.ParseArgs(args)
	if err != nil {
		fmt.Println("Error parsing issues", err)
		return
	}

	if err = s.Execute(issues); err != nil {
		fmt.Println("Error executing issues", err)
		return
	}

	fmt.Println("Done")
}

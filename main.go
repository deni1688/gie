package main

import (
	"deni1688/gitissue/domain"
	"deni1688/gitissue/infra"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type config map[string]string

func (r config) Get(key string) string {
	if val, ok := r[key]; ok {
		return val
	}

	return ""
}

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
		fmt.Println("Provider not implemented yet")
		os.Exit(1)
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

func loadConfig() (config, error) {
	file, err := os.Open(os.Getenv("HOME") + "/.config/gitissue.json")
	if err != nil {
		return nil, err
	}

	var c config
	err = json.NewDecoder(file).Decode(&c)
	if errors.Is(err, &os.PathError{}) {
		return nil, err
	}

	return c, nil
}

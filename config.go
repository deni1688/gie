package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type config struct {
	Host     string   `json:"host"`
	Token    string   `json:"token"`
	Prefix   string   `json:"prefix"`
	Query    string   `json:"query"`
	WebHooks []string `json:"webhooks"`
}

func (r *config) Load(customPath string) error {
	var p string
	if customPath != "" {
		if !strings.Contains(customPath, ".json") {
			return errors.New("only json files are supported")
		}

		p = customPath
	} else {
		p = os.Getenv("HOME") + "/.config/gogie.json"
	}

	file, err := os.Open(p)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(r)
}

func (r *config) Setup() error {
	cp := os.Getenv("HOME") + "/.config/gogie.json"

	if _, err := os.Stat(cp); err == nil {
		return errors.New("config file already exists at " + cp)
	}

	file, err := os.Create(cp)
	if err != nil {
		return err
	}

	fmt.Printf("Creating config file at %s. Navigate to the file and fill in the details.\n", cp)

	return json.NewEncoder(file).Encode(r)
}

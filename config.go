package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type config struct {
	Host     string   `json:"host"`
	Token    string   `json:"token"`
	Prefix   string   `json:"prefix"`
	Query    string   `json:"query"`
	WebHooks []string `json:"webhooks"`
}

func (r *config) Load() error {
	file, err := os.Open(os.Getenv("HOME") + "/.config/gitissue.json")
	if err != nil {
		return err
	}

	return json.NewDecoder(file).Decode(r)
}

func (r *config) Setup() error {
	cp := os.Getenv("HOME") + "/.config/gitissue.json"

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

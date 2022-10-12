package main

import (
	"encoding/json"
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
			return fmt.Errorf("only json files are supported")
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

	if err = json.NewDecoder(file).Decode(r); err != nil {
		return err
	}

	return nil
}

func (r *config) Setup() error {
	defaultConfigPath := os.Getenv("HOME") + "/.config/gogie.json"

	if _, err := os.Stat(defaultConfigPath); err == nil {
		return fmt.Errorf("config file already exists at %s", defaultConfigPath)
	}

	file, err := os.Create(defaultConfigPath)
	if err != nil {
		return err
	}

	fmt.Printf("Creating config file at %s. Navigate to the file and fill in the details.\n", defaultConfigPath)

	if err = json.NewEncoder(file).Encode(r); err != nil {
		return err
	}

	return nil
}

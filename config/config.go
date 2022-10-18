package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Host     string   `json:"host"`
	Token    string   `json:"token"`
	Prefix   string   `json:"prefix"`
	Query    string   `json:"query"`
	WebHooks []string `json:"webhooks"`
	Exclude  []string `json:"exclude"`
}

func (r *Config) Load(customPath string) error {
	var p string
	if customPath != "" {
		if !strings.Contains(customPath, ".json") {
			return fmt.Errorf("only json files are supported")
		}

		p = customPath
	} else {
		p = defaultConfigPath()
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

func (r *Config) Setup() error {
	configPath := defaultConfigPath()
	if _, err := os.Stat(configPath); err == nil {
		return fmt.Errorf("config file already exists at %s", configPath)
	}

	file, err := os.Create(configPath)
	if err != nil {
		return err
	}

	fmt.Printf("Creating config file at %s. Navigate to the file and fill in the details.\n", configPath)

	if err = json.NewEncoder(file).Encode(r); err != nil {
		return err
	}

	return nil
}

func defaultConfigPath() string {
	return os.Getenv("HOME") + "/.config/gie.json"
}

package main

import (
	"encoding/json"
	"os"
)

type config struct {
	Host     string `json:"host"`
	Token    string `json:"token"`
	Prefix   string `json:"prefix"`
	Query    string `json:"query"`
	Provider string `json:"provider"`
}

func (r *config) Load() error {
	file, err := os.Open(os.Getenv("HOME") + "/.config/gitissue.json")
	if err != nil {
		return err
	}

	return json.NewDecoder(file).Decode(r)
}

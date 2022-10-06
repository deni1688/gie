package main

import (
	"encoding/json"
	"errors"
	"os"
)

type config map[string]string

func (r config) Get(key string) string {
	if val, ok := r[key]; ok {
		return val
	}
	return ""
}

func loadConfig() (config, error) {
	file, err := os.Open(os.Getenv("HOME") + "/.config/gitissue.json")
	if errors.Is(err, &os.PathError{}) {
		return nil, err
	}

	var c config
	return c, json.NewDecoder(file).Decode(&c)
}

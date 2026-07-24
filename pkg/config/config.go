package config

import (
	"encoding/json"
	"os"
)

type Redirect struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Status int    `json:"status"`
}

type Rewrite struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Settings struct {
	Redirects []Redirect `json:"redirects"`
	Rewrites  []Rewrite  `json:"rewrites"`
}

func Load(path string) (*Settings, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var s Settings
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

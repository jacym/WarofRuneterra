package main

import "github.com/jacym/WarofRuneterra/server/stat"

type (
	Config struct {
		Port string `env:"PORT" envDefault:"8080"`
	}

	Item struct {
		ID     string       `json:"id"`
		Href   string       `json:"href"`
		Win    bool         `json:"win"`
		Result *stat.Reward `json:"result"`
	}

	State struct {
		// todo: guard with mutex
		items map[string]*Item
	}
)

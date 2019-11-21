package main

type (
	Config struct {
		Port string `env:"PORT" envDefault:"8080"`
	}

	LoRCard struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Region string `json:"region"`
	}
)

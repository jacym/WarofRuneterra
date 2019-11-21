package main

type (
	Config struct {
		Port string `env:"PORT" envDefault:"8080"`
	}

	LoRCard struct {
		ID     string
		Name   string
		Region string
	}
)

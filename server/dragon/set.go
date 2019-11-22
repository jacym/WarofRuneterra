package dragon

import (
	"encoding/json"
	"log"
	"os"

	"github.com/caarlos0/env"
)

func cfg() (*Config, error) {
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func CardSet() (cards []Card) {
	conf, err := cfg()
	if err != nil {
		log.Println(err)
		return
	}

	f, err := os.Open(conf.SetFilePath)
	if err != nil {
		log.Println(err)
		return
	}

	dec := json.NewDecoder(f)

	err = dec.Decode(&cards)
	if err != nil {
		log.Println(err)
	}

	return
}

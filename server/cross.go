package main

import "github.com/jacym/WarofRuneterra/server/dragon"

func findCard(set []dragon.Card, code string) *dragon.Card {
	for _, c := range set {
		if c.Code == code {
			return &c
		}
	}

	return nil
}

func crossCards(set []dragon.Card, references []string) (result []*LoRCard) {
	result = make([]*LoRCard, len(references))

	for _, code := range references {
		if origin := findCard(set, code); origin != nil {
			result = append(result, &LoRCard{
				ID:     origin.Code,
				Name:   origin.Name,
				Region: origin.Region,
			})
		}
	}

	return
}

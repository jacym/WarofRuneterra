package main

import (
	"encoding/json"

	"github.com/jacym/WarofRuneterra/server/dragon"
	"github.com/jacym/WarofRuneterra/server/stat"
)

func findCard(set []dragon.Card, code string) *dragon.Card {
	for _, c := range set {
		if c.Code == code {
			return &c
		}
	}

	return nil
}

func crossCards(set []dragon.Card, references []string) (result []*dragon.Card) {
	result = make([]*dragon.Card, 0)

	for _, code := range references {
		if origin := findCard(set, code); origin != nil {
			result = append(result, origin)
		}
	}

	return
}

func regions(cards []*dragon.Card) stat.Regions {
	keys := make(map[string]bool)

	for _, entry := range cards {
		if _, ok := keys[entry.Region]; !ok {
			keys[entry.Region] = true
		}
	}

	names := make(stat.Regions, 0)

	for k := range keys {
		names = append(names, k)
	}

	return names
}

func save(item *Item) {
	// danger: guard with mutext!
	state.items[item.ID] = item
}

func (i *Item) Encode() string {
	b, err := json.Marshal(i)
	if err != nil {
		return ""
	}

	return string(b)
}

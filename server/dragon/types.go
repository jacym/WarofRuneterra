package dragon

type (
	// Config holds the path to read the set of LoR cards from
	Config struct {
		SetFilePath string `env:"SET_FILEPATH" envDefault:"data/set1-en_us.json"`
	}

	// Card represents a single card from LoR
	Card struct {
		Code        string `json:"cardCode"`
		Name        string `json:"name"`
		Region      string `json:"region"`
		Description string `json:"description"`
	}
)

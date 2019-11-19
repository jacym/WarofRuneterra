package main

import "fmt"

var (
	cfg Config
)

func cards(w http.ResponseWriter, r *http.Request) {
	var cardList []string

	defer r.Body.Close()

	enc := json.NewDecoder(r.Body)
	err := enc.Decode(&cardList)

	if err != nil {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	for _, c := range cardlist {
		log.Printf("card: %s\n", c)
	}
}

func main() {
	fmt.Println("WLoR v0.0.1")

	// config
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	// router
	r := mux.NewRouter()
	r.HandleFunc("/cards", cards)

	srv := &http.Server{
		Handler: r,
		Addr: "127.0.0.1:" + cfg.Port

		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

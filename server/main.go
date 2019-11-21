package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/caarlos0/env"
	"github.com/gorilla/mux"
	"github.com/jacym/WarofRuneterra/server/dragon"
)

var (
	cfg Config
	set []dragon.Card
)

func readCardSet() {
	set = dragon.CardSet()
}

func indexCardSet(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)

	if err := enc.Encode(&set); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

func submitPlayerCards(w http.ResponseWriter, r *http.Request) {
	var cardList []string

	defer r.Body.Close()

	enc := json.NewDecoder(r.Body)
	err := enc.Decode(&cardList)

	if err != nil {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	for _, c := range cardList {
		log.Printf("card: %s\n", c)
	}
}

func cards(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		indexCardSet(w, r)
		break
	case "POST":
		submitPlayerCards(w, r)
		break
	default:
		http.Error(w, http.StatusText(405), 405)
		break
	}
}

func main() {
	log.Println("WLoR v0.0.1")

	// config
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	// card set
	readCardSet()

	// router
	r := mux.NewRouter()
	r.HandleFunc("/cards", cards).Methods("GET", "POST")

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:" + cfg.Port,

		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Listening on %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

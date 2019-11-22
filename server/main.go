package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/caarlos0/env"
	"github.com/gorilla/mux"
	"github.com/jacym/WarofRuneterra/server/dragon"

	packr "github.com/gobuffalo/packr/v2"
)

var (
	cfg       Config
	box       *packr.Box
	templates *template.Template
	set       []dragon.Card
)

func initReadTemplates() (box *packr.Box, tmpl *template.Template) {
	box = packr.New("templates", "./views")
	tmpl = template.New("_all")

	files := []string{
		"partial.html",
		"region.html",
		"show.html",
	}

	for _, t := range files {
		contents, err := box.FindString(t)

		if err != nil {
			log.Println(err)
			continue
		}

		name := strings.TrimSuffix(t, filepath.Ext(t))
		template.Must(tmpl.New(name).Parse(contents))
	}

	return
}

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
	var references []string

	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&references)

	if err != nil {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	enc := json.NewEncoder(w)
	cardList := crossCards(set, references)

	if err := enc.Encode(cardList); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
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

	// templates (html views)
	box, templates = initReadTemplates()

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

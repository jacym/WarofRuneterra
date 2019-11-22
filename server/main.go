package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/caarlos0/env"
	"github.com/gorilla/mux"
	"github.com/jacym/WarofRuneterra/server/dragon"
	"github.com/jacym/WarofRuneterra/server/stat"

	packr "github.com/gobuffalo/packr/v2"
)

var (
	cfg       Config
	box       *packr.Box
	templates *template.Template
	set       []dragon.Card
	state     State
	flaker    *snowflake.Node
)

func init() {
	state = State{
		items: make(map[string]*Item, 7),
	}
}

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
	// todo: add win/lose
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

	// todo: shove this inside stat points.go
	cardRegions := regions(cardList)

	re := stat.WithRegion(cardRegions)

	win := true // todo: fk

	points := re.Calc(win, cardList)

	item := &Item{
		ID:     flaker.Generate().String(),
		Points: points,
		Win:    win,
	}

	save(item) // todo: add/check error

	if err := enc.Encode(&item); err != nil {
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

func show(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	item, ok := state.items[vars["id"]]

	if !ok {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	templates.ExecuteTemplate(w, "show", item)
	return
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

	// snowflake starts falling down
	if node, err := snowflake.NewNode(1); err == nil {
		flaker = node
	} else {
		panic(err) // fk it
	}

	// router
	r := mux.NewRouter()
	r.HandleFunc("/cards", cards).Methods("GET", "POST")
	r.HandleFunc("/view/{id}", show)

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

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
)

var PORT int
var DBNAME string

func getSlugHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	url, err := getURL(slug)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, err.Error())
	}
	w.Header().Add("Location", url)
	w.WriteHeader(http.StatusMovedPermanently)
}

func shortenerHandler(w http.ResponseWriter, r *http.Request) {
	urlToShorten := r.FormValue("url")

	if urlToShorten == "" {
		http.Error(w, "Missing form value 'url'", http.StatusBadRequest)
		return
	}

	if _, err := url.ParseRequestURI(urlToShorten); err != nil {
		http.Error(w, "Invalid url", http.StatusBadRequest)
		return
	}

	res := shortenURL(urlToShorten)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, res)
}

func startServer() {
	r := mux.NewRouter()
	r.HandleFunc("/", shortenerHandler).Methods("POST")
	r.HandleFunc("/{slug}", getSlugHandler).Methods("GET")

	http.Handle("/", r)
	http.ListenAndServe(fmt.Sprintf(":%v", PORT), nil)
}

func setupFlag() {
	portPtr := flag.Int("port", 8080, "port to be used")
	dbnamePtr := flag.String("db", "data.db", "name of database to be used")

	flag.Parse()

	PORT = *portPtr
	DBNAME = *dbnamePtr
	fmt.Println("PORT", PORT)
	fmt.Println("DBNAME", DBNAME)
}

func main() {
	rand.Seed(time.Now().Unix())
	setupFlag()

	openDB(DBNAME)
	defer closeDB()

	startServer()
}

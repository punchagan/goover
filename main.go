package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var DB_PATH string = "/tmp/digest.json"

func main() {

	flag.StringVar(&DB_PATH, "db-path", "/tmp/digest.json", "Full path to the db.")
	flag.Parse()

	router := mux.NewRouter().StrictSlash(true)

	// Web methods
	router.HandleFunc("/", Index)
	router.HandleFunc("/static/{filename}", StaticHandler)

	// API methods
	router.HandleFunc("/view", ArticleIndex)
	router.HandleFunc("/random", RandomArticle)
	router.HandleFunc("/add", AddArticle)

	log.Fatal(http.ListenAndServe(":8080", router))
}

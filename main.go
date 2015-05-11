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
	router.HandleFunc("/view", ArticleIndex)

	log.Fatal(http.ListenAndServe(":8080", router))
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func ArticleIndex(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Invalid query params.")
		return
	}
	tags, ok := r.Form["tags"]
	if ok {
		tags = strings.Split(tags[0], ",")
	}

	entries := GetEntries(tags)
	sendJSONResponse(w, entries)

}

func sendJSONResponse(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"mime"
	"net/http"
	"strings"
	"time"
	"path"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("templates/index.html")
	if err != nil {
		fmt.Fprintln(w, err)
	} else {
		fmt.Fprintln(w, string(data))
	}
}

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]

	data, err := ioutil.ReadFile("static/" + filename)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err)
	} else {
		mime_type := mime.TypeByExtension(path.Ext(filename))
		w.Header().Set("Content-Type", mime_type + "; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, string(data))
	}
}

func ArticleIndex(w http.ResponseWriter, r *http.Request) {
	tags, err := parseTags(w, r)
	if err != nil {
		response := map[string]error{"error": err}
		sendJSONResponse(w, response, http.StatusBadRequest)
	}
	entries := GetEntries(tags)
	sendJSONResponse(w, entries, http.StatusOK)

}

func RandomArticle(w http.ResponseWriter, r *http.Request) {
	tags, err := parseTags(w, r)
	if err != nil {
		response := map[string]error{"error": err}
		sendJSONResponse(w, response, http.StatusBadRequest)
	}
	entry := randomEntry(GetEntries(tags))
	sendJSONResponse(w, entry, http.StatusOK)
}

func sendJSONResponse(w http.ResponseWriter, response interface{}, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func parseTags(w http.ResponseWriter, r *http.Request) (tags []string, err error) {
	err = r.ParseForm()
	if err != nil {
		return nil, err
	}
	tags, ok := r.Form["tags"]
	if ok {
		tags = strings.Split(tags[0], ",")
	}
	return tags, nil
}

func randomEntry(entries []Article) Article {
	rand.Seed(time.Now().Unix())
	return entries[rand.Intn(len(entries))]
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"mime"
	"net/http"
	"path"
	"sort"
	"strings"
	"time"

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
		w.Header().Set("Content-Type", mime_type+"; charset=UTF-8")
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
	sort.Sort(sort.Reverse(entries))
	sendJSONResponse(w, entries, http.StatusOK)

}

func RandomArticle(w http.ResponseWriter, r *http.Request) {
	tags, err := parseTags(w, r)
	if err != nil {
		response := map[string]error{"error": err}
		sendJSONResponse(w, response, http.StatusBadRequest)
	}

	entries := GetEntries(tags)
	if len(entries) == 0 {
		sendJSONResponse(w, nil, http.StatusNotFound)
	} else {
		entry := randomEntry(entries)
		sendJSONResponse(w, entry, http.StatusOK)
	}
}

func AddArticle(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		response := map[string]error{"error": err}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}
	urls, url_ok := r.Form["url"]
	titles, title_ok := r.Form["title"]
	contents, _ := r.Form["content"]
	tags, _ := parseTags(w, r)
	// FIXME: Why don't we allow setting other fields?
	if !(url_ok && title_ok) {
		response := map[string]string{"error": "Missing parameter"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}
	url := urls[0]
	title := titles[0]
	var content string
	if len(contents) == 0 {
		content = ""
	} else {
		content = contents[0]
	}
	article := Article{
		Url:     url,
		Title:   title,
		Content: content,
		Tags:    tags,
		Added:   CustomTime(time.Now()),
	}

	err = AddEntry(article)
	if err != nil {
		response := map[string]error{"error": err}
		sendJSONResponse(w, response, http.StatusNotFound)
	} else {
		response := map[string]bool{"success": true}
		sendJSONResponse(w, response, http.StatusNotFound)
	}
}

func EditArticle(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		response := map[string]error{"error": err}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}
	urls, url_ok := r.Form["url"]
	tags, err := parseTags(w, r)
	// FIXME: Why don't we allow editing other fields?
	if !url_ok || err != nil {
		response := map[string]string{"error": "Missing parameter"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}
	entries := GetEntryMap()

	for _, url := range urls {
		article, ok := entries[url]
		if ok {
			article = article.AddRemoveTags(tags)
			// Save the updated article to the db.
			err = AddEntry(article) //fixme: AddEntry -> UpdateEntry
		}
	}

	if err != nil {
		response := map[string]error{"error": err}
		sendJSONResponse(w, response, http.StatusNotFound)
	} else {
		response := map[string]bool{"success": true}
		sendJSONResponse(w, response, http.StatusNotFound)
	}
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
		tags = strings.Split(strings.Trim(tags[0], ", "), ",")
	}

	return tags, nil
}

func randomEntry(entries Articles) Article {
	rand.Seed(time.Now().Unix())
	return entries[rand.Intn(len(entries))]
}

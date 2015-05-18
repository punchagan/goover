package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"time"
)

func GetEntries(tags []string) (entries Articles) {
	entry_map := GetEntryMap()

	for _, article := range entry_map {
		if article.HasTags(tags) {
			entries = append(entries, article)
		}
	}
	return entries
}

func GetEntryMap() (entries map[string]Article) {
	data := readDb(DB_PATH)

	if data == nil {
		return nil
	}

	entries = make(map[string]Article)

	for id, value := range data.(map[string]interface{}) {
		switch value.(type) {

		case map[string]interface{}:

			v := value.(map[string]interface{})
			// FIXME: marshal() -> unmarshal() to change type,
			// kinda sucks!
			b, _ := json.Marshal(v)
			var article Article
			json.Unmarshal(b, &article)
			article.Id = id

			entries[id] = article

		}
	}
	return entries
}

func AddEntry(article Article) (err error) {
	// FIXME: DB needs to be locked.
	data := readDb(DB_PATH)
	if data != nil {
		D := data.(map[string]interface{})
		D[article.Id] = article
		var json_data []byte
		json_data, err = json.Marshal(D)
		ioutil.WriteFile(DB_PATH, json_data, 0755)
	} else {
		err = errors.New("could not add article.")
	}
	return err
}

func readDb(path string) interface{} {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Opening db failed with error: %s", err)
		return nil
	}

	var parsed interface{}
	err = json.Unmarshal(data, &parsed)
	if err != nil {
		log.Printf("Corrupt db, parsing failed with error: %s", err)
		return nil
	}

	return parsed
}

func getId() string {
	return getRandomString(20)
}

func getRandomString(n int) string {
	choices := "0123456789abcdef"
	rand.Seed(time.Now().Unix())
	r := ""
	for i := 0; i < n; i++ {
		r = r + string(choices[rand.Intn(len(choices))])
	}
	return r
}

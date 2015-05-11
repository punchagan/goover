package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

func GetEntries(tags []string) (entries []Article) {
	data := readDb()

	if data == nil {
		return nil
	}

	for _, value := range data.(map[string]interface{}) {
		switch value.(type) {

		case map[string]interface{}:

			v := value.(map[string]interface{})
			// FIXME: marshal() -> unmarshal() to change type,
			// kinda sucks!
			b, _ := json.Marshal(v)
			var article Article
			json.Unmarshal(b, &article)

			if article.HasTags(tags) {
				entries = append(entries, article)
			}

		}
	}
	return entries
}

func readDb() interface{} {
	path := "/tmp/digest.json"
	fmt.Println(path)

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

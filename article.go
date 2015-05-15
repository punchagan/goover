package main

import (
	"encoding/json"
	"time"
)

type CustomTime time.Time

const TIME_FORMAT = "2006-01-02T15:04:05-0700"

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(ct))
}

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	b = b[1 : len(b)-1]
	t, err := time.Parse(TIME_FORMAT, string(b))
	if err != nil {
		return err
	}
	*ct = CustomTime(t)
	return nil
}

func (ct CustomTime) String() (s string) {
	return time.Time(ct).String()
}

type Article struct {
	Url       string     `json:"url"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	Blog      string     `json:"blog"`
	Author    string     `json:"author"`
	Published CustomTime `json:"date_published"`
	Added     CustomTime `json:"date_added"`
	Tags      []string   `json:"tags"`
}

func (article Article) hasTag(tag string) bool {
	// Convert list of tags to map for fast lookup
	atags := make(map[string]bool)
	for _, t := range article.Tags {
		atags[t] = true
	}
	return atags[tag]
}

func (article Article) HasTags(tags []string) bool {
	for _, tag := range tags {
		if tag[0] != '!' && !article.hasTag(tag) {
			return false
		}
		if tag[0] == '!' && article.hasTag(tag[1:]) {
			return false
		}
	}
	return true
}

// Articles List

type Articles []Article

func (articles Articles) Len() int {
	return len(articles)
}

func (articles Articles) Less(i, j int) bool {
	return time.Time(articles[i].Published).Before(time.Time(articles[j].Published))
}

func (articles Articles) Swap(i, j int) {
	articles[i], articles[j] = articles[j], articles[i]
}

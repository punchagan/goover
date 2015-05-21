package main

import (
	"time"
)

type CustomTime time.Time

const TIME_FORMAT = "2006-01-02T15:04:05-0700"

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return []byte("\"" + ct.String() + "\""), nil
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
	return time.Time(ct).Format(TIME_FORMAT)
}

type Article struct {
	Id        string     `json:"id"`
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
		if len(tag) == 0 {
			continue
		}
		if tag[0] != '!' && !article.hasTag(tag) {
			return false
		}
		if tag[0] == '!' && article.hasTag(tag[1:]) {
			return false
		}
	}
	return true
}

func (article Article) AddRemoveTags(tags []string) Article {
	for _, tag := range tags {
		if len(tag) == 0 {
			continue
		}
		if tag[0] != '!' {
			if !article.hasTag(tag) {
				article.Tags = append(article.Tags, tag)
			}
		} else {
			if article.hasTag(tag[1:]) {
				article.Tags = RemoveTag(article.Tags, tag[1:])
			}
		}
	}
	return article
}

func RemoveTag(tags []string, tag string) []string {
	T := make(map[string]bool)

	// make a map out of the list
	for _, t := range tags {
		T[t] = true
	}

	delete(T, tag)
	tags = make([]string, len(T))

	for t, _ := range T {
		tags = append(tags, t)
	}

	return tags

}

// Articles List

type Articles []Article

func (articles Articles) Len() int {
	return len(articles)
}

func (articles Articles) Less(i, j int) bool {
	return time.Time(articles[i].Added).Before(time.Time(articles[j].Added))
}

func (articles Articles) Swap(i, j int) {
	articles[i], articles[j] = articles[j], articles[i]
}

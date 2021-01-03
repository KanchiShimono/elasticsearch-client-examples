package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v7"
)

type Fields struct {
	Directors       []string `json:"directors"`
	ReleaseDate     string   `json:"release_date"`
	Rating          float64  `json:"rating"`
	Genres          []string `json:"genres"`
	ImageURL        string   `json:"image_url"`
	Plot            string   `json:"plot"`
	Title           string   `json:"title"`
	Rank            int64    `json:"rank"`
	RunningTimeSecs int64    `json:"running_time_secs"`
	Actors          []string `json:"actors"`
	Year            int64    `json:"year"`
}

type Source struct {
	Fields `json:"fields"`
}

type SearchResult struct {
	Index  string  `json:"_index"`
	ID     string  `json:"_id"`
	Type   string  `json:"_type"`
	Score  float64 `json:"_score"`
	Source `json:"_source"`
}

type Client struct {
	ES *elasticsearch.Client
}

func NewClient(es *elasticsearch.Client) Client {
	return Client{es}
}

func (c Client) Search(q string, s int64) ([]SearchResult, error) {
	buf := bytes.Buffer{}
	query := map[string]interface{}{
		"size": s,
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  q,
				"fields": []string{"fields.title^4", "fields.plot^2", "fields.actors", "fields.directors"},
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}
	res, err := c.ES.Search(
		c.ES.Search.WithContext(context.Background()),
		c.ES.Search.WithIndex("movies"),
		c.ES.Search.WithBody(&buf),
		c.ES.Search.WithTrackTotalHits(true),
		c.ES.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("[%s] %s: %s",
			res.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"],
		)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, err
	}

	results := []SearchResult{}
	for _, hit := range resp["hits"].(map[string]interface{})["hits"].([]interface{}) {
		r := SearchResult{}

		jsonByte, err := json.Marshal(hit)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(jsonByte, &r); err != nil {
			return nil, err
		}

		results = append(results, r)
	}
	return results, nil
}

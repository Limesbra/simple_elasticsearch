package elastic

import (
	"encoding/json"
	"ex01/model"
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v8"
)

type Client struct {
	*elasticsearch.Client
}

func NewClient() *Client {
	return &Client{connectToElasticsearch()}
}

func connectToElasticsearch() *elasticsearch.Client {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	return es
}

func (c *Client) getTotal() int {

	res, err := c.Count(
		c.Count.WithIndex("places"),
	)

	if err != nil {
		return 0
	}

	defer res.Body.Close()

	var response struct {
		Count int `json:"count"`
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		log.Fatal(err)
	}

	total := response.Count

	return total
}

func (c *Client) MakeBigger() {

	req := esapi.IndicesPutSettingsRequest{
		Index: []string{"strings"},
		Body: strings.NewReader(`
		{
			"index": {
				"max_result_window": 15000
			}
		}
	`),
	}

	res, err := c.Indices.PutSettings(req.Body)

	if err != nil {
		fmt.Printf("Error setting max_result_window: %v\n", err)
		return
	}
	defer res.Body.Close()

	// Проверьте ответ
	if res.IsError() {
		log.Fatal(res.String())
	}
}

func (c *Client) makeRequest(limit, offset int, places []model.Place) []model.Place {
	res, err := c.Search(
		c.Search.WithIndex("places"),
		c.Search.WithSize(limit),
		c.Search.WithFrom(offset),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatal(res.String())
	}

	var response struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
			Hits []struct {
				ID     string      `json:"_id"`
				Source model.Place `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		log.Fatal(err)
	}

	for _, hit := range response.Hits.Hits {
		places = append(places, hit.Source)
	}

	return places
}

func (c *Client) GetPlaces(limit, offset int) ([]model.Place, int, error) {
	if offset < 0 {
		return nil, 0, fmt.Errorf("invalid 'page' value: 'foo'")
	}
	places := make([]model.Place, 0)
	if limit == 0 {
		return places, 0, nil
	}

	c.MakeBigger()

	total := c.getTotal()

	if offset > total {
		return nil, 0, fmt.Errorf("invalid 'page' value: 'foo'")
	}

	for offset < total {
		places = c.makeRequest(limit, offset, places)
		offset += limit
	}

	return places, total, nil
}

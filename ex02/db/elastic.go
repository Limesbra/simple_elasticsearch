package elastic

import (
	"encoding/json"
	"ex02/model"
	"fmt"
	"log"

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

func (c *Client) makeRequest(limit, offset int, places []model.Place) []model.Place {
	res, err := c.Search(
		c.Search.WithIndex("places"),
		c.Search.WithSize(limit),
		c.Search.WithFrom(offset),
	)
	if err != nil {
		return nil
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

	total := c.getTotal()

	if offset > total {
		return nil, 0, fmt.Errorf("invalid 'page' value: 'foo'")
	}

	places = c.makeRequest(limit, offset, places)

	return places, total, nil
}

func AddMetaInfo(places []model.Place, total int, page int, totalPages int) model.Api {
	var meta model.Api

	meta.Name = "places"
	meta.Total = total
	meta.Places = places
	meta.Prev_page = page - 1
	meta.Next_page = page + 1
	meta.Last_page = totalPages

	return meta
}

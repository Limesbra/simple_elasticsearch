package elastic

import (
	"encoding/json"
	"ex03/model"
	"fmt"
	"log"
	"strings"

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

func (c *Client) makeRequest(lat, lon float64, places []model.Place) []model.Place {
	query := fmt.Sprintf(`{
		"query": {
			"match_all": {}
		},
		"sort": [
			{
				"_geo_distance": {
					"location": {
						"lat": %f,
						"lon": %f
					},
					"order": "asc",
					"unit": "km",
					"mode": "min",
					"distance_type": "arc",
					"ignore_unmapped": true
				}
			}
		],
		"size": 3
	}`, lat, lon)

	res, err := c.Search(
		c.Search.WithIndex("places"),
		c.Search.WithBody(strings.NewReader(query)),
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

func (c *Client) GetPlaces(lat, lon float64) ([]model.Place, error) {
	if lat < 0 || lon < 0 {
		return nil, fmt.Errorf("geopoint can not be less than 0")
	}
	places := make([]model.Place, 0)

	places = c.makeRequest(lat, lon, places)

	return places, nil
}

func AddInfo(places []model.Place) model.Api {
	var meta model.Api

	meta.Name = "places"
	meta.Places = places

	return meta
}

package main

// why dont have "place" - https://www.elastic.co/guide/en/elasticsearch/reference/7.17/removal-of-types.html#_why_are_mapping_types_being_removed

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

type Geopoint struct {
	Long float64 `json:"lon"`
	Lat  float64 `json:"lat"`
}

// type Geopoint struct {
// 	Coordinate []float64
// }

type Restaurant struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Address  string   `json:"address"`
	Phone    string   `json:"phone"`
	Location Geopoint `json:"location"`
}

var mapping = `{
    "mappings": {
		"properties": {
			"name": {
				"type": "text"
			},
			"address": {
				"type": "text"
			},
			"phone": {
				"type": "text"
			},
			"location": {
				"type": "geo_point"
			}
		}
	}
}`

func LoadPlacesFromFile(ctx context.Context) context.Context {
	var mu sync.Mutex
	var wg sync.WaitGroup
	// Open the CSV file
	file, err := os.Open("/Users/limesbra/bootcamp/go/Go_Day03-1/materials/data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)
	reader.Comma = '\t'

	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Create a slice of Restaurant structs
	restaurants := make([]Restaurant, 0, len(records))
	for count, record := range records {
		if count == 0 {
			continue
		}

		recordCopy := record

		wg.Add(1)
		go func() {
			defer wg.Done()
			// coordinate := make([]float64, 2)
			// coordinate[0] = convertToFloat(recordCopy[4])
			// coordinate[1] = convertToFloat(recordCopy[5])
			restaurant := Restaurant{
				ID:      int(convertToFloat(recordCopy[0])),
				Name:    recordCopy[1],
				Address: recordCopy[2],
				Phone:   recordCopy[3],
				Location: Geopoint{
					Long: convertToFloat(recordCopy[4]),
					Lat:  convertToFloat(recordCopy[5]),
				},
			}
			mu.Lock()
			restaurants = append(restaurants, restaurant)
			mu.Unlock()
		}()
	}

	wg.Wait()

	fmt.Printf("✅ Data loaded from CSV file: %d records\n", len(restaurants))
	return context.WithValue(ctx, "restaurants", restaurants)
}

func convertToFloat(str string) float64 {
	result, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Println("Error converting string to float:", err)
		return 0
	}
	return result
}

func ConnectWithElasticSearch(ctx context.Context) context.Context {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("✅ Connected to ElasticSearch")
	return context.WithValue(ctx, "es", es)
}

func IndexData(ctx context.Context) {
	restaurants := ctx.Value("restaurants").([]Restaurant)
	es := ctx.Value("es").(*elasticsearch.Client)

	_, err := es.Indices.Delete([]string{"places"})
	if err != nil {
		log.Fatal(err)
	}

	_, err = es.Indices.Create(
		"places",
		es.Indices.Create.WithBody(strings.NewReader(mapping)),
	)

	if err != nil {
		log.Fatal(err)
	}

	bulkIndexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client:     es,
		Index:      "places",
		NumWorkers: 5,
	})

	if err != nil {
		log.Fatal(err)
	}

	for docuumentID, restaurant := range restaurants {

		data, err := json.Marshal(restaurant)

		if err != nil {
			log.Fatal(err)
		}

		err = bulkIndexer.Add(
			ctx,
			esutil.BulkIndexerItem{
				Action:     "index",
				DocumentID: strconv.Itoa(docuumentID),
				Body:       strings.NewReader(string(data)),
			},
		)

		if err != nil {
			log.Fatal(err)
		}
	}

	bulkIndexer.Close(ctx)
	biStats := bulkIndexer.Stats()

	fmt.Printf("✅ Data indexed on Elasticsearch: %d \n", biStats.NumIndexed)

}

func main() {

	ctx := context.Background()

	ctx = LoadPlacesFromFile(ctx)
	ctx = ConnectWithElasticSearch(ctx)

	IndexData(ctx)
}

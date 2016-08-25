package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	bigquery "google.golang.org/api/bigquery/v2"
)

func main() {
	f, err := os.Open("personsDataSchema.json")
	if err != nil {
		log.Fatalf("unable to open file. %v", err)
	}

	var value []*bigquery.TableFieldSchema

	err = json.NewDecoder(f).Decode(&value)
	if err != nil {
		log.Fatalf("unable to unmarshal json. %v", err)
	}

	fmt.Println(len(value))
}

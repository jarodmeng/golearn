package main

import (
	"fmt"
	"log"
)

func main() {
	const projectId = "google.com:jarodmeng"
	const datasetId = "testing"

	bqService, err := createBqService(defaultSecretFile, defaultBqTokenFile, defaultBqScope)
	if err != nil {
		log.Fatalf("Unable to create BigQuery service: %v", err)
	}

	r, err := bqService.Tables.List(projectId, datasetId).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve table list. %v", err)
	}

	fmt.Println(len(r.Tables))

	for _, t := range r.Tables {
		fmt.Println(t.Id)
	}
}

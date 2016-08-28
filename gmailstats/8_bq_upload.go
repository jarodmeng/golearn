package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	bigquery "google.golang.org/api/bigquery/v2"
)

func main() {
	const projectId = "google.com:jarodmeng"
	const datasetId = "testing"
	const tableId = "gmailheader3"

	const dataFile = "messages.json"
	const schemaFile = "messagesSchema.json"

	bqService, err := createBqService(defaultSecretFile, defaultBqTokenFile, defaultBqScope)
	if err != nil {
		log.Fatalf("Unable to create BigQuery service: %v", err)
	}

	sf, err := os.Open(schemaFile)
	defer sf.Close()
	if err != nil {
		log.Fatalf("unable to open file. %v", err)
	}

	var value []*bigquery.TableFieldSchema

	err = json.NewDecoder(sf).Decode(&value)
	if err != nil {
		log.Fatalf("unable to unmarshal json. %v", err)
	}

	schema := &bigquery.TableSchema{
		Fields: value,
	}

	configuration := &bigquery.JobConfigurationLoad{
		CreateDisposition: "CREATE_IF_NEEDED",
		WriteDisposition:  "WRITE_APPEND",
		DestinationTable: &bigquery.TableReference{
			ProjectId: projectId,
			DatasetId: datasetId,
			TableId:   tableId,
		},
		SourceFormat: "NEWLINE_DELIMITED_JSON",
		Schema:       schema,
	}

	job := &bigquery.Job{
		Configuration: &bigquery.JobConfiguration{
			Load: configuration,
		},
	}

	df, err := os.Open(dataFile)
	defer df.Close()
	if err != nil {
		log.Fatalf("Unable to open the file.", err)
	}

	j, err := bqService.Jobs.Insert(projectId, job).Media(df).Do()

	fmt.Println(j.Status.State)
}

package main

import (
	"github.com/jarodmeng/googleauth"
	bigquery "google.golang.org/api/bigquery/v2"
)

func createBqService(secretFile string, tokenFile string, scope string) (*bigquery.Service, error) {
	bqClient, err := googleauth.CreateClientFromFile(secretFile, tokenFile, scope)
	if err != nil {
		return nil, err
	}

	bqService, err := bigquery.New(bqClient)
	if err != nil {
		return nil, err
	}

	return bqService, nil
}

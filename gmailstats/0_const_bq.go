// token cache file and scope information for BigQuery API calls
package main

import bigquery "google.golang.org/api/bigquery/v2"

const (
	defaultBqTokenFile = "bq_token.json"
	defaultBqScope     = bigquery.BigqueryScope
)

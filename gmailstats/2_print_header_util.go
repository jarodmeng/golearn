package main

import (
	"errors"
	"fmt"
	"log"

	gmail "google.golang.org/api/gmail/v1"
)

func extractHeader(headers []*gmail.MessagePartHeader, part_header string) (string, error) {
	for _, h := range headers {
		if h.Name == part_header {
			return h.Value, nil
		}
	}
	return "", errors.New("Cannot find the part header")
}

func printHeader(mr *gmail.Message, part_header string) {
	headerValue, err := extractHeader(mr.Payload.Headers, part_header)
	if err != nil {
		log.Fatalf("Unable to extract the part header. %v", err)
	}
	fmt.Printf("%s\n", headerValue)
}

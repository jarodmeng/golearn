package main

import (
	"fmt"
	"log"

	gmail "google.golang.org/api/gmail/v1"
)

func main() {
	const numMessages = 500

	// Stage 0: set up plumbing
	gmailService, err := createGmailService(defaultSecretFile, defaultGmailTokenFile, defaultGmailScope)
	if err != nil {
		log.Fatalf("Unable to create Gmail service: %v", err)
	}

	// Stage 1: get message ids
	// this api call only returns message ids and thread ids
	messages := make([]*gmail.Message, 0)
	user := defaultGmailUser
	call := gmailService.Users.Messages.List(user).MaxResults(numMessages).Q("after:2016/01/01")
	r0, err := call.Do()
	if err != nil {
		log.Fatalf("Unable to retrieve messages. %v", err)
	}
	nextToken := r0.NextPageToken
	for nextToken != "" {
		r1, err := call.PageToken(nextToken).Do()
		if err != nil {
			log.Fatalf("Unable to retrieve messages. %v", err)
		}
		fmt.Printf("profiling %d messages.\n", len(r1.Messages))
		messages = append(messages, r1.Messages...)
		nextToken = r1.NextPageToken
	}
	fmt.Printf("total # of messages: %d.\n", len(messages))
}

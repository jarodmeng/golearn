package main

import (
	"fmt"
	"log"
)

func main() {
	const numMessages = 100

	gmailService, err := createGmailService(defaultSecretFile, defaultGmailTokenFile, defaultGmailScope)
	if err != nil {
		log.Fatalf("Unable to create Gmail service: %v", err)
	}
	user := defaultGmailUser

	r, err := gmailService.Users.Messages.List(user).MaxResults(numMessages).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve messages. %v", err)
	}

	fmt.Printf("Total # of messages received: %d\n", len(r.Messages))

	partHeader := "From"
	for _, m := range r.Messages {
		mr, _ := gmailService.Users.Messages.Get(user, m.Id).Do()
		printHeader(mr, partHeader)
	}

	fmt.Println("Finished all IDs.\n")
}

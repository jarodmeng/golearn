package main

import (
	"fmt"
	"log"
)

func main() {
	gmailService, err := createGmailService(defaultSecretFile, defaultGmailTokenFile, defaultGmailScope)
	if err != nil {
		log.Fatalf("Unable to create Gmail service: %v", err)
	}
	user := defaultGmailUser

	r, err := gmailService.Users.Messages.List(user).MaxResults(1).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve messages. %v", err)
	}

	for _, m := range r.Messages {
		mr, _ := gmailService.Users.Messages.Get(user, m.Id).Do()
		for _, h := range mr.Payload.Headers {
			fmt.Printf("%s: %s.\n", h.Name, h.Value)
		}
	}
}

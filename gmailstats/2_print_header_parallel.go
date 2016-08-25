package main

import (
	"fmt"
	"log"
	"sync"
)

func main() {
	const numMessages = 100
	const numWorkers = 8

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

	ids := make(chan string, len(r.Messages))
	for _, m := range r.Messages {
		ids <- m.Id
	}
	close(ids)

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	partHeader := "Mailing-list"
	for i := 0; i < numWorkers; i++ {
		go func() {
			for id := range ids {
				mr, _ := gmailService.Users.Messages.Get(user, id).Do()
				printHeader(mr, partHeader)
			}
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println("Finished all IDs.\n")
}

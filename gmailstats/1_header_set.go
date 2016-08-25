// Get a number of messages, get their headers and tabulate the occurances of
// each particular header.
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

	var headerCounter = struct {
		sync.RWMutex
		header map[string]int
	}{header: make(map[string]int)}

	for i := 0; i < numWorkers; i++ {
		go func() {
			for id := range ids {
				mr, _ := gmailService.Users.Messages.Get(user, id).Do()
				for _, h := range mr.Payload.Headers {
					headerCounter.Lock()
					headerCounter.header[h.Name]++
					headerCounter.Unlock()
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()

	fmt.Println("Finished all IDs.\n")
	for k, v := range headerCounter.header {
		fmt.Println(k, ": ", v)
	}
}

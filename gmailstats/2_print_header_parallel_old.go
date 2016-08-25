package main

import (
	"fmt"
	"log"

	gmail "google.golang.org/api/gmail/v1"
)

type Worker struct {
	ID    int
	Count int
	Work  chan string
	Quit  chan bool
}

func (w *Worker) Start(srv *gmail.Service, userId string) {
	go func() {
		partHeader := "From"
		for {
			select {
			case work := <-w.Work:
				mr, _ := srv.Users.Messages.Get(userId, work).Do()
				w.Count++
				fmt.Printf("Worker %d has done work %d times.\n", w.ID, w.Count)
				printHeader(mr, partHeader)
			case <-w.Quit:
				fmt.Printf("Worker %d has been stopped.\n", w.ID)
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	go func() {
		fmt.Printf("Stopping worker %d.\n", w.ID)
		w.Quit <- true
	}()
}

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

	in := make(chan string)
	go func() {
		for _, m := range r.Messages {
			in <- m.Id
		}
		close(in)
	}()

	for i := 0; i < numWorkers; i++ {
		worker := Worker{
			ID:    i + 1,
			Count: 0,
			Work:  in,
			Quit:  make(chan bool),
		}
		worker.Start(gmailService, user)
	}

	var input string
	fmt.Scanln(&input)
}

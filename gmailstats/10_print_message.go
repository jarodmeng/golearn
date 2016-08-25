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

	r, err := gmailService.Users.Messages.List(user).MaxResults(50).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve messages. %v", err)
	}

	for _, m := range r.Messages {
		mr, _ := gmailService.Users.Messages.Get(user, m.Id).Do()
		// fmt.Printf("Payload.MimeType: %s.\n", mr.Payload.MimeType)
		// fmt.Printf("HistoryId: %d.\n", mr.HistoryId)
		// fmt.Printf("Id: %s.\n", mr.Id)
		// fmt.Printf("InternalDate: %d.\n", mr.InternalDate)
		// fmt.Printf("Raw: %s.\n", mr.Raw)
		// fmt.Printf("SizeEstimate: %d.\n", mr.SizeEstimate)
		// fmt.Printf("Snippet: %s.\n", mr.Snippet)
		// fmt.Printf("ThreadId: %s.\n", mr.ThreadId)
		// fmt.Printf("Payload.Body.Data: %s.\n", mr.Payload.Body.Data)
		// fmt.Printf("Payload.Body.Size: %d.\n", mr.Payload.Body.Size)
		fmt.Printf("Size of Payload.Parts: %d.\n", len(mr.Payload.Parts))
		for _, pp := range mr.Payload.Parts {
			fmt.Printf("MimeType: %s.\n", pp.MimeType)
		}
		// pp0 := mr.Payload.Parts[0].Parts[0]
		// fmt.Printf("Payload.Parts[0].Parts[0].Body.Data: %s.\n", pp0.Body.Data)
		// data0, _ := base64.StdEncoding.DecodeString(pp0.Body.Data)
		// fmt.Printf("Payload.Parts[0].Parts[0].Body.Data: %s.\n", data0)
		// fmt.Printf("Payload.Parts[0].Parts[0].Body.size: %d.\n", pp0.Body.Size)
		// pp1 := mr.Payload.Parts[0].Parts[1]
		// fmt.Printf("Payload.Parts[0].Parts[0].Body.Data: %s.\n", pp1.Body.Data)
		// data1, _ := base64.StdEncoding.DecodeString(pp1.Body.Data)
		// fmt.Printf("Payload.Parts[0].Parts[0].Body.Data: %s.\n", data1)
		// fmt.Printf("Payload.Parts[0].Parts[0].Body.size: %d.\n", pp1.Body.Size)
	}
}

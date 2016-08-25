package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

// a struct to hold a gmail message's meta data
type headerMetadata struct {
	MessageId   string   `json:"messageid,omitempty"`
	ThreadId    string   `json:"threadid,omitempty"`
	Timestamp   int64    `json:"timestamp,omitempty"`
	Time        string   `json:"time,omitempty"`
	FromEmail   string   `json:"fromemail,omitempty"`
	ToEmails    []string `json:"toemails,omitempty"`
	CcEmails    []string `json:"ccemails,omitempty"`
	BccEmails   []string `json:"bccemails,omitempty"`
	MailingList string   `json:"mailinglist,omitempty"`
	Subject     string   `json:"subject,omitempty"`
}

func main() {
	f, _ := os.OpenFile("output.json", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	defer f.Close()
	dec := json.NewDecoder(f)
	ms := make([]headerMetadata, 0)
	for {
		var m headerMetadata
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		ms = append(ms, m)
	}
	fmt.Println(len(ms))
}

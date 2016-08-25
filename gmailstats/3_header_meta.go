package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	gmail "google.golang.org/api/gmail/v1"
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

// merge a slice of channels into one
// when all input channels are exhausted, close the output channel
func merge(ms []chan headerMetadata) chan headerMetadata {
	out := make(chan headerMetadata)
	var wg sync.WaitGroup

	output := func(c chan headerMetadata) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	wg.Add(len(ms))

	for _, c := range ms {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func matchEmail(s string) string {
	const emailRawRegex = `(?:^|(?:^.+ <))([[:alnum:]-+@.]+)(?:$|(?:>$))`
	emailRegex, _ := regexp.Compile(emailRawRegex)
	email := emailRegex.FindStringSubmatch(s)
	if len(email) < 2 {
		return s
	}
	return email[1]
}

func matchAllEmails(s string) []string {
	out := make([]string, 0)
	emailStrings := strings.Split(s, ", ")
	for _, e := range emailStrings {
		email := matchEmail(e)
		out = append(out, email)
	}
	return out
}

func matchMailingList(s string) string {
	const mlRawRegex = `^list ([[:alnum:]-@.]+);.+`
	mlRegex := regexp.MustCompile(mlRawRegex)
	ml := mlRegex.FindStringSubmatch(s)
	if len(ml) < 2 {
		return s
	}
	return ml[1]
}

func formatMicroseconds(ms int64) string {
	const timeLayout = "Mon 2006/01/02 15:04:05 PST"
	s := ms / 1000
	ns := int64(math.Mod(float64(ms), 1000)) * 1e6
	t := time.Unix(s, ns)
	return t.Format(timeLayout)
}

func main() {
	const numWorkers = 8
	const numMessages = 500

	// Stage 0: set up plumbing
	gmailService, err := createGmailService(defaultSecretFile, defaultGmailTokenFile, defaultGmailScope)
	if err != nil {
		log.Fatalf("Unable to create Gmail service: %v", err)
	}

	f, err := os.OpenFile("output.json", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("unable to read/create json file. %v", err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	existingIds := make(map[string]bool)
	for {
		var m headerMetadata
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		existingIds[m.MessageId] = true
	}
	fmt.Println(len(existingIds))

	// Stage 1: get message ids
	// this api call only returns message ids and thread ids
	messages := make([]*gmail.Message, 0)
	user := defaultGmailUser
	call := gmailService.Users.Messages.List(user).MaxResults(numMessages).Q("-is:chat after:2015/01/01")
	r0, err := call.Do()
	if err != nil {
		log.Fatalf("Unable to retrieve messages. %v", err)
	}
	nextToken := r0.NextPageToken
	pageCounter := 0
	for nextToken != "" {
		pageCounter++
		r1, err := call.PageToken(nextToken).Do()
		if err != nil {
			log.Fatalf("Unable to retrieve messages. %v", err)
		}
		fmt.Printf("%d: profiling %d messages.\n", pageCounter, len(r1.Messages))
		messages = append(messages, r1.Messages...)
		nextToken = r1.NextPageToken
	}
	fmt.Printf("total # of messages: %d.\n", len(messages))

	// Stage 2: feed message ids into a channel
	// start a goroutine that takes each gmail message and puts its message Id
	// to a string channel
	// close the string channel when the messages are exhausted
	ids := make(chan string)
	go func() {
		for _, m := range messages {
			if _, ok := existingIds[m.Id]; !ok {
				ids <- m.Id
			}
		}
		close(ids)
	}()

	// Stage 3: get detailed information on each message and send to a channel
	// this is the main work section
	// multiple workers are started, each with a separate goroutine. those
	// goroutines all use the same input channel that's created and being
	// populated using the goroutine above.
	// within each goroutine, the message id is used to get the detailed
	// information of the message; the detail information is then used to
	// populate a headerMetadata struct; lastly the struct is sent to the
	// output channel.
	metaSlices := make([]chan headerMetadata, numWorkers)
	for i := 0; i < numWorkers; i++ {
		metaSlices[i] = make(chan headerMetadata)
	}
	for i := 0; i < numWorkers; i++ {
		go func(i int) {
			out := metaSlices[i]
			defer close(out)
			for id := range ids {
				mr, _ := gmailService.Users.Messages.Get(user, id).Do()
				meta := headerMetadata{
					MessageId: mr.Id,
					ThreadId:  mr.ThreadId,
					Timestamp: mr.InternalDate / 1000,
					Time:      formatMicroseconds(mr.InternalDate),
				}
				for _, h := range mr.Payload.Headers {
					switch h.Name {
					case "From":
						meta.FromEmail = matchEmail(h.Value)
					case "To":
						meta.ToEmails = matchAllEmails(h.Value)
					case "Cc":
						meta.CcEmails = matchAllEmails(h.Value)
					case "Bcc":
						meta.BccEmails = matchAllEmails(h.Value)
					case "Mailing-list":
						meta.MailingList = matchMailingList(h.Value)
					case "Subject":
						meta.Subject = h.Value
					}
				}
				out <- meta
			}
		}(i)
	}

	// Stage 5: merge channels from stage 4 into one
	metas := merge(metaSlices)

	// Stage 6: write each message's meta info into a json row
	counter := 0
	for m := range metas {
		err = json.NewEncoder(f).Encode(m)
		if err != nil {
			log.Fatalf("unable to write json file. %v", err)
		}
		counter++
		fmt.Printf("processed message id %s.\n", m.MessageId)
	}

	fmt.Printf("Finished %d IDs.\n", counter)
}

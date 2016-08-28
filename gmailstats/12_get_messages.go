package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/jarodmeng/gmailstats"
)

var (
	mr = flag.Int64("mr", 100, "A positive integer")
)

func main() {
	flag.Parse()
	gs := gmailstats.New()
	gs, err := gs.ListMessages().MaxResults(*mr).Do()
	if err != nil {
		log.Fatalf("Error: %v.\n", err)
	}
	fmt.Println(len(gs.MessageIds))

	gs = gs.GetMessages().Write().Do()
	fmt.Println(len(gs.Messages))
	s, _ := json.Marshal(gs.Messages[50])
	fmt.Println(string(s))
}

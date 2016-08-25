package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/jarodmeng/gmailstats"
)

var (
	mr      = flag.Int64("mr", 100, "A positive integer")
	verbose = flag.Bool("verbose", false, "")
)

func main() {
	flag.Parse()
	gs := gmailstats.New()
	gs, err := gs.ListMessages().MaxResults(*mr).Q("from:mfeng@google.com").Do()
	if err != nil {
		log.Fatalf("Error: %v.\n", err)
	}
	fmt.Println(len(gs.MessageIds))
	if *verbose {
		for _, mid := range gs.MessageIds {
			fmt.Println(mid.MessageId)
		}
	}
}

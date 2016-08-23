package main

import "sync"

// merge any number of sending channel into one sending channel
func merge(cs ...<-chan int) <-chan int {
	// the merged sending channel returned by merge()
	out := make(chan int)
	// set up a wait group, that upon invoking the Add() method would wait for
	// operations on each input channel to be completed
	var wg sync.WaitGroup

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan int) {
		// the range clause relies on the closing of the input channel to indicate
		// the end of the loop. otherwise, it would hang in wait for the input
		// channel c to send yet another value
		for n := range c {
			out <- n
		}
		// once all values from the input channel is copied to the output channel,
		// call the Done() method to let the wait group know that the operation on
		// this channel is done.
		wg.Done()
	}

	// number of wait hooks equals to the number of input channels
	wg.Add(len(cs))
	// start a copy operation for each input channel in a goroutine.
	// this loop returns to main() immediately once those goroutines are started
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

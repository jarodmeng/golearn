package main

import (
	"fmt"
	"sync"
)

func merge(done <-chan struct{}, cs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	// this function is applied onto each input channel
	output := func(c <-chan int) {
		// the for loop starts when it receives a value from the input c channel
		// the received value is stored in n
		// the receiving action triggers the input channel (c1 or c2) to receive
		// a value from the in channel which in turn receives a value from the
		// argument list of the gen() function
		// the second output() call starts another for loop that ask to receive
		// a value from the other input channel (c2 or c1); the input channel then
		// asks the in channel to send another value which the in channel is going
		// to ask from the other value in the argument list of the gen() call
		// since gen() only has 2 arguments, the second output() call would
		// exhaust the argument list which leads to the goroutine started by gen()
		// to exit and cause the in channel to close
		// when the in channel closes, both c1 and c2 channel close since the
		// goroutines that start them depend on the in channel
		for n := range c {
			// the operation now waits for either of the two options
			// 1) out channel is ready to receive and n is sent to the out channel
			// 2) done channel receives some input
			// when either case is satisfied, the select statement completes and
			// it moves onto the next value in the for loop
			select {
			case out <- n:
			case <-done:
			}
		}
		fmt.Println("output() goroutine exits, calling Done().")
		// when for loop completes (when the input c channel is closed), notify
		// the wait group that this goroutine is completed
		wg.Done()
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		fmt.Println("merge() exits, closing channel.")
		close(out)
	}()

	return out
}

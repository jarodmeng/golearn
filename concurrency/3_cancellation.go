package main

import "fmt"

func main() {
	in := gen(2, 3)

	// Distribute the sq work across two goroutines that both read from in.
	c1 := sq(in)
	c2 := sq(in)

	// Consume the first value from output.
	done := make(chan struct{})
	out := merge(done, c1, c2)
	// this statment notifies the out channel that it's ready to receive
	// the signal triggers either c1 or c2 channel to send out a value that's
	// retrieved from the in channel
	fmt.Println(<-out) // 4 or 9
	// the retrieval of out value causes the sending channel (c1 or c2) to
	// complete the hanging select{} statement and go to the next for loop value;
	// since the in channel is now closed (see comments in the merge() function),
	// the for loop exits and Done() is called on the goroutine.

	// for the other sending channel (c2 or c1), its n value still hasn't been
	// sent out yet, so it hangs in the select{} statement and Done() is not
	// called.

	// send an empty struct (so it doesn't cost any memory space) to the done
	// channel
	// this causes the other sending channel (c2 or c1) to complete its hanging
	// select{} statement and move on to the next for loop value
	// since the channel is already closed, the for loop exits and calls Done()
	// at this stage, both Done() are called and Wait() completes; so merge()
	// exits and closes the out channel (not that it's important since no
	// subsequent for loop depends on the out channel, but it's generally a good
	// practice to close a channel when it's not in use)
	done <- struct{}{}

	// by now the main() goroutine exits and the only channel that's still open
	// is the done channel

	// for debugging purpose only
	// var input string
	// fmt.Scanln(&input)
}

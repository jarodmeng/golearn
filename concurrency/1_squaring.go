package main

import "fmt"

// this is the most simplistic version of the pipeline
func main() {
	// c is a sending channel that sends the argument list of gen() one integer
	// at a time
	// gen() starts a goroutine and returns the channel, so this statement returns
	// to main() immediately
	c := gen(2, 3)
	// out is a sending channel that sends out the square of each element that c
	// sends out
	// sq() also starts a goroutine and returns the channel, so this statement
	// returns to main() immediately too
	out := sq(c)

	// Consume the output.
	// this statement triggers a send from out, which in turn triggers a send from
	// c which then receives from the argument list of gen()
	fmt.Println(<-out) // 4
	// this statment triggers the second round of pipeline
	fmt.Println(<-out) // 9

	// since main() goroutine has exhausted the argument list of gen(), the
	// goroutine in gen() closes the c channel and exits the goroutine. the
	// closing of the c channel in turn closes the out channel and completes the
	// goroutine started by sq()
}

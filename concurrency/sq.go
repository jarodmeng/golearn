package main

import "fmt"

func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
		fmt.Println("sq() goroutine exits, closing channel.")
	}()
	return out
}

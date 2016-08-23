package main

import "fmt"

func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
		fmt.Println("gen() goroutine exits, closing channel.")
	}()
	return out
}

package main

import "fmt"

var WorkerQueue chan chan WorkRequest

func StartDispatcher(nworkers int) {
	// First, initialize the channel
	WorkerQueue = make(chan chan WorkRequest, nworkers)

	// Create all workers
	for i := 0; i < nworkers; i++ {
		fmt.Println("Starting worker", i+1)
		worker := NewWorker(i+1, WorkerQueue)
		worker.Start()
	}

	go func() {
		for {
			select {
			// get a WorkRequest from WorkQueue, or wait for a WorkRequest to be
			// available
			case work := <-WorkQueue:
				fmt.Println("Received work request")
				go func() {
					// get a worker, or wait for a worker to be available
					worker := <-WorkerQueue

					fmt.Println("Dispatching work request")
					// send the WorkRequest to the worker when it's ready to receive
					worker <- work
				}()
			}
		}
	}()
}

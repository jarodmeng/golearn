package main

import (
	"fmt"
	"time"
)

type Worker struct {
	ID          int
	Work        chan WorkRequest
	WorkerQueue chan chan WorkRequest
	QuitChan    chan bool
}

// NewWorker creates, and returns a new Worker object. Its only argument
// is a channel to which the work can add itself whenever it has done its
// work.
func NewWorker(id int, workerQueue chan chan WorkRequest) Worker {
	// Create, and return the worker
	worker := Worker{
		ID:          id,
		Work:        make(chan WorkRequest),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool),
	}
	return worker
}

// This function "starts" the worker by starting a goroutine, that is
// an infinite "for-select" loop
func (w *Worker) Start() {
	go func() {
		for {
			// Add itself into the worker queue
			// A worker is basically a channel of WorkRequest
			w.WorkerQueue <- w.Work

			select {
			// when the worker receives a WorkRequest
			case work := <-w.Work:
				// Receive a work request
				fmt.Printf("worker%d: Received work request, delaying for %f seconds\n", w.ID, work.Delay.Seconds())

				time.Sleep(work.Delay)
				fmt.Printf("worker:%d: Hello, %s!\n", w.ID, work.Name)
			// when the worker receives a signal to stop
			case <-w.QuitChan:
				// We have been asked to stop
				fmt.Printf("worker%d stopping\n", w.ID)
				// return the function, hence stop the worker
				return
			}
		}
	}()
}

// Stop tells the worker to stop listening to work requests
// Note that the worker will only stop *after* it has finished its work
func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

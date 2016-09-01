package main

import (
	"errors"
	"fmt"
	"math"
)

/*
worker methods
*/
func (w *worker) processWork(in *job) (out *result, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	if math.Mod(float64(in.num), 6.0) == 0.0 {
		panic(errors.New("It's a 6!"))
	}
	out = &result{num: in.num * in.num}

	return
}

func (w *worker) work() {
	// when a worker finishes its jobs, sign itself off
	defer func() {
		w.manager.finish.Done()
	}()

	// add the worker itself to the manager's team
	w.manager.team <- w
	// when the worker receives a job via its job queue channel; the workers knows
	// when to exit when its job queue is closed by the manager
	for j := range w.jobQueue {
		r, err := w.processWork(j)
		// when there's an error, trigger manager to pause and record the failed job
		// otherwise, send the result to the output channel for collating
		if err != nil {
			fmt.Printf("Unable to process job: %v.\n", err)
			w.manager.triggerPause()
			w.manager.fail.Lock()
			w.manager.fail.jobs = append(w.manager.fail.jobs, j)
			w.manager.fail.Unlock()
		} else {
			w.manager.output <- r
		}

		// regardless of the state, when the worker is free again, add itself back
		// to the team so that it can take the next job available
		w.manager.team <- w
	}
}

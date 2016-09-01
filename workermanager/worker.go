package main

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"time"
)

/*
job and result types correspond to input and output of the workflow
*/
type job struct {
	num int
}

type result struct {
	num int
}

type failJobs struct {
	sync.RWMutex
	jobs []*job
}

/*
manager and worker types definition
*/
type manager struct {
	nWorkers int             // how many workers the manager deploys
	team     chan *worker    // workers organized as a channel
	finish   *sync.WaitGroup // worker coordination group
	input    chan *job       // input as a channel of job
	pause    chan bool       // pause indicates when manager should pause
	output   chan *result    // output as a channel of result
	result   []*result       // a slice of result to be returned
	fail     failJobs        // a slice of failed jobs
	retry    int             // retry limits
}

type worker struct {
	id       int       // worker id
	manager  *manager  // worker is aware of its manager
	jobQueue chan *job // each worker has its own job queue
}

// NewManager initiates a new manager instance and returns the pointer to it
func NewManager(nWorkers int) *manager {
	m := &manager{
		nWorkers: nWorkers,
		finish:   &sync.WaitGroup{},
		retry:    0,
	}

	return m
}

func (m *manager) setZeroes() {
	m.team = make(chan *worker, m.nWorkers)
	m.input = make(chan *job)
	m.pause = make(chan bool, m.nWorkers)
	m.output = make(chan *result)
	m.fail = failJobs{jobs: make([]*job, 0)}
}

// SetRetry sets a manager's retry limit
func (m *manager) SetRetry(n int) *manager {
	m.retry = n
	return m
}

// SetNumWorkers sets a manager's number of workers
func (m *manager) SetNumWorkers(nWorkers int) *manager {
	m.nWorkers = nWorkers

	return m
}

func (m *manager) convertInput(in []*job) {
	// feed input job into input channel one by one
	for _, j := range in {
		m.input <- j
	}
	// close input channel when input slice is exhausted
	close(m.input)
}

func (m *manager) newWorker(id int) *worker {
	// initialize a worker instance that's aware of the manager
	w := &worker{
		id:       id,
		manager:  m,
		jobQueue: make(chan *job),
	}
	// register the worker with the manager
	w.manager.finish.Add(1)

	return w
}

func (m *manager) populateWorkers() {
	// create a pool of workers and set them to work
	for i := 0; i < m.nWorkers; i++ {
		w := m.newWorker(i)
		go w.work()
	}
}

func (m *manager) pauseWork(n int64) {
	fmt.Printf("Manager is sleeping for %d seconds.\n", n)
	time.Sleep(time.Duration(n) * time.Second)
}

func (m *manager) triggerPause() {
	m.pause <- true
}

func (m *manager) processInput() {
forloop:
	for {
		select {
		// when the manager receives a job through the input channel
		case job, more := <-m.input:
			// jump out of outer for loop when the input channel is closed
			if !more {
				break forloop
			}
			// get a worker from the team channel
			w := <-m.team
			// send the received job to the avaiblable worker's job queue
			w.jobQueue <- job
		// when the manager receives a pause request
		case <-m.pause:
			m.pauseWork(3)
		}
	}

	// when the input channel is exhausted, "lay off" workers by closing their
	// job queues. this for loop finishes when the team channel is closed
	// when the finish wait group is reduced to zero (in waitWorkers()).
	for w := range m.team {
		close(w.jobQueue)
	}

	// when all workers are "laid off", close the output channel so that the
	// collate job knows when to finish
	close(m.output)
}

func (m *manager) waitWorkers() {
	m.finish.Wait()
	close(m.team)
}

func (m *manager) collateOutput() {
	for r := range m.output {
		m.result = append(m.result, r)
	}
}

func (m *manager) Manage(in []*job) *manager {
	m.setZeroes()         // set all intermediate channels/slices to zero values
	go m.convertInput(in) // convert input slice to input channel for the manager
	m.populateWorkers()   // hire workers and set them to work
	go m.processInput()   // start processing jobs from the input channel
	go m.waitWorkers()    // start a goroutine to wait for workers to finish
	m.collateOutput()     // start collating results from the output channel

	// when there's failed jobs and retry is still poisitive, retry
	if len(m.fail.jobs) != 0 && m.retry > 0 {
		m.retry--
		fmt.Println("Retrying...")
		m = m.Manage(m.fail.jobs)
	}

	return m
}

// Report returns output and failed jobs from a manager
func (m *manager) Report() ([]*result, []*job) {
	return m.result, m.fail.jobs
}

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

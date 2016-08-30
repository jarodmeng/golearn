package main

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"time"
)

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

type manager struct {
	nWorkers int
	team     chan *worker
	finish   *sync.WaitGroup
	input    chan *job
	pause    chan bool
	output   chan *result
	result   []*result
	fail     failJobs
	retry    int
}

type worker struct {
	id       int
	manager  *manager
	jobQueue chan *job
}

func newManager(n int) *manager {
	m := &manager{
		nWorkers: n,
		team:     make(chan *worker, n),
		finish:   &sync.WaitGroup{},
		input:    make(chan *job),
		pause:    make(chan bool, n),
		output:   make(chan *result),
		fail:     failJobs{jobs: make([]*job, 0)},
		retry:    0,
	}

	return m
}

func (m *manager) reset() {
	m.team = make(chan *worker, m.nWorkers)
	m.input = make(chan *job)
	m.pause = make(chan bool, m.nWorkers)
	m.output = make(chan *result)
	m.fail = failJobs{jobs: make([]*job, 0)}
}

func (m *manager) setRetry(n int) *manager {
	m.retry = n
	return m
}

func (m *manager) convertInput(in []*job) {
	m.input = make(chan *job)
	for _, j := range in {
		m.input <- j
	}
	close(m.input)
}

func (m *manager) newWorker(id int) *worker {
	w := &worker{
		id:       id,
		manager:  m,
		jobQueue: make(chan *job),
	}
	w.manager.finish.Add(1)
	return w
}

func (m *manager) populateWorkers() {
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
	// for job, more := <-m.input; more; job, more = <-m.input {
	// for job := range m.input {
forloop:
	for {
		select {
		case job, more := <-m.input:
			if !more {
				break forloop
			}
			w := <-m.team
			w.jobQueue <- job
		case <-m.pause:
			m.pauseWork(3)
		}
	}

	for w := range m.team {
		close(w.jobQueue)
	}

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

func (m *manager) manage(in []*job) *manager {
	m.reset()
	go m.convertInput(in)
	m.populateWorkers()
	go m.processInput()
	go m.waitWorkers()
	m.collateOutput()

	if len(m.fail.jobs) != 0 && m.retry > 0 {
		m.retry--
		fmt.Println("Retrying...")
		m = m.manage(m.fail.jobs)
	}

	return m
}

func (m *manager) report() ([]*result, []*job) {
	return m.result, m.fail.jobs
}

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
	defer func() {
		w.manager.finish.Done()
	}()

	w.manager.team <- w
	for j := range w.jobQueue {
		r, err := w.processWork(j)
		if err != nil {
			fmt.Printf("Unable to process job: %v.\n", err)
			w.manager.triggerPause()
			w.manager.fail.Lock()
			w.manager.fail.jobs = append(w.manager.fail.jobs, j)
			w.manager.fail.Unlock()
		} else {
			w.manager.output <- r
		}

		w.manager.team <- w
	}
}

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	input := make([]*job, 0)
	for _, n := range nums {
		input = append(input, &job{num: n})
	}

	out, fail := newManager(2).setRetry(4).manage(input).report()

	fmt.Printf("Output is a %d-element slice.\n", len(out))
	for _, r := range out {
		fmt.Println(r.num)
	}

	fmt.Printf("Fail is a %d-element slice.\n", len(fail))
	for _, f := range fail {
		fmt.Println(f.num)
	}
}

package main

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

type work struct {
	num int
}

type result struct {
	num int
}

type manager struct {
	nWorkers int
	team     chan *worker
	finish   *sync.WaitGroup
	input    <-chan *work
	output   chan *result
}

type worker struct {
	id        int
	manager   *manager
	workQueue chan *work
}

func newManager(n int, in chan *work) *manager {
	if n <= 0 {
		log.Fatalln("# of workers cannot be less than 1.")
	}

	m := &manager{
		nWorkers: n,
		team:     make(chan *worker, n),
		finish:   &sync.WaitGroup{},
		input:    in,
		output:   make(chan *result),
	}

	return m
}

func (m *manager) newWorker(id int) *worker {
	w := &worker{
		id:        id,
		manager:   m,
		workQueue: make(chan *work),
	}

	return w
}

func processWork(in *work) (out *result, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	if in.num == 5 {
		panic(errors.New("It's a 5!"))
	}
	out = &result{num: in.num * in.num}
	return
}

func (m *manager) work() chan *result {
	for i := 0; i < m.nWorkers; i++ {
		w := m.newWorker(i)
		w.work()
		w.manager.finish.Add(1)
	}

	go func() {
		for work := range m.input {
			w := <-m.team
			w.workQueue <- work
		}

		for w := range m.team {
			close(w.workQueue)
		}

		// fmt.Println("Closing output.")
		close(m.output)
	}()

	go func() {
		m.finish.Wait()
		// fmt.Println("Closing team.")
		close(m.team)
	}()

	return m.output
}

func (w *worker) work() {
	go func() {
		w.manager.team <- w
		for wk := range w.workQueue {
			r, err := processWork(wk)
			if err != nil {
				fmt.Printf("Unable to process work: %v.\n", err)
			} else {
				w.manager.output <- r
			}
			w.manager.team <- w
		}
		// fmt.Printf("Worker %d signs off.\n", w.id)
		w.manager.finish.Done()
	}()
}

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	input := make(chan *work)

	go func() {
		for _, n := range nums {
			input <- &work{num: n}
		}
		// fmt.Println("Closing input.")
		close(input)
	}()

	m := newManager(1, input)
	out := m.work()

	for r := range out {
		fmt.Println(r.num)
	}
}
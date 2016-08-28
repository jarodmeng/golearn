package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"sync"
)

type work struct {
	num int
}

type result struct {
	num int
}

type fail struct {
	sync.RWMutex
	fail []*work
}

type manager struct {
	nWorkers int
	team     chan *worker
	finish   *sync.WaitGroup
	input    <-chan *work
	output   chan *result
	fail     fail
}

type worker struct {
	id        int
	manager   *manager
	workQueue chan *work
}

func newManager(n int, in []*work) *manager {
	if n <= 0 {
		log.Fatalln("# of workers cannot be less than 1.")
	}

	c := make(chan *work)
	go func() {
		for _, wk := range in {
			c <- wk
		}
		close(c)
	}()

	m := &manager{
		nWorkers: n,
		team:     make(chan *worker, n),
		finish:   &sync.WaitGroup{},
		input:    c,
		output:   make(chan *result),
		fail:     fail{fail: make([]*work, 0)},
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
	if math.Mod(float64(in.num), 5.0) == 0.0 {
		panic(errors.New("It's a 5!"))
	}
	out = &result{num: in.num * in.num}
	return
}

func (m *manager) work() ([]*result, []*work) {
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

	out := make([]*result, 0)
	for r := range m.output {
		out = append(out, r)
	}

	return out, m.fail.fail
}

func (w *worker) work() {
	go func() {
		w.manager.team <- w
		for wk := range w.workQueue {
			r, err := processWork(wk)
			if err != nil {
				fmt.Printf("Unable to process work: %v.\n", err)
				w.manager.fail.Lock()
				w.manager.fail.fail = append(w.manager.fail.fail, wk)
				w.manager.fail.Unlock()
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
	input := make([]*work, 0)
	for _, n := range nums {
		input = append(input, &work{num: n})
	}

	m := newManager(2, input)
	out, fail := m.work()
	if len(fail) > 0 {
		for i := 0; i < 3; i++ {
			if len(fail) == 0 {
				break
			}

			m = newManager(2, fail)
			a, b := m.work()
			if len(a) > 0 {
				out = append(out, a...)
			}
			fail = b
		}
	}

	fmt.Printf("Output is a %d-element slice.\n", len(out))
	for _, r := range out {
		fmt.Println(r.num)
	}

	fmt.Printf("Fail is a %d-element slice.\n", len(fail))
	for _, f := range fail {
		fmt.Println(f.num)
	}
}

package main

import "sync"

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

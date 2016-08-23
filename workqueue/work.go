package main

import "time"

type WorkRequest struct {
	Name  string        // name to print
	Delay time.Duration // duration to wait for
}

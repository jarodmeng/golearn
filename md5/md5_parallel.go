package main

import (
	"crypto/md5"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type result struct {
	path string
	sum  [md5.Size]byte
	err  error
}

func sumFiles(done <-chan struct{}, root string) (<-chan result, <-chan error) {
	// For each regular file, start a goroutine that sums the file and sends
	// the result on c.  Send the result of the walk on errc.
	c := make(chan result)
	errc := make(chan error, 1)

	mainFunc := func() {
		var wg sync.WaitGroup

		// this function is applied to every file in the root directory
		walkFunc := func(path string, info os.FileInfo, err error) error {
			// if walk fails for any file in the root directory, terminate the walk
			// and return the error
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}

			wg.Add(1)

			// this function dictates how each file should be processed and synced
			// with external channels (c to store results and done to indicate exit)
			fileSum := func() {
				data, err := ioutil.ReadFile(path)
				select {
				// when c is ready to receive, send it the result for the file
				// in question
				case c <- result{path, md5.Sum(data), err}:
				// alternatively, if done is closed, exit the select immediately and
				// don't wait for c to be ready
				case <-done:
				}
				wg.Done()
			}
			go fileSum()

			// Abort the walk if done is closed.
			select {
			// if done is closed, cancel the walk and return an error
			case <-done:
				return errors.New("walk canceled")
			// if done is not ready, fall through to return nil
			default:
				return nil
			}
		}

		err := filepath.Walk(root, walkFunc)

		// Walk has returned, so all calls to wg.Add are done.  Start a
		// goroutine to close c once all the sends are done.
		go func() {
			wg.Wait()
			close(c)
		}()

		// No select needed here, since errc is buffered.
		errc <- err
	}
	go mainFunc()

	return c, errc
}

func MD5All(root string) (map[string][md5.Size]byte, error) {
	// MD5All closes the done channel when it returns; it may do so before
	// receiving all the values from c and errc.
	done := make(chan struct{})
	defer close(done)

	c, errc := sumFiles(done, root)

	m := make(map[string][md5.Size]byte)
	for r := range c {
		if r.err != nil {
			return nil, r.err
		}
		m[r.path] = r.sum
	}
	if err := <-errc; err != nil {
		return nil, err
	}
	return m, nil
}

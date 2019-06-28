// Package worker implements the function to handle different jobs concurrently.
package worker

import (
	"fmt"
	"sync"
)

// JobFunc is a func that makes some job and returns an error if it's failed.
type JobFunc func() error

// Handle handles concurrently the jobFuncs at the same time, depending on provided maxConcurrentJobsCount.
// maxErrCount shows how many errors are acceptable. If there are more errors than the maxErrCount then
// the function will return an error.
func Handle(jobFuncs []JobFunc, maxConcurrentJobsCount, maxErrCount int) error {
	// To make sure that all jobs will be finished..
	var wg sync.WaitGroup
	wg.Add(len(jobFuncs))

	// errCount shows how many jobs already have returned errors. If errCount will be same or greater
	// than the maxErrCount then the handling will be stopped.
	errCount := 0
	// There is a need to lock changes of the errCount between goroutines.
	var m sync.Mutex

	// To control maximum concurrent jobs handling at the same time.
	queue := make(chan JobFunc, maxConcurrentJobsCount)
	go func(<-chan JobFunc) {
		for job := range queue {
			err := job()
			if err != nil {
				m.Lock()
				errCount++
				m.Unlock()
			}
			wg.Done()
		}
	}(queue)

	for _, jobFunc := range jobFuncs {
		queue <- jobFunc
	}

	// The goroutine controls the jobs errors.
	errorQuit := make(chan struct{})
	go func() {
		for {
			m.Lock()
			if errCount > maxErrCount {
				errorQuit <- struct{}{}
			}
			m.Unlock()
		}
	}()

	// The goroutine waits until all jobs will be handled.
	done := make(chan struct{})
	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

	select {
	// The handling is failed due to too many errors.
	case <-errorQuit:
		return getMaxCountReachedError(maxErrCount)
	// The handling is finished.
	case <-done:
		if errCount > maxErrCount {
			return getMaxCountReachedError(maxErrCount)
		}
		return nil
	}
}

func getMaxCountReachedError(maxErrCount int) error {
	return fmt.Errorf("error stack overflow. Max acceptable count is %d", maxErrCount)
}

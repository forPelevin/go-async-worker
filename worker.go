// Package worker implements the function to handle different jobs concurrently.
package worker

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// JobFunc is a func that makes some job and returns an error if it's failed.
type JobFunc func() error

// Handle handles concurrently the jobFuncs at the same time, depending on provided maxConcurrentJobsCount.
// maxErrCount shows how many errors are acceptable. If there are more errors than the maxErrCount then
// the function will return an error.
func Handle(jobFuncs []JobFunc, maxConcurrentJobsCount, maxErrCount int) error {
	// To make sure that all jobs will be finished..
	var wg sync.WaitGroup

	// Channel to quit when the max error count is reached.
	errorQuit := make(chan struct{})
	// errCount shows how many jobs already have returned errors. If errCount will be same or greater
	// than the maxErrCount then the handling will be stopped.
	var errCount int32

	// To control maximum concurrent jobs handling at the same time.
	queue := make(chan struct{}, maxConcurrentJobsCount)
	for _, jobFunc := range jobFuncs {
		queue <- struct{}{}
		wg.Add(1)
		go func(job JobFunc, queue <-chan struct{}) {
			if atomic.LoadInt32(&errCount) > int32(maxErrCount) {
				errorQuit <- struct{}{}
				return
			}

			err := job()
			if err != nil {
				atomic.AddInt32(&errCount, 1)
			}

			<-queue
			wg.Done()
		}(jobFunc, queue)
	}

	// The goroutine waits until all jobs will be handled.
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(queue)
		done <- struct{}{}
	}()

	select {
	case <-errorQuit:
		// The handling is failed due to too many errors.
		return getMaxCountReachedError(maxErrCount)
	case <-done:
		// The handling is finished.
		return nil
	}
}

func getMaxCountReachedError(maxErrCount int) error {
	return fmt.Errorf("error stack overflow. Max acceptable count is %d", maxErrCount)
}

package worker

import (
	"errors"
	"log"
	"sync"
	"testing"
)

func TestHandle(t *testing.T) {
	var m sync.Mutex
	i := 0
	incrFunc := func() error {
		m.Lock()
		i++
		m.Unlock()

		return nil
	}
	errFunc := func() error {
		return errors.New("something went wrong in the job func")
	}

	testCases := []struct {
		description            string
		jobs                   []JobFunc
		maxConcurrentJobsCount int
		maxErrCount            int
		expectsError           bool
		expectedResult         int
	}{
		{
			"simple case with 3 jobs without errors",
			[]JobFunc{incrFunc, incrFunc, incrFunc},
			1,
			0,
			false,
			3,
		},
		{
			"case when 2 of the 3 jobs are failed",
			[]JobFunc{errFunc, errFunc, incrFunc},
			3,
			2,
			false,
			1,
		},
		{
			"case when 1 of the 5 are failed, and it was unacceptable",
			[]JobFunc{incrFunc, errFunc, incrFunc, errFunc, errFunc},
			1,
			1,
			true,
			2,
		},
	}

	for _, tc := range testCases {
		// Reset the counter for correct test case handling.
		i = 0

		err := Handle(tc.jobs, tc.maxConcurrentJobsCount, tc.maxErrCount)
		if err != nil && tc.expectsError == false {
			log.Fatalf("Does not expect that the case `%s` will have an error. Error: %v", tc.description, err)
		}

		if err == nil && tc.expectsError == true {
			log.Fatalf("Expected that the case `%s` will have an error, but it not", tc.description)
		}

		if i != tc.expectedResult {
			log.Fatalf(
				"Expected that the case `%s` will have counter %d, but got %d",
				tc.description,
				tc.expectedResult,
				i,
			)
		}
	}
}

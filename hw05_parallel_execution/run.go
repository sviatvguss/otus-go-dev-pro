package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	wg := &sync.WaitGroup{}
	var errorsCount int32
	var ignoreErrors bool

	if len(tasks) == 0 {
		return nil
	}

	if m <= 0 {
		ignoreErrors = true
	}

	tasksCh := make(chan Task, len(tasks))
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range tasksCh {
				if atomic.LoadInt32(&errorsCount) <= int32(m) {
					err := task()
					if !ignoreErrors && err != nil {
						atomic.AddInt32(&errorsCount, 1)
					}
				}
			}
		}()
	}

	for _, task := range tasks {
		tasksCh <- task
	}
	close(tasksCh)

	wg.Wait()
	if !ignoreErrors && errorsCount > int32(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}

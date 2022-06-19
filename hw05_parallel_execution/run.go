package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	tasksCh := make(chan Task, n)
	stopCh := make(chan struct{})
	doneCh := make(chan struct{})
	var ignoreErrors bool
	mu := sync.Mutex{}
	wg := &sync.WaitGroup{}

	if m <= 0 {
		ignoreErrors = true
	}

	go func() {
		defer close(tasksCh)
		for _, t := range tasks {
			tasksCh <- t
		}
	}()

	for task := range tasksCh {
		t := task
		wg.Add(1)
		go func() {
			defer wg.Done()
			if ignoreErrors {
				_ = t()
			} else {
				err := t()
				if err != nil {
					mu.Lock()
					defer mu.Unlock()
					if m > 0 {
						if m > 1 {
							m--
						} else {
							select {
							case stopCh <- struct{}{}:
							default:
							}
						}
					}
				}
			}

		}()
	}

	go func() {
		wg.Wait()
		select {
		case doneCh <- struct{}{}:
		default:
		}
	}()

	for {
		select {
		case <-stopCh:
			return ErrErrorsLimitExceeded
		case <-doneCh:
			return nil
		}
	}
}

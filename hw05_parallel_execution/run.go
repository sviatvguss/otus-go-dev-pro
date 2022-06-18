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
	mu := sync.Mutex{}
	wg := &sync.WaitGroup{}

	go func() {
		defer close(tasksCh)
		for _, t := range tasks {
			tasksCh <- t
		}
	}()

	for task := range tasksCh {
		task := task
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := task()
			if m > 0 {
				if err != nil {
					if m > 1 {
						mu.Lock()
						m--
						mu.Unlock()
					} else {
						select {
						case stopCh <- struct{}{}:
						default:
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

package main

import (
	"errors"
	"time"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	cp <- 100
	close(cp)
	for i := 0; i < 100; i++ {
		time.Sleep(time.Millisecond * 10)
		step <- struct{}{}
	}
	close(step)
	done <- struct{}{}

	// Place your code here.
	return nil
}

package main

import (
	"errors"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	from, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	info, err := from.Stat()
	if err != nil {
		return err
	}
	len := info.Size()
	switch {
	case !info.Mode().IsRegular():
		return ErrUnsupportedFile
	case len < offset:
		return ErrOffsetExceedsFileSize
	}
	defer from.Close()

	to, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer to.Close()

	buffSize := len / 10
	cp <- len
	close(cp)
	buf := make([]byte, buffSize)
	for {
		n, err := from.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		wn, err := to.Write(buf[:n])
		if err != nil {
			return err
		}
		step <- wn
	}
	close(step)
	done <- struct{}{}

	// Place your code here.
	return nil
}

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

	from.Seek(offset, 0)
	var buffSize, bytesCount int64
	if limit != 0 {
		bytesCount = limit
	} else {
		bytesCount = len - offset
	}
	if bytesCount < 10 {
		buffSize = bytesCount
	} else {
		buffSize = bytesCount / 10
	}
	if limit != 0 {
		if limit > len-offset {
			cp <- len - offset
		} else {
			cp <- limit
		}
	} else {
		cp <- bytesCount
	}

	close(cp)
	buf := make([]byte, buffSize)
	var limitCounter int
	for {
		n, err := from.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		if limit != 0 {
			limitCounter += n
			if limitCounter > int(limit) {
				n -= limitCounter - int(limit)
			}
		}
		wn, err := to.Write(buf[:n])
		if err != nil {
			return err
		}
		step <- wn
		if limitCounter > int(limit) {
			break
		}
	}
	close(step)
	done <- struct{}{}

	// Place your code here.
	return nil
}

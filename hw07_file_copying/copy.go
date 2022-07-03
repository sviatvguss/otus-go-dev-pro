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
	switch {
	case !info.Mode().IsRegular():
		return ErrUnsupportedFile
	case info.Size() < offset:
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
		bytesCount = info.Size() - offset
	}
	if bytesCount < 10 {
		buffSize = bytesCount
	} else {
		buffSize = bytesCount / 10
	}

	if cp != nil {
		defer close(cp)
		if limit != 0 {
			if limit > info.Size()-offset {
				cp <- info.Size() - offset
			} else {
				cp <- limit
			}
		} else {
			cp <- bytesCount
		}
	}

	buf := make([]byte, buffSize)
	var limitCounter int
	for {
		n, err := from.Read(buf)
		if err != nil && !errors.Is(err, io.EOF) {
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
		if step != nil {
			step <- wn
		}
		if limit != 0 && limitCounter > int(limit) {
			break
		}
	}
	if step != nil && done != nil {
		close(step)
		done <- struct{}{}
	}

	return nil
}

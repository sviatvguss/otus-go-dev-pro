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

	buf, bytesToWrite := getBuffer(info.Size())

	if count != nil {
		sendToProgressBar(info.Size(), bytesToWrite)
	}

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
			step <- int64(wn)
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

func getBuffer(fileSize int64) (buf []byte, bytesCount int64) {
	if limit != 0 {
		bytesCount = limit
	} else {
		bytesCount = fileSize - offset
	}

	var buffSize int64
	if bytesCount < 10 {
		buffSize = bytesCount
	} else {
		buffSize = bytesCount / 10
	}
	buf = make([]byte, buffSize)
	return
}

func sendToProgressBar(fileSize, bytesCount int64) {
	defer close(count)
	switch {
	case limit != 0:
		if limit > fileSize-offset {
			count <- fileSize - offset
		} else {
			count <- limit
		}
	default:
		count <- bytesCount
	}
}

package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
)

type Environment map[string]string

func readLine(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	line, err := reader.ReadBytes('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}

	line = bytes.ReplaceAll(line, []byte("\n"), []byte(""))
	line = bytes.ReplaceAll(line, []byte("\x00"), []byte("\n"))

	return strings.TrimRight(string(line), " \t"), nil
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	var wg sync.WaitGroup
	var mx sync.Mutex

	env := Environment{}

	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, f := range fis {
		if !f.Mode().IsRegular() {
			continue
		}
		wg.Add(1)
		go func(fi os.FileInfo) {
			var line string
			var err error
			defer wg.Done()
			name := fi.Name()
			if fi.Size() > 0 {
				line, err = readLine(path.Join(dir, name))
				if err != nil {
					fmt.Printf("can not read data from %v\n", name)
					return
				}
			}
			mx.Lock()
			env[name] = line
			mx.Unlock()
		}(f)
	}
	wg.Wait()

	return env, nil
}

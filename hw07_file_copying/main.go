package main

import (
	"flag"
	"fmt"

	"github.com/cheggaaa/pb/v3"
)

var (
	from, to      string
	limit, offset int64
)

var cp chan int64
var step chan int
var done chan struct{}

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()
	// Place your code here.
	fmt.Printf("from = %v, to = %v, limit = %v, offset = %v\n", from, to, limit, offset)

	cp = make(chan int64)
	step = make(chan int)
	done = make(chan struct{})
	go func() {
		err := Copy(from, to, offset, limit)
		if err != nil {
			fmt.Println(fmt.Errorf("An error occured: %w", err))
			done <- struct{}{}
			return
		}
	}()

	len := <-cp
	bar := pb.StartNew(int(len))
	for {
		select {
		case v, ok := <-step:
			if ok {
				bar.Add(v)
			}
		case <-done:
			bar.Finish()
			return
		}
	}
}

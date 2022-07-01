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

var cp chan int
var step chan struct{}
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

	cp = make(chan int)
	step = make(chan struct{})
	done = make(chan struct{})
	go func() {
		err := Copy(from, to, offset, limit)
		if err != nil {

		}
	}()

	// create and start new bar
	bar := pb.StartNew(<-cp)

	for {
		select {
		case _, ok := <-step:
			if ok {
				bar.Increment()
			}
		case <-done:
			bar.Finish()
			return
		}
	}
}

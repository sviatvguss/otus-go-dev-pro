package main

import (
	"fmt"
	"os"
)

// https://untroubled.org/daemontools-encore/envdir.8.html
const exitCode = 111

func main() {
	env, err := ReadDir(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(exitCode)
	}
	os.Exit(RunCmd(os.Args[2:], env))
}

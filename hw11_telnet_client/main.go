package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const usage = "usage: telnet host port [--timeout=2s]\n"

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "used to specify dial timeout")
	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), usage)
	}
	flag.Parse()
	flag.Usage()
	args := flag.Args()
	if len(args) < 2 {
		log.Fatal(usage)
	}

	host := args[0]
	port := args[1]
	addr := net.JoinHostPort(host, port)
	client := NewTelnetClient(
		addr,
		*timeout,
		os.Stdin,
		os.Stdout)
	defer client.Close()

	if err := client.Connect(); err != nil {
		log.Println(err)
		return
	}

	log.Printf("Connected to %s...", addr)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		client.Send()
		err := client.Send()
		if err != nil {
			log.Println(err)
		}
		cancel()
	}()

	go func() {
		client.Receive()
		log.Println("...Connection was closed by peer")
		cancel()
	}()

	<-ctx.Done()
}

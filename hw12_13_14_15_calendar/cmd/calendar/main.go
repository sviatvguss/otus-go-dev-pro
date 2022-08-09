package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sviatvguss/otus-go-dev-pro/hw12_13_14_15_calendar/internal/app"
	"github.com/sviatvguss/otus-go-dev-pro/hw12_13_14_15_calendar/internal/logger"
	httpInternal "github.com/sviatvguss/otus-go-dev-pro/hw12_13_14_15_calendar/internal/server/http"
	storageInternal "github.com/sviatvguss/otus-go-dev-pro/hw12_13_14_15_calendar/internal/storage/fabric"
)

var cfg string

func init() {
	flag.StringVar(&cfg, "config", "./configs/config.toml", "Path to config file")
}

func main() {
	flag.Parse()
	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig(cfg)
	l, err := logger.New(config.Logger.Type, config.Logger.Directory, config.Logger.Level)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer l.Close()

	s, err := storageInternal.Create(config.Storage)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	calendar := app.New(l, s)
	server := httpInternal.NewServer(l, calendar, config.Server.Host, config.Server.Port)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			l.Error("failed to stop http server: " + err.Error())
		}
	}()

	l.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		l.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

package main

import (
	"log"

	"github.com/spf13/viper"
	"github.com/sviatvguss/otus-go-dev-pro/hw12_13_14_15_calendar/internal/storage"
)

type Config struct {
	Logger  LoggerConf
	Storage storage.StorageConfig
	Server  ServerConf
}

type LoggerConf struct {
	Level     string
	Directory string
	Type      string
}

type ServerConf struct {
	Host string
	Port string
}

func NewConfig(cfg string) Config {
	viper.SetConfigFile(cfg)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}
	c := Config{}
	if err := viper.Unmarshal(&c); err != nil {
		log.Fatalf("can not read confg: %v", err)
	}
	return c
}

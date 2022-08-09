package app

import (
	"context"

	"github.com/sviatvguss/otus-go-dev-pro/hw12_13_14_15_calendar/internal/logger"
	"github.com/sviatvguss/otus-go-dev-pro/hw12_13_14_15_calendar/internal/storage"
)

type App struct { // TODO
	storage storage.Storage
	logger  *logger.Logger
}

func New(logger *logger.Logger, storage storage.Storage) *App {
	return &App{
		storage: storage,
		logger:  logger,
	}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

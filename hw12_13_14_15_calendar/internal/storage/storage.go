package storage

import (
	"errors"
)

var (
	ErrNoEvent          = errors.New("no event")
	ErrEventsWithSameId = errors.New("same id")
	ErrAccessDenied     = errors.New("access denied")
)

type Storage interface {
	Add(event Event) (int64, error)
	Update(id int64, event Event) error
	Delete(id int64) error
	Get(id int64) (Event, error)
	EventsForDate(date string) ([]Event, error)
	EventsForWeek(date string) ([]Event, error)
	EventsForMonth(date string) ([]Event, error)
}

type StorageConfig struct {
	Name string
	Type string
	User string
	Pass string
	Host string
	Port string
}

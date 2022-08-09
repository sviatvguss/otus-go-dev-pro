package storagefabric

import (
	"errors"

	"github.com/sviatvguss/otus-go-dev-pro/hw12_13_14_15_calendar/internal/storage"
	inmemStorage "github.com/sviatvguss/otus-go-dev-pro/hw12_13_14_15_calendar/internal/storage/memory"
	sqlStorage "github.com/sviatvguss/otus-go-dev-pro/hw12_13_14_15_calendar/internal/storage/sql"
)

var ErrStorageProblem = errors.New("problem with storage occurred")

func Create(conf storage.StorageConfig) (storage.Storage, error) {
	var s storage.Storage
	t := conf.Type
	if t == "inmem" {
		s = inmemStorage.New()
	} else if t == "db" {
		st, err := sqlStorage.New(
			conf.User,
			conf.Pass,
			conf.Name,
			conf.Port,
			conf.Host,
		)
		if err != nil {
			return st, err
		}
	} else {
		return s, ErrStorageProblem
	}

	return s, nil
}

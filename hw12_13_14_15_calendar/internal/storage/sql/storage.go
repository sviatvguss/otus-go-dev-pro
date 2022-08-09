package sqlstorage

import (
	"context"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sviatvguss/otus-go-dev-pro/hw12_13_14_15_calendar/internal/storage"
	timeInternal "github.com/sviatvguss/otus-go-dev-pro/hw12_13_14_15_calendar/internal/time"
)

type Storage struct {
	db *sqlx.DB
}

func New(user, password, dbname, port, host string) (*Storage, error) {
	source := fmt.Sprintf("name=%s dbname=%s password=%s port=%s host=%s sslmode=disable",
		user,
		dbname,
		password,
		port,
		host)
	db, err := sqlx.Open("pgx", source)
	if err != nil {
		return nil, err
	}
	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) Close(ctx context.Context) error {
	return s.db.Close()
}

func (s *Storage) Add(e storage.Event) (int64, error) {
	sql := "INSERT INTO events (title, started_at, ended_at, description, user_id) VALUES (:title, :started_at, :ended_at, :description, :user_id)"

	startedAt, err := time.Parse(timeInternal.DateTimeFormat, e.DateStart)
	if err != nil {
		return 0, err
	}

	endedAt, err := time.Parse(timeInternal.DateTimeFormat, e.DateEnd)
	if err != nil {
		return 0, err
	}

	result, err := s.db.NamedExec(sql, map[string]interface{}{
		"title":       e.Title,
		"started_at":  startedAt.Unix(),
		"ended_at":    endedAt.Unix(),
		"description": e.Description,
		"user_id":     e.UserId,
	})
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Storage) Update(id int64, e storage.Event) error {
	sql := "UPDATE events SET title = :title, started_at := started_at, ended_at := ended_at, description := description WHERE id = :id"

	startedAt, err := time.Parse(timeInternal.DateTimeFormat, e.DateStart)
	if err != nil {
		return err
	}

	endedAt, err := time.Parse(timeInternal.DateTimeFormat, e.DateEnd)
	if err != nil {
		return err
	}

	result, err := s.db.NamedExec(sql, map[string]interface{}{
		"id":          id,
		"title":       e.Title,
		"description": e.Description,
		"started_at":  startedAt.Unix(),
		"ended_at":    endedAt.Unix(),
	})
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return storage.ErrNoEvent
	} else if count > 1 {
		return storage.ErrEventsWithSameId
	}

	return nil
}

func (s *Storage) Delete(id int64) error {
	sql := "DELETE FROM events WHERE id = :id"

	result, err := s.db.NamedExec(sql, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return storage.ErrNoEvent
	} else if count > 1 {
		return storage.ErrEventsWithSameId
	}

	return nil
}

func (s *Storage) Get(id int64) (storage.Event, error) {
	sql := "SELECT * FROM events WHERE id = :id"

	result, err := s.db.NamedQuery(sql, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return storage.Event{}, err
	}

	events, err := getQueryResults(result)
	if err != nil {
		return storage.Event{}, err
	}

	count := len(events)
	if count == 0 {
		return storage.Event{}, storage.ErrNoEvent
	} else if count > 1 {
		return storage.Event{}, storage.ErrEventsWithSameId
	}

	return events[0], nil
}

func (s *Storage) EventsForDate(date string) ([]storage.Event, error) {
	d, err := time.Parse(timeInternal.DateFormat, date)
	if err != nil {
		return nil, err
	}
	start, end := timeInternal.DayStartAndEnd(d)
	return s.getEventsInDateRange(start, end)
}

func (s *Storage) EventsForWeek(date string) ([]storage.Event, error) {
	d, err := time.Parse(timeInternal.DateFormat, date)
	if err != nil {
		return nil, err
	}
	start, end := timeInternal.WeekStartAndEnd(d)
	return s.getEventsInDateRange(start, end)
}

func (s *Storage) EventsForMonth(date string) ([]storage.Event, error) {
	d, err := time.Parse(timeInternal.DateFormat, date)
	if err != nil {
		return nil, err
	}
	start, end := timeInternal.MonthStartAndEnd(d)
	return s.getEventsInDateRange(start, end)
}

func (s *Storage) getEventsInDateRange(from time.Time, to time.Time) ([]storage.Event, error) {
	sql := "SELECT * FROM events WHERE started_at >= :from AND ended_at <= :to"

	queryResult, err := s.db.NamedQuery(sql, map[string]interface{}{
		"from": from.Unix(),
		"to":   to.Unix(),
	})
	if err != nil {
		return nil, err
	}

	return getQueryResults(queryResult)
}

func getQueryResults(queryResult *sqlx.Rows) ([]storage.Event, error) {
	result := make([]storage.Event, 0)

	for queryResult.Next() {
		row := make(map[string]interface{})
		err := queryResult.MapScan(row)
		if err != nil {
			return nil, err
		}
		result = append(result, createEventFromRow(row))
	}

	return result, nil
}

func createEventFromRow(row map[string]interface{}) storage.Event {
	description := ""
	if _, ok := row["description"].(string); ok {
		description = row["description"].(string)
	}
	return storage.Event{
		ID:          row["id"].(int64),
		Title:       row["title"].(string),
		DateStart:   row["started_at"].(string),
		DateEnd:     row["ended_at"].(string),
		Description: description,
		UserId:      row["user_id"].(int),
	}
}

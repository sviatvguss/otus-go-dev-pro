package memorystorage

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/sviatvguss/otus-go-dev-pro/hw12_13_14_15_calendar/internal/storage"
)

func TestStorage(t *testing.T) {
	t.Run("storage: Add", func(t *testing.T) {
		s := New()
		e := storage.Event{
			Title:       "add test",
			DateStart:   "09.08.2022 12:00:00",
			DateEnd:     "09.08.2022 17:25:00",
			Description: "supercalifragilisticexpialidocious",
			UserId:      111,
		}
		id, err := s.Add(e)
		require.Nil(t, err)
		require.EqualValues(t, 1, id)
	})

	t.Run("storage: Get", func(t *testing.T) {
		s := New()
		e := storage.Event{
			Title:       "get test",
			DateStart:   "04.07.2022 22:00:00",
			DateEnd:     "01.08.2022 17:45:00",
			Description: "You don’t spell it…you feel it.",
			UserId:      222,
		}
		id, _ := s.Add(e)
		ge, err := s.Get(id)
		require.Nil(t, err)
		require.Equal(t, "get test", ge.Title)
		require.Equal(t, "You don’t spell it…you feel it.", ge.Description)
		require.EqualValues(t, 1, ge.ID)
	})

	t.Run("storage: Update", func(t *testing.T) {
		s := New()
		e := storage.Event{
			Title:       "update test",
			DateStart:   "31.05.2021 22:00:00",
			DateEnd:     "01.07.2022 07:30:00",
			Description: "Some people care too much. I think it’s called love.",
			UserId:      333,
		}
		id, _ := s.Add(e)
		ue := storage.Event{
			Title:       "update test",
			DateStart:   "12.01.2022 13:00:00",
			DateEnd:     "13.01.2022 14:30:00",
			Description: "Some people care too much.",
			UserId:      333,
		}
		err := s.Update(id, ue)
		require.Nil(t, err)
		ge, _ := s.Get(id)
		require.Equal(t, "update test", ge.Title)
		require.Equal(t, "Some people care too much.", ge.Description)
		require.EqualValues(t, 1, ge.ID)
	})

	t.Run("storage: Delete", func(t *testing.T) {
		s := New()
		e := storage.Event{
			Title:       "delete test",
			DateStart:   "18.03.2022 22:00:00",
			DateEnd:     "01.04.2022 22:00:00",
			Description: "Another one",
			UserId:      444,
		}
		id, _ := s.Add(e)
		err := s.Delete(id)
		require.Nil(t, err)
		_, err = s.Get(id)
		require.ErrorIs(t, err, storage.ErrNoEvent)
	})

	t.Run("storage: events for date", func(t *testing.T) {
		s := New()
		e1 := storage.Event{
			Title:       "yes",
			DateStart:   "03.04.2022 05:00:00",
			DateEnd:     "03.04.2022 22:30:00",
			Description: "that one",
			UserId:      555,
		}
		e2 := storage.Event{
			Title:       "yes",
			DateStart:   "03.04.2022 23:00:00",
			DateEnd:     "04.04.2022 21:30:00",
			Description: "and that one",
			UserId:      555,
		}
		e3 := storage.Event{
			Title:       "no",
			DateStart:   "02.04.2022 11:00:00",
			DateEnd:     "02.04.2022 21:30:00",
			Description: "wrong event",
			UserId:      555,
		}
		s.Add(e1)
		s.Add(e2)
		s.Add(e3)
		events, err := s.EventsForDate("03.04.2022")
		require.Nil(t, err)
		require.Len(t, events, 2)
	})

	t.Run("storage: events for week", func(t *testing.T) {
		s := New()
		e1 := storage.Event{
			Title:       "yes",
			DateStart:   "14.02.2022 11:00:00",
			DateEnd:     "14.02.2022 21:30:00",
			Description: "that one",
			UserId:      666,
		}
		e2 := storage.Event{
			Title:       "no",
			DateStart:   "01.11.2023 22:00:00",
			DateEnd:     "01.12.2023 08:30:00",
			Description: "wrong",
			UserId:      666,
		}
		e3 := storage.Event{
			Title:       "no",
			DateStart:   "22.12.2023 22:00:00",
			DateEnd:     "29.12.2023 23:30:00",
			Description: "wrong",
			UserId:      666,
		}
		s.Add(e1)
		s.Add(e2)
		s.Add(e3)
		events, err := s.EventsForWeek("14.02.2022")
		require.Nil(t, err)
		require.Len(t, events, 1)
	})

	t.Run("storage: events for month", func(t *testing.T) {
		s := New()
		e1 := storage.Event{
			Title:       "yes",
			DateStart:   "13.02.2022 10:30:00",
			DateEnd:     "13.02.2022 14:30:00",
			Description: "that one",
			UserId:      777,
		}
		e2 := storage.Event{
			Title:       "no",
			DateStart:   "12.12.2022 22:00:00",
			DateEnd:     "22.12.2022 23:30:00",
			Description: "wrong",
			UserId:      777,
		}
		e3 := storage.Event{
			Title:       "yes",
			DateStart:   "21.02.2022 01:00:00",
			DateEnd:     "27.02.2022 23:30:00",
			Description: "and that one",
			UserId:      777,
		}
		s.Add(e1)
		s.Add(e2)
		s.Add(e3)
		events, err := s.EventsForMonth("01.02.2022")
		require.Nil(t, err)
		require.Len(t, events, 2)
	})
}

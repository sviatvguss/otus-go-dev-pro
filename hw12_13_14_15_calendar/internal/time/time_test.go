package time

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	dateFormat     = "02.01.2006"
	datetimeFormat = "02.01.2006 15:04"
)

func TestTime(t *testing.T) {
	t.Run("day", func(t *testing.T) {
		day, _ := time.Parse(dateFormat, "09.08.2022")
		start, end := DayStartAndEnd(day)
		require.Equal(t, "09.08.2022 00:00", start.Format(datetimeFormat))
		require.Equal(t, "09.08.2022 23:59", end.Format(datetimeFormat))
	})

	t.Run("week", func(t *testing.T) {
		day, _ := time.Parse(dateFormat, "05.08.2022")
		start, end := WeekStartAndEnd(day)
		require.Equal(t, "01.08.2022 00:00", start.Format(datetimeFormat))
		require.Equal(t, "07.08.2022 23:59", end.Format(datetimeFormat))
	})

	t.Run("month", func(t *testing.T) {
		day, _ := time.Parse(dateFormat, "09.08.2022")
		start, end := MonthStartAndEnd(day)

		require.Equal(t, "01.08.2022 00:00", start.Format(datetimeFormat))
		require.Equal(t, "31.08.2022 23:59", end.Format(datetimeFormat))
	})
}

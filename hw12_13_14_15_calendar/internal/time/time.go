package time

import "time"

const (
	DateFormat        = "02.01.2006"
	DateTimeFormat    = "02.01.2006 15:04:05"
	LogDateTimeFormat = "02/01/2006:15:04:05 +0700"
)

func DateInRange(date time.Time, start time.Time, end time.Time) bool {
	return date.After(start) && date.Before(end)
}

func DayStartAndEnd(date time.Time) (time.Time, time.Time) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, date.Location())
	return start, end
}

func WeekStartAndEnd(date time.Time) (time.Time, time.Time) {
	if date.Weekday() != time.Monday {
		if date.Weekday() == time.Sunday {
			date = date.Add(-time.Hour * 24 * 6)
		} else {
			weekday := int(date.Weekday()) - 1
			date = date.Add(-time.Hour * 24 * time.Duration(weekday))
		}
	}
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	date = date.Add(time.Hour * 24 * 6)
	end := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, date.Location())
	return start, end
}

func MonthStartAndEnd(date time.Time) (time.Time, time.Time) {
	start := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)
	return start, end
}

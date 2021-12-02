package datetime

import "time"

func BeginningOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0,0,0,0, time.UTC)
}

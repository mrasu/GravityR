package lib

import "time"

func NormalizeTimeByHour(t time.Time) time.Time {
	return time.Date(
		t.Year(), t.Month(), t.Day(),
		t.Hour(), 0, 0, 0,
		t.Location(),
	)
}

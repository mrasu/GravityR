package lib

import "time"

func NormalizeTimeByHour(t time.Time) time.Time {
	return time.Date(
		t.Year(), t.Month(), t.Day(),
		t.Hour(), 0, 0, 0,
		t.Location(),
	)
}

func GenerateTimeRanges(start, end time.Time, interval time.Duration) []time.Time {
	var res []time.Time
	for t := start; t.Before(end); t = t.Add(interval) {
		res = append(res, t)
	}

	return res
}

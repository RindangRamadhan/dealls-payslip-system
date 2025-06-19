package utils

import "time"

func CountWorkingDays(start, end time.Time) int {
	count := 0
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		if d.Weekday() != time.Saturday && d.Weekday() != time.Sunday {
			count++
		}
	}
	return count
}

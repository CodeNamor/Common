package datetime

import "time"

// CalcDays returns the number of days between two dates.
func CalcDays(from time.Time, to time.Time) int {
	invalid := InvalidTime()
	if from.Equal(invalid) {
		return 0
	}
	if to.Equal(invalid) {
		now := time.Now()
		to = now
	}
	totalDays := float32(to.Sub(from) / (24 * time.Hour))

	return int(totalDays)
}

package datetime

import "time"

// CalcBusinessDays calculates the business days between two dates, not accounting for the holidays
// The calculation is inclusive for both dates
func CalcBusinessDays(from time.Time, to time.Time) int {
	invalid := InvalidTime()
	if from.Equal(invalid) {
		return 0
	}
	if to.Equal(invalid) {
		now := time.Now()
		to = now
	}
	totalDays := float32(to.Sub(from) / (24 * time.Hour))
	weekDays := float32(from.Weekday()) - float32(to.Weekday())
	businessDays := int(1 + (totalDays*5-weekDays*2)/7)
	if to.Weekday() == time.Saturday {
		businessDays--
	}
	if from.Weekday() == time.Sunday {
		businessDays--
	}

	return businessDays
}

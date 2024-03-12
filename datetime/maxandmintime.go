package datetime

import "time"

// MinTime returns the time.Time representation of 0001/01/01.
func MinTime() time.Time {
	t, _ := time.Parse(time.RFC3339, "")
	return t
}

// MaxTime returns the time.Time representation of 9999/12/31.
func MaxTime() time.Time {
	t, _ := time.Parse(time.RFC3339, "9999-12-31T23:59:59.999Z")
	return t
}

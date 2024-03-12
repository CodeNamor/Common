package datetime

import (
	"time"
)

// ToSearchDateFormat returns a formatted date ("2006-01-02")
func ToSearchDateFormat(t time.Time) string {
	return t.Format(SearchDateFormat)
}

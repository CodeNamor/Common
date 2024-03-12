package datetime

import (
	"strings"
	"time"
)

const SearchDateFormat = "2006-01-02"

// StringToDate converts a string date format to time.Time if supported, otherwise a utils.InvalidTime() is returned.
func StringToDate(dateStr string) time.Time {

	dateStr = strings.TrimSpace(dateStr)

	if dateStr == "" {
		return InvalidTime()
	}

	for _, f := range []string{
		time.RFC3339,     // 2018-01-01T12:00:00Z00:00
		SearchDateFormat, // 2018-01-01
		"2006-1-2T15:4:5Z",
		"2006-1-2",
		"1/2/2006",
		"2006-01-02T15:04:05.000+0000", //2006-01-02T15:04:05Z07:00"
		"2006-01-02T15:04:05",
		"15:04",
	} {
		if parsedDate, err := time.Parse(f, dateStr); err == nil {
			return parsedDate
		}
	}
	return InvalidTime()
}

// ToSearchDateFormatFromString converts a string date to "2006-01-02" format
func ToSearchDateFormatFromString(date string) string {
	if strings.TrimSpace(date) != "" {
		return ToSearchDateFormat(StringToDate(date))
	}
	return ""
}

// InvalidTime is used to compare Time to the result of an invalid string value conversion
func InvalidTime() time.Time {
	t, _ := time.Parse(time.RFC3339, "")
	return t
}

// ConvertDateFormat converts a date string from one layout to another layout. An invalid 'from' layout will return the original date string. An invalid 'to' layout will return the invalid 'to' layout string.
func ConvertDateFormat(date, from, to string) string {
	t, err := time.Parse(from, date)
	if err != nil {
		return date
	}
	return t.Format(to)
}

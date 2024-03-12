package datetime

import (
	"fmt"
	"testing"
	"time"

	"github.com/CodeNamor/Common/utils"
	"github.com/stretchr/testify/assert"
)

func Test_StringToDate(t *testing.T) {
	tests := []struct {
		name  string
		date  string
		valid bool
	}{
		{"W3C DTF", "2018-01-01", true},
		{"W3C DTF no zeros", "2018-1-1", true},
		{"RFC3339", "2018-04-16T13:08:58-05:00", true},
		{"Invalid1", "x", false},
		{"EmptyInput", "", false},
		{"ISO 8601", "01/01/2018", true},
		{"Test1", "7/9/2014 12:00:00 AM", true},
	}
	invalidTime := InvalidTime()
	for _, test := range tests {
		sw := utils.StopWatch{}
		sw.Start()
		date := StringToDate(test.date)
		//visual
		fmt.Println(date, sw.Elapsed())
		if date.Equal(invalidTime) && test.valid {
			t.Fatalf("Test %v failed, could not parse %v", test.name, test.date)
		}
		if !date.Equal(invalidTime) && !test.valid {
			t.Fatalf("Test %v failed, got %v for %v, but was not supposed to", test.name, date, test.date)
		}
	}
}

// Test_ToSearchDateFormatFromString Testing ToSearchDateFormatFromString Function
func Test_ToSearchDateFormatFromString(t *testing.T) {

	inputDate := "20160608"
	expectedDate := "2016-06-08"
	result := ToSearchDateFormatFromString(inputDate)
	if result != expectedDate {
		t.Errorf("Test Failed. expected '%s' but got '%s'.", expectedDate, result)
	}
}

func Test_ConvertDateFormat(t *testing.T) {
	tests := []struct {
		name, date, from, to, expected string
	}{
		{
			name:     "convert a valid date string's format to a different format",
			date:     "20190214",
			from:     "20060102",
			to:       time.RFC3339,
			expected: "2019-02-14T00:00:00Z",
		},
		{
			name:     "convert a valid date string with an invalid format",
			date:     "20190214",
			from:     "abc",
			to:       time.RFC3339,
			expected: "20190214",
		},
		{
			name:     "converting an invalid date string will return itself",
			date:     "asdf",
			from:     time.RFC3339,
			to:       time.RFC3339,
			expected: "asdf",
		},
		{
			name:     "converting a valid date string with an invalid outgoing format will return the invalid format",
			date:     "20190214",
			from:     "20060102",
			to:       "asdf",
			expected: "asdf",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fdate := ConvertDateFormat(test.date, test.from, test.to)
			assert.Equal(t, test.expected, fdate)
		})
	}
}

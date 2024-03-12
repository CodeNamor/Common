package datetime

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_CalcBusinessDays(t *testing.T) {
	today := time.Now()
	var lastBusinessDay string
	if today.Weekday() == time.Monday {
		lastBusinessDay = time.Now().Add(-72 * time.Hour).Format(time.RFC3339)
	} else {
		lastBusinessDay = time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
	}

	for _, test := range []struct {
		start, end string
		days       int
	}{
		{"2016-06-01T00:00:00Z", "2016-06-06T00:00:00Z", 4},
		{"2016-08-29T00:00:00Z", "2016-09-04T00:00:00Z", 5},
		{"2016-08-27T00:00:00Z", "2016-09-04T00:00:00Z", 5},
		{"2016-08-27T00:00:00Z", "2016-08-29T00:00:00Z", 1},
		{"2016-01-01T00:00:00Z", "2016-01-31T00:00:00Z", 21},
		{"2016-02-01T00:00:00Z", "2016-02-29T00:00:00Z", 21},
		{"2016-03-01T00:00:00Z", "2016-03-31T00:00:00Z", 23},
		{"2016-04-01T00:00:00Z", "2016-04-30T00:00:00Z", 21},
		{"2016-05-01T00:00:00Z", "2016-05-31T00:00:00Z", 22},
		{"2016-06-01T00:00:00Z", "2016-06-30T00:00:00Z", 22},
		{"2016-07-01T00:00:00Z", "2016-07-31T00:00:00Z", 21},
		{"2016-08-01T00:00:00Z", "2016-08-31T00:00:00Z", 23},
		{"2016-09-01T00:00:00Z", "2016-09-30T00:00:00Z", 22},
		{"2016-10-01T00:00:00Z", "2016-10-31T00:00:00Z", 21},
		{"2016-11-01T00:00:00Z", "2016-11-30T00:00:00Z", 22},
		{"2016-12-01T00:00:00Z", "2016-12-31T00:00:00Z", 22},
		{"2016-01-01T00:00:00Z", "2016-12-31T00:00:00Z", 261},
		{"x", "2018-01-01T00:00:00Z", 0},
		{lastBusinessDay, "x", 2},
		{"2016-05-18T00:00:00Z", "2016-05-20T00:00:00Z", 3},
		{"2017-10-17T00:00:00Z", "2017-11-17T00:00:00Z", 24},
		{"", "", 0},
		{lastBusinessDay, "", 2},
	} {
		t.Run(fmt.Sprintf("%s-%s", test.start, test.end), func(t *testing.T) {
			assert.Equal(t, test.days, CalcBusinessDays(StringToDate(test.start), StringToDate(test.end)))
		})
	}
}

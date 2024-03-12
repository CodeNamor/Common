package datetime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	start, _   = time.Parse("2006-01-02", "2019-01-01")
	end, _     = time.Parse("2006-01-02", "2019-01-03")
	between, _ = time.Parse("2006-01-02", "2019-01-02")
)

func TestTimeBoundaries(t *testing.T) {
	tests := []struct {
		name         string
		timeboundary TimeBound
		start, end   *time.Time
		expected     bool
	}{
		{
			name:         "Start and end bounds are nil. There are no bounds defined.",
			timeboundary: createMockedTimeBoundary(&start, &end),
			start:        nil,
			end:          nil,
			expected:     true,
		},
		{
			name:         "Start and end bounds are provided. Record's start and end times are within bounds.",
			timeboundary: createMockedTimeBoundary(&start, &end),
			start:        &start,
			end:          &end,
			expected:     true,
		},
		{
			name:         "Start and end bounds are provided. Record's start time in bounds but the end is out of bounds.",
			timeboundary: createMockedTimeBoundary(&start, addDays(end, 1)),
			start:        &start,
			end:          &end,
			expected:     true,
		},
		{
			name:         "Start and end bounds are provided. Record's time is between the start and end.",
			timeboundary: createMockedTimeBoundary(&between, &end),
			start:        &start,
			end:          &end,
			expected:     true,
		},
		{
			name:         "Record's time is within the bounds of the start and end. Record's time is the same as end.",
			timeboundary: createMockedTimeBoundary(&end, &end),
			start:        &start,
			end:          &end,
			expected:     true,
		},
		{
			name:         "Record's time is within the bounds of the start and end. Record's time is the same as end.",
			timeboundary: createMockedTimeBoundary(&end, &end),
			start:        &start,
			end:          &end,
			expected:     true,
		},
		{
			name:         "Start and end bounds are provided. Record's end time is before the start bound.",
			timeboundary: createMockedTimeBoundary(&start, addDays(start, -1)),
			start:        &start,
			end:          &end,
			expected:     false,
		},
		{
			name:         "Start and end bounds are provided. Record's start time is after the end bound.",
			timeboundary: createMockedTimeBoundary(addDays(end, 1), &end),
			start:        &start,
			end:          &end,
			expected:     false,
		},
		{
			name:         "Only start bound is defined. Record's end time is the same as the start bound.",
			timeboundary: createMockedTimeBoundary(&start, &start),
			start:        &start,
			end:          nil,
			expected:     true,
		},
		{
			name:         "Only start bound is defined. Record's end time is after the start bound.",
			timeboundary: createMockedTimeBoundary(&start, &end),
			start:        &start,
			end:          nil,
			expected:     true,
		},
		{
			name:         "Only start bound is defined. Record's end time is before the start bound.",
			timeboundary: createMockedTimeBoundary(&start, addDays(start, -1)),
			start:        &start,
			end:          nil,
			expected:     false,
		},
		{
			name:         "Only end bound is defined. Record's start time is the same as the end bound.",
			timeboundary: createMockedTimeBoundary(&end, &end),
			start:        nil,
			end:          &end,
			expected:     true,
		},
		{
			name:         "Only end bound is defined. Record's start time is before the end bound.",
			timeboundary: createMockedTimeBoundary(addDays(end, -1), &end),
			start:        nil,
			end:          &end,
			expected:     true,
		},
		{
			name:         "Only end bound is defined. Record's start time is after the end bound.",
			timeboundary: createMockedTimeBoundary(addDays(end, 1), &end),
			start:        nil,
			end:          &end,
			expected:     false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// 1. Test the bounds
			actual, err := TimeRangeInBounds(test.timeboundary, test.start, test.end)

			// 2. Make sure no errors occur
			assert.Nil(t, err, "These tests should be successful.")

			// 3. Validate the response
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestTimeBoundaries_expectedErrors(t *testing.T) {
	tests := []struct {
		name         string
		timeboundary TimeBound
		start, end   *time.Time
	}{
		{
			name:         "Record's start time is a zero value",
			timeboundary: createMockedTimeBoundary(nil, &end),
			start:        &start,
			end:          &end,
		},
		{
			name:         "Record's end time is a zero value",
			timeboundary: createMockedTimeBoundary(nil, &end),
			start:        &start,
			end:          &end,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// 1. Test the bounds
			_, err := TimeRangeInBounds(test.timeboundary, test.start, test.end)

			// 2. Make sure there is an error
			assert.NotNil(t, err, "These tests should generate errors.")
		})
	}
}

// **************************
// Mocking the timeboundary
// **************************
type MockedTimeBoundary struct {
	mock.Mock
}

func (mock *MockedTimeBoundary) Start() ComparableDate {
	args := mock.MethodCalled("Start")
	return args.Get(0).(ComparableDate)
}

func (mock *MockedTimeBoundary) End() ComparableDate {
	args := mock.MethodCalled("End")
	return args.Get(0).(ComparableDate)
}

func createMockedTimeBoundary(start, end *time.Time) *MockedTimeBoundary {
	boundary := new(MockedTimeBoundary)
	if start != nil {
		boundary.On("Start").Return(ComparableDate{*start})
	} else {
		boundary.On("Start").Return(ComparableDate{time.Time{}})
	}

	if end != nil {
		boundary.On("End").Return(ComparableDate{*end})
	} else {
		boundary.On("End").Return(ComparableDate{time.Time{}})
	}

	return boundary
}

func addDays(currentTime time.Time, days int) *time.Time {
	t := currentTime.Add(time.Hour * time.Duration(24*days))
	return &t
}

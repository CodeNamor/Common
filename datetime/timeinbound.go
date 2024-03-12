package datetime

import (
	"errors"
	"time"
)

// ComparableDate provides an interface to compare a record's date with a provided date.
type ComparableDate struct {
	Date time.Time
}

// EqualOrAfter determines whether the record's date is equal or after the provided date.
func (dt ComparableDate) EqualOrAfter(date time.Time) bool {
	return dt.Date.Equal(date) || dt.Date.After(date)
}

// EqualOrBefore determines whether the record's date is equal or before the provided date.
func (dt ComparableDate) EqualOrBefore(date time.Time) bool {
	return dt.Date.Equal(date) || dt.Date.Before(date)
}

// TimeBound provides an interface for records to specify their date boundaries.
type TimeBound interface {
	Start() ComparableDate
	End() ComparableDate
}

// TimeRangeInBounds determines whether the record's time range is included in the specified time ranges
func TimeRangeInBounds(record TimeBound, start, end *time.Time) (bool, error) {
	if record == nil {
		return false, errors.New("time bound check failed. the record is nil")
	}
	if record.End().Date.IsZero() || record.Start().Date.IsZero() {
		return false, errors.New("the record cannot contain time.Time zero values")
	}

	startBound := time.Time{}
	endBound := time.Time{}

	if start != nil {
		startBound = *start
	}

	if end != nil {
		endBound = *end
	}

	if startBound.IsZero() && endBound.IsZero() {
		return true, nil
	} else if !startBound.IsZero() && endBound.IsZero() {
		return record.End().EqualOrAfter(startBound), nil
	} else if !endBound.IsZero() && startBound.IsZero() {
		return record.Start().EqualOrBefore(endBound), nil
	} else {
		return record.Start().EqualOrBefore(endBound) && record.End().EqualOrAfter(startBound), nil
	}
}

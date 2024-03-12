package datetime

import (
	"testing"
)

func Test_ToSearchDateFormat(t *testing.T) {

	inputDate := "20160608"
	expectedDate := "2016-06-08"
	result := ToSearchDateFormat(StringToDate(inputDate))
	if result != expectedDate {
		t.Errorf("Test Failed. expected '%s' but got '%s'.", expectedDate, result)
	}
}

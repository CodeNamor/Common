package utils

import (
	"log"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_StringInSlice(t *testing.T) {
	test := []string{"a", "b", "c", "d"}
	if !StringInSlice("a", test) {
		t.Error("Could not find a")
	}
	if StringInSlice("f", test) {
		t.Error("Found f?!?")
	}
	if StringInSlice("C", test) {
		t.Error("Found C?!?")
	}
}

func Test_HostName(t *testing.T) {
	h := HostName()
	if h == "" {
		t.Fatal("HostName was not delivered")
	}
}

func Test_StringFloat64Sub(t *testing.T) {
	tests := []struct {
		value1    string
		value2    string
		precision int
		expected  string
	}{
		{"2.013", "1", 2, "1.01"},
		{"2.013", "1x", 3, "0.000"},
		{"2.013x", "1", 2, "0.00"},
		{"2.013", "3", 2, "-0.99"},
		{"2.016", "3", 2, "-0.98"},
	}
	//test
	for _, test := range tests {
		r := StringFloat64Sub(test.value1, test.value2, test.precision)
		if r != test.expected {
			t.Fatalf("%v - %v resulted in %v, expected was %v", test.value1, test.value2, r, test.expected)
		}
	}
}

func Test_StringToDecimalString(t *testing.T) {
	tests := []struct {
		name     string //test name
		s        string //string to convert
		n        int    //precision
		expected string //result expected
	}{
		{"2.013", "2.013", 2, "2.01"},
		{"2.016", "2.016", 2, "2.02"},
		{"2.015", "2.015", 2, "2.02"},
		{"2.014", "2.014", 2, "2.01"},
		{"EmptyString", "", 2, ""},
		{"BadString", "2.x14", 2, ""},
		{"NegativePrecision", "2.014", -2, "2"},
		{"2.014 precision 1", "2.014", 1, "2.0"},
		{"2", "2", 2, "2.00"},
	}
	testName := ""

	for _, test := range tests {
		if testName == "" || testName == test.name {
			r := StringToDecimalString(test.s, test.n)
			if r != test.expected {
				t.Fatalf("Test %v failed. Expected %v, got %v.", test.name, test.expected, r)
			}
		}
	}

}

func Test_IsDateInSpan(t *testing.T) {
	tests := []struct {
		name      string
		startDate time.Time
		endDate   time.Time
		date      time.Time
		expected  bool
	}{
		{
			name:      "date on start date",
			startDate: time.Date(2018, time.January, 10, 0, 0, 0, 0, time.UTC),
			endDate:   time.Date(2018, time.February, 10, 0, 0, 0, 0, time.UTC),
			date:      time.Date(2018, time.January, 10, 0, 0, 0, 0, time.UTC),
			expected:  true,
		},
		{
			name:      "date on end date",
			startDate: time.Date(2018, time.January, 10, 0, 0, 0, 0, time.UTC),
			endDate:   time.Date(2018, time.February, 10, 0, 0, 0, 0, time.UTC),
			date:      time.Date(2018, time.February, 10, 0, 0, 0, 0, time.UTC),
			expected:  true,
		},
		{
			startDate: time.Date(2018, time.January, 10, 0, 0, 0, 0, time.UTC),
			endDate:   time.Date(2018, time.February, 10, 0, 0, 0, 0, time.UTC),
			date:      time.Date(2018, time.January, 15, 0, 0, 0, 0, time.UTC),
			expected:  true,
		},
		{
			name:      "date before start date",
			startDate: time.Date(2018, time.January, 10, 0, 0, 0, 0, time.UTC),
			endDate:   time.Date(2018, time.February, 10, 0, 0, 0, 0, time.UTC),
			date:      time.Date(2018, time.January, 9, 0, 0, 0, 0, time.UTC),
			expected:  false,
		},
		{
			name:      "date after end date",
			startDate: time.Date(2018, time.January, 10, 0, 0, 0, 0, time.UTC),
			endDate:   time.Date(2018, time.February, 10, 0, 0, 0, 0, time.UTC),
			date:      time.Date(2018, time.February, 11, 0, 0, 0, 0, time.UTC),
			expected:  false,
		},
	}
	testName := ""

	for _, test := range tests {
		if testName == "" || testName == test.name {
			result := IsDateInSpan(test.startDate, test.endDate, test.date)
			if test.expected != result {
				t.Fatalf("%v: expected result of %v, got %v", test.name, test.expected, result)
			}
		}
	}
}

func Test_GetFirstN(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		n        int
		expected string
	}{
		{"More than N characters", "abcdef", 3, "abc"},
		{"Exactly N characters", "abc", 3, "abc"},
		{"Less than N characters", "ab", 3, "ab"},
		{"N is 0", "abc", 0, ""},
		{"Negative N", "abc", -1, ""},
		{"Empty string", "", 3, ""},
	}
	testName := ""

	for _, test := range tests {
		if testName == "" || testName == test.name {
			result := GetFirstN(test.input, test.n)
			if test.expected != result {
				t.Fatalf("%v: expected %v, got %v", test.name, test.expected, result)
			}
		}
	}
}

func Test_GivesTheDefaultValueIfActualIsEmpty(t *testing.T) {
	actualString := ""
	actualValidString := "B4321"
	defaultString := "A12345"
	result := DefaultStringIfActualIsBlank(actualString, defaultString)
	if result != defaultString {
		t.Errorf("Expected to get an defaultstring string. Got '%s' want '%s'", result, defaultString)
	}
	result = DefaultStringIfActualIsBlank(actualValidString, defaultString)
	if result != actualValidString {
		t.Errorf("Expected to get an defaultstring string. Got '%s' want '%s'", result, actualValidString)
	}
}

func Test_Distinct(t *testing.T) {

	fruitsList := []string{"Apple", "Orange", "Apple", "Mango", "Banana", "Orange"}

	result := Distinct(fruitsList)

	assert.Equal(t, 4, len(result))
	assert.Equal(t, fruitsList[0], result[0])
	assert.Equal(t, fruitsList[1], result[1])
	assert.Equal(t, fruitsList[3], result[2])
	assert.Equal(t, fruitsList[4], result[3])
}

func Test_ConvertBoolToYN(t *testing.T) {

	tests := []struct {
		input    bool
		expected string
	}{
		{true, "Y"}, {false, "N"},
	}

	for _, r := range tests {
		result := ConvertBoolToYN(r.input)
		assert.Equal(t, r.expected, result)
	}
}

func Test_IsListEmpty(t *testing.T) {

	tests := []struct {
		input    []string
		expected bool
	}{
		{[]string{""}, true},
		{[]string{"Dummy"}, false},
	}

	for _, r := range tests {
		result := IsListEmpty(r.input)
		assert.Equal(t, r.expected, result)
	}
}

func Test_StringWithHyphenSeprated(t *testing.T) {
	tests := []struct {
		name     string
		stringA  string
		stringB  string
		expected string
	}{
		{
			name:     "should return a string with hyphen separated strings",
			stringA:  "concatenate first stringA",
			stringB:  "string B",
			expected: "concatenate first stringA - string B",
		},
		{
			name:     "should return a string when stringB have hyphen though stringA is an empty string",
			stringA:  "",
			stringB:  "string B",
			expected: " - string B",
		},
		{
			name:     "should return a string stringA with hyphen though stringB is an empty string",
			stringA:  "concatenate first string",
			stringB:  "",
			expected: "concatenate first string - ",
		},
		{
			name:     "should return an empty string with no hyphen when both strings are empty",
			stringA:  "",
			stringB:  "",
			expected: "",
		},
	}

	testName := ""

	for _, test := range tests {
		if testName == "" || testName == test.name {
			result := StringsWithHyphenseparated(test.stringA, test.stringB)
			log.Printf("Test %v: result :%v", test.name, result)
			if result != test.expected {
				t.Fatalf("Test %v: Expected %v , got %v.", test.name, test.expected, result)
			}
		}
	}

}

func Test_SeekProjectRoot(t *testing.T) {
	assert.True(t, strings.HasSuffix(SeekProjectRoot(""), "common/"))
}

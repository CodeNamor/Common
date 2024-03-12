package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// HostName returns HostName associated with that particular System eg: cncv.
func HostName() string {
	var name string
	if n, err := os.Hostname(); err != nil {
		log.Print("Error retrieving HostName.")
		name = ""
	} else {
		name = n
	}

	return name
}

// IsDateInSpan returns true if the provided date is between or on the provided start and end dates
// otherwise returns false.
func IsDateInSpan(start, end, date time.Time) bool {
	return (start.Before(date) || start.Equal(date)) && (end.After(date) || end.Equal(date))
}

// SeekProjectRoot searches up the file structure to find the directory containing go.mod and returns the absolute path, always ending in "/.
func SeekProjectRoot(path string) string {
	if fileExists(path + "go.mod") {
		if file, err := filepath.Abs(path); err != nil {
			panic("Error determining absolute path of " + path)
		} else {
			return file + "/"
		}
	}
	return SeekProjectRoot("../" + path)
}

func fileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

//*********************************************************************
//STRING HELPERS
//*********************************************************************

// StringToDecimalString takes a string input, attempts to convert it to decimal
// If the conversion is successful, returns the string representation of the decimal
// n specifies the precision, how many digits after the decimal point.
func StringToDecimalString(s string, n int) string {
	if n < 0 {
		n = 0
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	number, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return ""
	}

	return strconv.FormatFloat(number, 'f', n, 64)
}

// StringInSlice checks if []string already has an element s
func StringInSlice(s string, list []string) bool {
	for _, e := range list {
		if e == s {
			return true
		}
	}
	return false
}

// StringFloat64Sub takes two input string parameters, a and b, treats them as float64 values
// and returns the difference (a-b) as a string format with n precision
func StringFloat64Sub(a string, b string, n int) string {
	result := 0.0
	af, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return strconv.FormatFloat(result, 'f', n, 64)
	}
	bf, err := strconv.ParseFloat(b, 64)
	if err != nil {
		return strconv.FormatFloat(result, 'f', n, 64)
	}
	return strconv.FormatFloat(af-bf, 'f', n, 64)
}

// GetFirstN returns the first n characters of s, or s itself, if s has less than n characters
// returns empty string if n is equal or less than 0, regardless off s
func GetFirstN(s string, n int) string {
	if n < 0 {
		return ""
	}
	result := s
	if len(s) > n && s != "" && n >= 0 {
		result = s[0:n]
	}
	return result
}

// DefaultStringIfActualIsBlank will give default string if actual is empty
// else return actual value
func DefaultStringIfActualIsBlank(actual string, defaultString string) string {

	if actual == "" {
		actual = defaultString
		return actual
	}
	return actual
}

// Distinct removes the duplicates from the List and returns a List with Unique values
func Distinct(arr []string) []string {
	distinct := make([]string, 0, len(arr))
	m := make(map[string]bool)
	for _, x := range arr {
		if _, ok := m[x]; !ok { // if didn't match existing
			m[x] = true
			distinct = append(distinct, x)
		}
	}
	return distinct
}

// ConvertBoolToYN converts a input if true to Y or false to N
func ConvertBoolToYN(b bool) string {
	if b {
		return "Y"
	}
	return "N"
}

// IsListEmpty checks if List is empty in the case of sourceSystem: [""]
// if empty returns true else false
func IsListEmpty(input []string) (emptyFlag bool) {

	filteredList := make([]string, 0)

	for _, r := range input {
		if r != "" {
			filteredList = append(filteredList, r)
		}
	}
	return len(filteredList) == 0
}

// StringsWithHyphenseparated accepts two inputs strings and validates if the strings are null
// and return empty string else return string in "a - b" format.
func StringsWithHyphenseparated(a, b string) string {

	if a != "" || b != "" {
		return fmt.Sprintf("%v - %v", a, b)
	}
	return ""
}

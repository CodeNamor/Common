package path

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Resolve(t *testing.T) {
	testcases := []struct {
		name     string
		paths    []string
		expected string
	}{
		{
			name:     "0 - nil paths, returns empty string",
			paths:    nil,
			expected: "",
		},
		{
			name:     "1 - empty slice, returns empty string",
			paths:    []string{},
			expected: "",
		},
		{
			name:     "2 - empty segments, returns empty string",
			paths:    []string{"", "", ""},
			expected: "",
		},
		{
			name: "3 - all relative paths, returns cleaned joined path",
			paths: []string{
				"a/b",
				"./c/d",
				"e/./../f",
				"g",
				"",
				"../h/i",
				"j/..",
				"k/l/m/",
				"n",
			},
			expected: "a/b/c/d/f/h/i/k/l/m/n",
		},
		{
			name: "4 - white space in or around segments is retained",
			paths: []string{
				" a / b ",
				" c",
				" ",
				"d ",
			},
			expected: " a / b / c/ /d ",
		},
		{
			name: "5 - abs + rel paths, returns cleaned abs path",
			paths: []string{
				"/abs/a",
				"b/c",
				"d",
			},
			expected: "/abs/a/b/c/d",
		},
		{
			name: "6 - rel1, abs1, abs2 + rel2 paths, returns cleaned abs2+ path",
			paths: []string{
				"rel1/foo",
				"/abs1/a",
				"/abs2/b",
				"c/./d",
				"e/f/..",
			},
			expected: "/abs2/b/c/d/e",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, Resolve(tc.paths...))
		})
	}
}
